package local

import (
	"strconv"
	"sync"

	"github.com/pkg/errors"

	cachePkg "github.com/maximgoltsov/botproject/internal/pkg/core/product/cache"
	"github.com/maximgoltsov/botproject/internal/pkg/core/product/models"
)

var lastId = uint64(0)

const poolSize = 10

var (
	ErrProductNotExists = errors.New("product does not exist")
	ErrProductExists    = errors.New("product exist")
)

func New() cachePkg.Interface {
	return &cache{
		mu:     sync.RWMutex{},
		data:   map[uint64]models.Product{},
		poolCh: make(chan struct{}, poolSize),
	}
}

type cache struct {
	mu     sync.RWMutex
	data   map[uint64]models.Product
	poolCh chan struct{}
}

func (c *cache) GetProducts() []models.Product {
	c.poolCh <- struct{}{}
	c.mu.RLock()
	defer func() {
		c.mu.RUnlock()
		<-c.poolCh
	}()

	result := make([]models.Product, 0, len(c.data))

	for idx := range c.data {
		result = append(result, c.data[idx])
	}

	return result
}

func (c *cache) GetProduct(id uint64) (models.Product, error) {
	c.poolCh <- struct{}{}
	c.mu.RLock()
	defer func() {
		<-c.poolCh
		c.mu.RUnlock()
	}()

	product, ok := c.data[id]

	if !ok {
		return models.Product{}, errors.Wrap(ErrProductExists, strconv.FormatUint(id, 10))
	}

	return product, nil
}

func (c *cache) UpsertProduct(p models.Product) error {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		<-c.poolCh
		c.mu.Unlock()
	}()

	if p.Id != 0 {
		return updateProduct(p, c)
	} else {
		lastId++
		p.Id = lastId
		return addProduct(p, c)
	}
}

func addProduct(p models.Product, c *cache) error {
	if _, ok := c.data[p.Id]; ok {
		return errors.Wrap(ErrProductExists, strconv.FormatUint(p.Id, 10))
	}

	c.data[p.Id] = p
	return nil
}

func updateProduct(p models.Product, c *cache) error {
	if _, ok := c.data[p.Id]; !ok {
		return errors.Wrap(ErrProductNotExists, strconv.FormatUint(p.Id, 10))
	}

	c.data[p.Id] = p
	return nil
}

func (c *cache) DeleteProduct(p models.Product) error {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		<-c.poolCh
		c.mu.Unlock()
	}()

	if _, ok := c.data[p.Id]; !ok {
		return errors.Wrap(ErrProductNotExists, strconv.FormatUint(p.Id, 10))
	}

	delete(c.data, p.Id)
	return nil
}

func (c *cache) DeleteProductById(id uint64) error {
	c.poolCh <- struct{}{}
	c.mu.Lock()
	defer func() {
		<-c.poolCh
		c.mu.Unlock()
	}()

	if _, ok := c.data[id]; !ok {
		return errors.Wrap(ErrProductNotExists, strconv.FormatUint(id, 10))
	}

	delete(c.data, id)
	return nil
}
