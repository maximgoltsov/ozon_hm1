package models

import "fmt"

type Product struct {
	Id      uint64
	Title   string
	Price   uint64
	Type_Id uint64
}

func (p *Product) String() string {
	return fmt.Sprintf("%v) %v (Price: %v)", p.Id, p.Title, p.Price)
}
