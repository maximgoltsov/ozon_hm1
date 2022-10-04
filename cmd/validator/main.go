package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiPkg "github.com/maximgoltsov/botproject/internal/api"
	"github.com/maximgoltsov/botproject/internal/kafka"
	"github.com/maximgoltsov/botproject/internal/pkg/logger"
	pb "github.com/maximgoltsov/botproject/pkg/api"
)

var (
	serverAddr = flag.String("addr", "localhost:8081", "The server address in the format of host:port")
)

func main() {
	logger.InitLogger()

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Logger.Fatal("can't connect to bot service", err)
	}
	defer conn.Close()

	productClient := pb.NewProductClient(conn)
	productTypeClient := pb.NewProductTypeClient(conn)

	go runKafka()
	go validatorGRPCServer(productClient, productTypeClient)
	runValidatorREST()
}

func validatorGRPCServer(productClient pb.ProductClient, productTypeClient pb.ProductTypeClient) {
	listener, err := net.Listen("tcp", ":8083")
	if err != nil {
		logger.Logger.Fatal("can't start gRPC server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProductServer(grpcServer, apiPkg.NewProductValidator(productClient))
	pb.RegisterProductTypeServer(grpcServer, apiPkg.NewProductTypeValidator(productTypeClient))

	if err = grpcServer.Serve(listener); err != nil {
		logger.Logger.Error("can't start gRPC server", err)
	}
}

func runValidatorREST() {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	httpMux := http.NewServeMux()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterProductHandlerFromEndpoint(ctx, mux, ":8083", opts); err != nil {
		panic(err)
	}

	httpMux.Handle("/", mux)

	s := &http.Server{
		Addr:    ":8084",
		Handler: httpMux,
	}

	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}

func runKafka() {
	brokers := []string{"localhost:9095"}

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true

	asyncProducer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		logger.Logger.Fatalf("error with async kafka producer: %v", err)
	}

	go func() {
		for msg := range asyncProducer.Errors() {
			logger.Logger.Errorf("error: %v", msg)
		}
	}()

	go func() {
		for msg := range asyncProducer.Successes() {
			logger.Logger.Infof("success: %v", msg)
		}
	}()

	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			for {
				data, err := json.Marshal(&kafka.Message{
					Service: "Validator",
					Action:  "Random Text " + time.Now().String(),
				})
				if err != nil {
					logger.Logger.Error("cant convert message", err)
					return
				}
				asyncProducer.Input() <- &sarama.ProducerMessage{
					Topic: "validator",
					Key:   sarama.StringEncoder(key),
					Value: sarama.ByteEncoder(data),
				}
				time.Sleep(time.Second * 5)
			}
		}(fmt.Sprintf("%v", i))
	}

	logger.Logger.Info("kafka producer started")
	wg.Wait()
}
