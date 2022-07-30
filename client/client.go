package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/maximgoltsov/botproject/pkg/api"
)

func main() {
	conns, err := grpc.Dial(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := pb.NewProductClient(conns)

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "custom", "hello")

	response, err := client.ProductList(ctx, &pb.ProductListRequest{})
	if err != nil {
		panic(err)
	}

	log.Printf("\nresponse: [%v]\n", response)
}
