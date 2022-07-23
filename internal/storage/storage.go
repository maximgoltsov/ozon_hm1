package storage

import (
	"strconv"

	"github.com/pkg/errors"
)

var data map[uint64]*Product

var ProductNotExists = errors.New("product does not exist")
var ProductExists = errors.New("product exist")

// Выполняется сама один раз
func init() {
	data = make(map[uint64]*Product)
	product, _ := NewProduct("Букет", 100)

	if err := addProduct(product); err != nil {
		panic(err)
	}
}

func GetProducts() []*Product {
	res := make([]*Product, 0, len(data))

	for idx := range data {
		res = append(res, data[idx])
	}

	return res
}

func GetProduct(id uint64) (*Product, error) {
	product, ok := data[id]

	if !ok {
		return nil, errors.Wrap(ProductExists, strconv.FormatUint(id, 10))
	}

	return product, nil
}

func UpsertProduct(p *Product) error {
	if p.id != 0 {
		updateProduct(p)
	} else {
		p.SetId()
		addProduct(p)
	}

	return nil
}

func addProduct(p *Product) error {
	if _, ok := data[p.GetId()]; ok {
		return errors.Wrap(ProductExists, strconv.FormatUint(p.GetId(), 10))
	}

	data[p.GetId()] = p
	return nil
}

func updateProduct(p *Product) error {
	if _, ok := data[p.GetId()]; !ok {
		return errors.Wrap(ProductNotExists, strconv.FormatUint(p.GetId(), 10))
	}

	data[p.GetId()] = p
	return nil
}

func DeleteProduct(p *Product) error {
	if _, ok := data[p.GetId()]; !ok {
		return errors.Wrap(ProductNotExists, strconv.FormatUint(p.GetId(), 10))
	}

	delete(data, p.GetId())
	return nil
}

func DeleteProductById(id uint64) error {
	if _, ok := data[id]; !ok {
		return errors.Wrap(ProductNotExists, strconv.FormatUint(id, 10))
	}

	delete(data, id)
	return nil
}
