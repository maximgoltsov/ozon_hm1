package product_type

import (
	"context"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/maximgoltsov/botproject/internal/pkg/core/productType/models"

	cachePkg "github.com/maximgoltsov/botproject/internal/pkg/core/productType/cache"
	postgresPkg "github.com/maximgoltsov/botproject/internal/pkg/core/productType/cache/postgres"
)

var ErrValidation = errors.New("invalid data")

type Interface interface {
	UpsertProductType(ctx context.Context, pt models.ProductType) (uint64, error)
	DeleteProductTypeById(ctx context.Context, id uint64) error
	GetProductType(ctx context.Context, id uint64) (models.ProductType, error)
	GetProductTypes(ctx context.Context) []models.ProductType
}

func New(pool *pgxpool.Pool) Interface {
	return &core{
		pool:  pool,
		cache: postgresPkg.New(pool),
	}
}

type core struct {
	pool  *pgxpool.Pool
	cache cachePkg.Interface
}

func (c *core) UpsertProductType(ctx context.Context, pt models.ProductType) (uint64, error) {
	if pt.Name == "" {
		return 0, errors.Wrap(ErrValidation, "filed: [name] cannot be empty")
	}

	return c.cache.UpsertProductType(ctx, pt)
}
func (c *core) DeleteProductTypeById(ctx context.Context, id uint64) error {
	return c.cache.DeleteProductTypeById(ctx, id)
}
func (c *core) GetProductType(ctx context.Context, id uint64) (models.ProductType, error) {
	return c.cache.GetProductType(ctx, id)
}
func (c *core) GetProductTypes(ctx context.Context) []models.ProductType {
	return c.cache.GetProductTypes(ctx)
}
