package edit

import (
	"context"
	"fmt"
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
	return "edit"
}

func (c *command) Description() string {
	return "<id> <title> <price> - edit title and price of product with passed id"
}

func (c *command) Process(ctx context.Context, args string) string {
	params := strings.Split(args, " ")[:]
	if len(params) != 3 {
		return ErrBadArgument.Error()
	}

	id, err := parserPkg.ParseId(params[0])
	if err != nil {
		return err.Error()
	}

	title, price, err := parserPkg.ParseTitleAndPrice(params[1:3])
	if err != nil {
		return err.Error()
	}

	_, err = c.product.UpsertProduct(ctx, models.Product{
		Id:    id,
		Title: title,
		Price: price,
	})

	if err != nil {
		if errors.Is(err, productPkg.ErrValidation) {
			return ErrBadArgument.Error()
		}
		return "internal error"
	}

	return fmt.Sprintln("product was edited successfully")
}
