package models

import "fmt"

type ProductType struct {
	Id   uint64
	Name string
}

func (pt *ProductType) String() string {
	return fmt.Sprintf("%v) %v", pt.Id, pt.Name)
}
