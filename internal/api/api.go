package api

import (
	"context"

	errorsPkg "github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	productPkg "github.com/maximgoltsov/botproject/internal/pkg/core/product"
	"github.com/maximgoltsov/botproject/internal/pkg/core/product/models"
	pb "github.com/maximgoltsov/botproject/pkg/api"
)

func New(product productPkg.Interface) pb.ProductServer {
	return &implementation{
		product: product,
	}
}

type implementation struct {
	pb.UnimplementedProductServer
	product productPkg.Interface
}

func (i *implementation) ProductCreate(ctx context.Context, in *pb.ProductCreateRequest) (*pb.ProductCreateResponse, error) {
	id, err := i.product.UpsertProduct(models.Product{
		Id:      0,
		Title:   in.GetTitle(),
		Price:   uint64(in.GetPrice()),
		Type_Id: in.GetTypeId(),
	})
	if err != nil {
		if errorsPkg.Is(err, productPkg.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductCreateResponse{Id: id}, nil
}

func (i *implementation) ProductDelete(ctx context.Context, in *pb.ProductDeleteRequest) (*pb.ProductDeleteResponse, error) {
	if err := i.product.DeleteProductById(in.GetId()); err != nil {
		if errorsPkg.Is(err, productPkg.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductDeleteResponse{}, nil
}

func (i *implementation) ProductUpdate(ctx context.Context, in *pb.ProductUpdateRequest) (*pb.ProductUpdateResponse, error) {
	id, err := i.product.UpsertProduct(models.Product{
		Id:      uint64(in.GetId()),
		Title:   in.GetTitle(),
		Price:   uint64(in.GetPrice()),
		Type_Id: in.GetTypeId(),
	})
	if err != nil {
		if errorsPkg.Is(err, productPkg.ErrValidation) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductUpdateResponse{Id: id}, nil
}

func (i *implementation) ProductGet(ctx context.Context, in *pb.ProductGetRequest) (*pb.ProductGetResponse, error) {
	product, err := i.product.GetProduct(in.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProductGetResponse{
		Id:     product.Id,
		Title:  product.Title,
		Price:  product.Price,
		TypeId: product.Type_Id,
	}, nil
}

func (i *implementation) ProductList(ctx context.Context, in *pb.ProductListRequest) (*pb.ProductListResponse, error) {
	products := i.product.GetProducts(in.GetLimit(), in.GetOffset(), in.GetDesc())

	result := make([]*pb.ProductListResponse_Product, 0, len(products))

	for idx := range products {
		result = append(result, &pb.ProductListResponse_Product{
			Id:     products[idx].Id,
			Title:  products[idx].Title,
			Price:  products[idx].Price,
			TypeId: products[idx].Type_Id,
		})
	}

	return &pb.ProductListResponse{
		Products: result,
	}, nil
}
