package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4/pgxpool"
	apiPkg "github.com/maximgoltsov/botproject/internal/api"
	botPkg "github.com/maximgoltsov/botproject/internal/pkg/bot"

	cmdAddPkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command/add"
	cmdDeletePkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command/delete"
	cmdEditPkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command/edit"
	cmdHelpPkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command/help"
	cmdListPkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command/list"
	productPkg "github.com/maximgoltsov/botproject/internal/pkg/core/product"
	productTypePkg "github.com/maximgoltsov/botproject/internal/pkg/core/productType"
	pb "github.com/maximgoltsov/botproject/pkg/api"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	const (
		Host     = "localhost"
		Port     = 7776
		User     = "user"
		Password = "password"
		DBname   = "products"

		MaxConnIdleTime = time.Minute
		MaxConnLifetime = time.Hour
		MinConns        = 2
		MaxConns        = 4
	)

	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, Port, User, Password, DBname)

	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		log.Fatal("can't connect to database", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("ping database error", err)
	}

	config := pool.Config()
	config.MaxConnIdleTime = MaxConnIdleTime
	config.MaxConnLifetime = MaxConnLifetime
	config.MinConns = MinConns
	config.MaxConns = MaxConns

	product := initProduct(pool)
	productType := initProductType(pool)
	bot := initBot(product)

	go runBot(bot)
	go runREST()
	runGRPCServer(product, productType)
}

func initProduct(pool *pgxpool.Pool) productPkg.Interface {
	var product productPkg.Interface
	{
		product = productPkg.New(pool)
	}
	return product
}

func initProductType(pool *pgxpool.Pool) productTypePkg.Interface {
	var productType productTypePkg.Interface
	{
		productType = productTypePkg.New(pool)
	}
	return productType
}

func initBot(product productPkg.Interface) botPkg.Interface {
	var bot botPkg.Interface
	{
		bot = botPkg.MustNew()

		commandAdd := cmdAddPkg.New(product)
		bot.RegisterHandler(commandAdd)

		commandDelete := cmdDeletePkg.New(product)
		bot.RegisterHandler(commandDelete)

		commandEdit := cmdEditPkg.New(product)
		bot.RegisterHandler(commandEdit)

		commandList := cmdListPkg.New(product)
		bot.RegisterHandler(commandList)

		commandHelp := cmdHelpPkg.New(map[string]string{
			commandAdd.Name():    commandAdd.Description(),
			commandDelete.Name(): commandDelete.Description(),
			commandEdit.Name():   commandEdit.Description(),
			commandList.Name():   commandList.Description(),
		})
		bot.RegisterHandler(commandHelp)
	}

	return bot
}

func runGRPCServer(product productPkg.Interface, productType productTypePkg.Interface) {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProductServer(grpcServer, apiPkg.New(product))
	pb.RegisterProductTypeServer(grpcServer, apiPkg.NewProductType(productType))

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func runBot(bot botPkg.Interface) {
	if err := bot.Run(); err != nil {
		log.Panic(err)
	}
}

func clientInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	err := invoker(ctx, method, req, reply, cc, opts...)

	return err
}

func withClientUnaryInterceptor() grpc.DialOption {
	return grpc.WithUnaryInterceptor(clientInterceptor)
}

func runREST() {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	httpMux := http.NewServeMux()

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	swaggerDir := strings.Join([]string{path, "/bin"}, "")
	httpMux.HandleFunc("/openapiv2/", openAPIServer(swaggerDir))

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()), withClientUnaryInterceptor()}
	if err := pb.RegisterProductHandlerFromEndpoint(ctx, mux, ":8081", opts); err != nil {
		panic(err)
	}
	if err := pb.RegisterProductTypeHandlerFromEndpoint(ctx, mux, ":8081", opts); err != nil {
		panic(err)
	}

	httpMux.Handle("/", mux)

	s := &http.Server{
		Addr:    ":8080",
		Handler: httpMux,
	}

	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}

func openAPIServer(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
			glog.Errorf("Not Found: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		glog.Infof("Serving %s", r.URL.Path)
		p := strings.TrimPrefix(r.URL.Path, "/openapiv2/")
		p = path.Join(dir, p)
		http.ServeFile(w, r, p)
	}
}
