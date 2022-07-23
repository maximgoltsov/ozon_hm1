package handlers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/maximgoltsov/botproject/internal/commander"
	"github.com/maximgoltsov/botproject/internal/storage"
)

const (
	helpCmd   = "help"
	listCmd   = "list"
	addCmd    = "add"
	deleteCmd = "delete"
	editCmd   = "edit"
)

var BadArgument = errors.New("bad arguments")

func AddHandlers(c *commander.Commander) {
	c.RegisterHandler(helpCmd, helpHandler)
	c.RegisterHandler(listCmd, getProductsHandler)
	c.RegisterHandler(addCmd, addProductHandler)
	c.RegisterHandler(deleteCmd, deleteProductHandler)
	c.RegisterHandler(editCmd, editProductHandler)
}

func helpHandler(s string) string {
	return strings.Join([]string{
		"/help - list commands",
		"/list - list products",
		"/add <title> <price> - add new product with title and price",
		"/delete <id> - delete product with passed id",
		"/edit <id> <title> <price> - edit title and price of product with passed id",
	}, "\n")
}

func editProductHandler(data string) string {
	id, title, price, err := parseEditArguments(data)
	if err != nil {
		return err.Error()
	}

	product, err := storage.GetProduct(id)
	if err != nil {
		return err.Error()
	}

	if err := product.SetTitle(title); err != nil {
		return err.Error()
	}

	if err := product.SetPrice(price); err != nil {
		return err.Error()
	}

	return fmt.Sprintln("product was edited successfully")
}

func getProductsHandler(s string) string {
	products := storage.GetProducts()

	sort.Slice(products, func(i, j int) bool {
		return products[i].GetId() < products[j].GetId()
	})

	res := make([]string, 0, len(products))
	for idx := range products {
		res = append(res, products[idx].String())
	}

	return strings.Join(res, "\n")
}

func addProductHandler(data string) string {
	title, price, err := parseAddArguments(data)
	if err != nil {
		return err.Error()
	}

	p := storage.Product{}

	if err := p.SetTitle(title); err != nil {
		return err.Error()
	}

	if err := p.SetPrice(price); err != nil {
		return err.Error()
	}

	if err := storage.UpsertProduct(&p); err != nil {
		return err.Error()
	}

	return fmt.Sprintln("product was added")
}

func deleteProductHandler(data string) string {
	id, err := parseDeleteArguments(data)
	if err != nil {
		return err.Error()
	}

	if err := storage.DeleteProductById(id); err != nil {
		return err.Error()
	}

	return fmt.Sprintf("product with id %v was deleted", id)
}
