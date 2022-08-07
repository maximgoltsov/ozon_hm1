package list

import (
	"sort"
	"strings"

	"github.com/pkg/errors"

	commandPkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command"
	productPkg "github.com/maximgoltsov/botproject/internal/pkg/core/product"
)

var ErrBadArgument = errors.New("invalid arguments")

func New(product productPkg.Interface) commandPkg.Interface {
	return &command{
		product: product,
	}
}

type command struct {
	product productPkg.Interface
}

func (c *command) Name() string {
	return "list"
}

func (c *command) Description() string {
	return " - list products"
}

func (c *command) Process(args string) string {

	products := c.product.GetProducts(0, 0, false)

	sort.Slice(products, func(i, j int) bool {
		return products[i].Id < products[j].Id
	})

	res := make([]string, 0, len(products))
	for idx := range products {
		res = append(res, products[idx].String())
	}

	if len(res) == 0 {
		return "List is empty"
	}

	return strings.Join(res, "\n")
}
