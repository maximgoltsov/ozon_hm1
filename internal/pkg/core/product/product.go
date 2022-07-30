package product

import (
	"github.com/pkg/errors"

	cachePkg "github.com/maximgoltsov/botproject/internal/pkg/core/product/cache"
	localCachePkg "github.com/maximgoltsov/botproject/internal/pkg/core/product/cache/local"
	"github.com/maximgoltsov/botproject/internal/pkg/core/product/models"
)

var ErrValidation = errors.New("invalid data")

type Interface interface {
	UpsertProduct(product models.Product) error
	DeleteProductById(id uint64) error
	GetProduct(id uint64) (models.Product, error)
	GetProducts() []models.Product
}

func New() Interface {
	return &core{
		cache: localCachePkg.New(),
	}
}

type core struct {
	cache cachePkg.Interface
}

func (c *core) UpsertProduct(product models.Product) error {
	if product.Title == "" {
		return errors.Wrap(ErrValidation, "field: [title] cannot be empty")
	}

	return c.cache.UpsertProduct(product)
}

func (c *core) DeleteProductById(id uint64) error {
	return c.cache.DeleteProductById(id)
}

func (c *core) GetProduct(id uint64) (models.Product, error) {
	return c.cache.GetProduct(id)
}

func (c *core) GetProducts() []models.Product {
	return c.cache.GetProducts()
}
