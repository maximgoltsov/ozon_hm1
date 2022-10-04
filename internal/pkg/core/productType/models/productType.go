package models

import "fmt"

type ProductType struct {
	Id   uint64 `db:"id"`
	Name string `db:"name"`
}

func (pt *ProductType) String() string {
	return fmt.Sprintf("%v) %v", pt.Id, pt.Name)
}
