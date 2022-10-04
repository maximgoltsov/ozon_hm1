package models

import "fmt"

type Product struct {
	Id      uint64 `db:"id"`
	Title   string `db:"title"`
	Price   uint64 `db:"price"`
	Type_Id uint64 `db:"type_id"`
}

func (p *Product) String() string {
	return fmt.Sprintf("%v) %v (Price: %v)", p.Id, p.Title, p.Price)
}
