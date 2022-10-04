package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/maximgoltsov/botproject/pkg/api"
)

func NewProductTypeValidator(productTypeClient pb.ProductTypeClient) pb.ProductTypeServer {
	return &productTypeValidatorImplementation{
		productTypeClient: productTypeClient,
	}
}

type productTypeValidatorImplementation struct {
	pb.UnimplementedProductTypeServer
	productTypeClient pb.ProductTypeClient
}

func (i *productTypeValidatorImplementation) ProductTypeCreate(ctx context.Context, in *pb.ProductTypeCreateRequest) (*pb.ProductTypeCreateResponse, error) {
	//Validation...
	name := in.GetName()
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "field [name] should not be empty")
	}

	return i.productTypeClient.ProductTypeCreate(ctx, in)
}

func (i *productTypeValidatorImplementation) ProductTypeDelete(ctx context.Context, in *pb.ProductTypeDeleteRequest) (*pb.ProductTypeDeleteResponse, error) {
	//Validation...
	return i.productTypeClient.ProductTypeDelete(ctx, in)
}

func (i *productTypeValidatorImplementation) ProductTypeUpdate(ctx context.Context, in *pb.ProductTypeUpdateRequest) (*pb.ProductTypeUpdateResponse, error) {
	//Validation...
	name := in.GetName()
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "field [name] should not be empty")
	}

	return i.productTypeClient.ProductTypeUpdate(ctx, in)
}

func (i *productTypeValidatorImplementation) ProductTypeGet(ctx context.Context, in *pb.ProductTypeGetRequest) (*pb.ProductTypeGetResponse, error) {
	//Validation...
	return i.productTypeClient.ProductTypeGet(ctx, in)
}

func (i *productTypeValidatorImplementation) ProductTypeList(ctx context.Context, in *pb.ProductTypeListRequest) (*pb.ProductTypeListResponse, error) {
	//Validation...
	return i.productTypeClient.ProductTypeList(ctx, in)
}
