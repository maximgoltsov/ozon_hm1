package repository

import (
	"context"

	product "github.com/maximgoltsov/botproject/internal/pkg/core/product/models"
	productType "github.com/maximgoltsov/botproject/internal/pkg/core/productType/models"
)

type ProductType interface {
	UpsertProductType(ctx context.Context, pt productType.ProductType) (uint64, error)
	DeleteProductTypeById(ctx context.Context, id uint64) error
	GetProductType(ctx context.Context, id uint64) (productType.ProductType, error)
	GetProductTypes(ctx context.Context, limit uint64, offset uint64, desc bool) []productType.ProductType
}

type Product interface {
	UpsertProduct(ctx context.Context, product product.Product) (uint64, error)
	DeleteProductById(ctx context.Context, id uint64) error
	GetProduct(ctx context.Context, id uint64) (product.Product, error)
	GetProducts(ctx context.Context, limit uint64, offset uint64, desc bool) []product.Product
}
