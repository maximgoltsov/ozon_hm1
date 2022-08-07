package product

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	cachePkg "github.com/maximgoltsov/botproject/internal/pkg/core/product/cache"
	postgresPkg "github.com/maximgoltsov/botproject/internal/pkg/core/product/cache/postgres"
	"github.com/maximgoltsov/botproject/internal/pkg/core/product/models"
)

var ErrValidation = errors.New("invalid data")

type Interface interface {
	UpsertProduct(product models.Product) (uint64, error)
	DeleteProductById(id uint64) error
	GetProduct(id uint64) (models.Product, error)
	GetProducts(limit uint64, offset uint64, desc bool) []models.Product
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

func (c *core) UpsertProduct(product models.Product) (uint64, error) {
	if product.Title == "" {
		return 0, errors.Wrap(ErrValidation, "field: [title] cannot be empty")
	}

	return c.cache.UpsertProduct(product)
}

func (c *core) DeleteProductById(id uint64) error {
	return c.cache.DeleteProductById(id)
}

func (c *core) GetProduct(id uint64) (models.Product, error) {
	return c.cache.GetProduct(id)
}

func (c *core) GetProducts(limit uint64, offset uint64, desc bool) []models.Product {
	return c.cache.GetProducts(limit, offset, desc)
}
