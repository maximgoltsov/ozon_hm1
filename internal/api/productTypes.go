package api

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	productTypePkg "github.com/maximgoltsov/botproject/internal/pkg/core/productType"
	"github.com/maximgoltsov/botproject/internal/pkg/core/productType/models"
	pb "github.com/maximgoltsov/botproject/pkg/api"
)

func NewProductType(productType productTypePkg.Interface) pb.ProductTypeServer {
	return &productTypeImplementation{
		productType: productType,
	}
}

type productTypeImplementation struct {
	pb.UnimplementedProductTypeServer
	productType productTypePkg.Interface
}

func (i *productTypeImplementation) ProductTypeCreate(ctx context.Context, in *pb.ProductTypeCreateRequest) (*pb.ProductTypeCreateResponse, error) {
	id, err := i.productType.UpsertProductType(ctx, models.ProductType{
		Id:   0,
		Name: in.GetName(),
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductTypeCreateResponse{Id: id}, nil
}

func (i *productTypeImplementation) ProductTypeDelete(ctx context.Context, in *pb.ProductTypeDeleteRequest) (*pb.ProductTypeDeleteResponse, error) {
	if err := i.productType.DeleteProductTypeById(ctx, in.GetId()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductTypeDeleteResponse{}, nil
}

func (i *productTypeImplementation) ProductTypeUpdate(ctx context.Context, in *pb.ProductTypeUpdateRequest) (*pb.ProductTypeUpdateResponse, error) {
	id, err := i.productType.UpsertProductType(ctx, models.ProductType{
		Id:   in.GetId(),
		Name: in.GetName(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductTypeUpdateResponse{Id: id}, nil
}

func (i *productTypeImplementation) ProductTypeGet(ctx context.Context, in *pb.ProductTypeGetRequest) (*pb.ProductTypeGetResponse, error) {
	productType, err := i.productType.GetProductType(ctx, in.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductTypeGetResponse{
		Id:   productType.Id,
		Name: productType.Name,
	}, nil
}

func (i *productTypeImplementation) ProductTypeList(ctx context.Context, in *pb.ProductTypeListRequest) (*pb.ProductTypeListResponse, error) {
	productTypes := i.productType.GetProductTypes(ctx, in.GetLimit(), in.GetOffset(), in.GetDesc())

	result := make([]*pb.ProductTypeListResponse_ProductType, 0, len(productTypes))

	for idx := range productTypes {
		result = append(result, &pb.ProductTypeListResponse_ProductType{
			Id:   productTypes[idx].Id,
			Name: productTypes[idx].Name,
		})
	}

	return &pb.ProductTypeListResponse{
		ProductTypes: result,
	}, nil
}
