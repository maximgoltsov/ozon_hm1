package cache

import (
	"context"

	"github.com/maximgoltsov/botproject/internal/pkg/core/productType/models"
)

type Interface interface {
	UpsertProductType(ctx context.Context, pt models.ProductType) (uint64, error)
	DeleteProductTypeById(ctx context.Context, id uint64) error
	GetProductType(ctx context.Context, id uint64) (models.ProductType, error)
	GetProductTypes(ctx context.Context) []models.ProductType
}
