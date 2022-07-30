package cache

import "github.com/maximgoltsov/botproject/internal/pkg/core/product/models"

type Interface interface {
	UpsertProduct(product models.Product) error
	DeleteProductById(id uint64) error
	GetProduct(id uint64) (models.Product, error)
	GetProducts() []models.Product
}
