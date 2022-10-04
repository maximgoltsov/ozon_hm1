package product

import (
	"context"

	"github.com/pkg/errors"

	models "github.com/maximgoltsov/botproject/internal/pkg/core/product/models"
	repositoryPkg "github.com/maximgoltsov/botproject/internal/pkg/repository"
)

var ErrValidation = errors.New("invalid data")

type Interface interface {
	UpsertProduct(ctx context.Context, product models.Product) (uint64, error)
	DeleteProductById(ctx context.Context, id uint64) error
	GetProduct(ctx context.Context, id uint64) (models.Product, error)
	GetProducts(ctx context.Context, limit uint64, offset uint64, desc bool) []models.Product
}

func New(repository repositoryPkg.Product) Interface {
	return &core{
		repository: repository,
	}
}

type core struct {
	repository repositoryPkg.Product
}

func (c *core) UpsertProduct(ctx context.Context, product models.Product) (uint64, error) {
	if product.Title == "" {
		return 0, errors.Wrap(ErrValidation, "field: [title] cannot be empty")
	}

	return c.repository.UpsertProduct(ctx, product)
}

func (c *core) DeleteProductById(ctx context.Context, id uint64) error {
	return c.repository.DeleteProductById(ctx, id)
}

func (c *core) GetProduct(ctx context.Context, id uint64) (models.Product, error) {
	return c.repository.GetProduct(ctx, id)
}

func (c *core) GetProducts(ctx context.Context, limit uint64, offset uint64, desc bool) []models.Product {
	return c.repository.GetProducts(ctx, limit, offset, desc)
}
