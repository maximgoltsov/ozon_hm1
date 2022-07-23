package storage

import "fmt"

var productLastId = uint64(0)

type Product struct {
	id    uint64
	title string
	price uint64
}

func NewProduct(title string, price uint64) (*Product, error) {
	product := Product{}

	if err := product.SetTitle(title); err != nil {
		return nil, err
	}

	if err := product.SetPrice(price); err != nil {
		return nil, err
	}

	if err := product.SetId(); err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *Product) SetId() error {
	productLastId++
	p.id = productLastId
	return nil
}

func (p *Product) SetTitle(title string) error {
	if len(title) == 0 {
		return fmt.Errorf("title should not be empty")
	}

	p.title = title
	return nil
}

func (p *Product) SetPrice(price uint64) error {
	if price == 0 {
		return fmt.Errorf("price should not be 0")
	}

	p.price = price
	return nil
}

func (p Product) String() string {
	return fmt.Sprintf("%d) %s (Price: %d)", p.id, p.title, p.price)
}

func (p Product) GetId() uint64 {
	return p.id
}

func (p Product) GetTitle() string {
	return p.title
}

func (p Product) GetPrice() uint64 {
	return p.price
}
