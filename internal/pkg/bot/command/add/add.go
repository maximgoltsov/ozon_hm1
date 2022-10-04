package add

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	commandPkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command"
	productPkg "github.com/maximgoltsov/botproject/internal/pkg/core/product"
	"github.com/maximgoltsov/botproject/internal/pkg/core/product/models"
	parserPkg "github.com/maximgoltsov/botproject/internal/pkg/parser"
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
	return "add"
}

func (c *command) Description() string {
	return "<title> <price> - add new product with title and price"
}

func (c *command) Process(ctx context.Context, args string) string {
	params := strings.Split(args, " ")
	if len(params) != 2 {
		return ErrBadArgument.Error()
	}

	title, price, err := parserPkg.ParseTitleAndPrice(params)
	if err != nil {
		return err.Error()
	}

	_, err = c.product.UpsertProduct(ctx, models.Product{
		Title: title,
		Price: price,
	})

	if err != nil {
		if errors.Is(err, productPkg.ErrValidation) {
			return ErrBadArgument.Error()
		}
		return "internal error"
	}

	return "success"
}
