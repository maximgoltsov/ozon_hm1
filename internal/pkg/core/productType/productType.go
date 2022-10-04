package product_type

import (
	"context"

	"github.com/pkg/errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/maximgoltsov/botproject/internal/pkg/core/productType/models"

	repositoryPkg "github.com/maximgoltsov/botproject/internal/pkg/repository"
	dbPkg "github.com/maximgoltsov/botproject/internal/pkg/repository/db"
)

var ErrValidation = errors.New("invalid data")

type Interface interface {
	UpsertProductType(ctx context.Context, pt models.ProductType) (uint64, error)
	DeleteProductTypeById(ctx context.Context, id uint64) error
	GetProductType(ctx context.Context, id uint64) (models.ProductType, error)
	GetProductTypes(ctx context.Context, limit uint64, offset uint64, desc bool) []models.ProductType
}

func New(pool *pgxpool.Pool) Interface {
	return &core{
		pool:       pool,
		repository: dbPkg.NewProductTypeRepository(pool),
	}
}

type core struct {
	pool       *pgxpool.Pool
	repository repositoryPkg.ProductType
}

func (c *core) UpsertProductType(ctx context.Context, pt models.ProductType) (uint64, error) {
	if pt.Name == "" {
		return 0, errors.Wrap(ErrValidation, "filed: [name] cannot be empty")
	}

	return c.repository.UpsertProductType(ctx, pt)
}
func (c *core) DeleteProductTypeById(ctx context.Context, id uint64) error {
	return c.repository.DeleteProductTypeById(ctx, id)
}
func (c *core) GetProductType(ctx context.Context, id uint64) (models.ProductType, error) {
	return c.repository.GetProductType(ctx, id)
}
func (c *core) GetProductTypes(ctx context.Context, limit uint64, offset uint64, desc bool) []models.ProductType {
	return c.repository.GetProductTypes(ctx, limit, offset, desc)
}
