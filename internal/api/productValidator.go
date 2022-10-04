package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/maximgoltsov/botproject/pkg/api"
)

func NewProductValidator(productClient pb.ProductClient) pb.ProductServer {
	return &productValidatorImplementation{
		productClient: productClient,
	}
}

type productValidatorImplementation struct {
	pb.UnimplementedProductServer
	productClient pb.ProductClient
}

func (i *productValidatorImplementation) ProductCreate(ctx context.Context, in *pb.ProductCreateRequest) (*pb.ProductCreateResponse, error) {
	title := in.GetTitle()
	if title == "" {
		return nil, status.Error(codes.InvalidArgument, "field [title] should not be empty")
	}

	return i.productClient.ProductCreate(ctx, in)
}

func (i *productValidatorImplementation) ProductDelete(ctx context.Context, in *pb.ProductDeleteRequest) (*pb.ProductDeleteResponse, error) {
	//Validation...
	return i.productClient.ProductDelete(ctx, in)
}

func (i *productValidatorImplementation) ProductUpdate(ctx context.Context, in *pb.ProductUpdateRequest) (*pb.ProductUpdateResponse, error) {
	//Validation...
	title := in.GetTitle()
	if title == "" {
		return nil, status.Error(codes.InvalidArgument, "field [title] should not be empty")
	}

	return i.productClient.ProductUpdate(ctx, in)
}

func (i *productValidatorImplementation) ProductGet(ctx context.Context, in *pb.ProductGetRequest) (*pb.ProductGetResponse, error) {
	//Validation...
	return i.productClient.ProductGet(ctx, in)
}

func (i *productValidatorImplementation) ProductList(ctx context.Context, in *pb.ProductListRequest) (*pb.ProductListResponse, error) {
	//Validation...
	return i.productClient.ProductList(ctx, in)
}
