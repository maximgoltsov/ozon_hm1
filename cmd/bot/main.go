package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	apiPkg "github.com/maximgoltsov/botproject/internal/api"
	botPkg "github.com/maximgoltsov/botproject/internal/pkg/bot"

	cmdAddPkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command/add"
	cmdDeletePkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command/delete"
	cmdEditPkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command/edit"
	cmdHelpPkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command/help"
	cmdListPkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command/list"
	productPkg "github.com/maximgoltsov/botproject/internal/pkg/core/product"
	pb "github.com/maximgoltsov/botproject/pkg/api"
)

func main() {
	var product productPkg.Interface
	{
		product = productPkg.New()
	}

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

	go runBot(bot)
	go runREST()
	runGRPCServer(product)
}

func runGRPCServer(product productPkg.Interface) {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProductServer(grpcServer, apiPkg.New(product))

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func runBot(bot botPkg.Interface) {
	if err := bot.Run(); err != nil {
		log.Panic(err)
	}
}

func runREST() {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterProductHandlerFromEndpoint(ctx, mux, ":8081", opts); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}

// func headerMatcherREST(key string) (string, bool) {
// 	switch key {
// 	case "Custom":
// 		return key, true
// 	default:
// 		return key, false
// 	}
// }
