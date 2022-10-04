package delete

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	commandPkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command"
	productPkg "github.com/maximgoltsov/botproject/internal/pkg/core/product"
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
	return "delete"
}

func (c *command) Description() string {
	return "<id> - delete product with passed id"
}

func (c *command) Process(ctx context.Context, args string) string {
	params := strings.Split(args, " ")
	if len(params) != 1 {
		return ErrBadArgument.Error()
	}

	id, err := parserPkg.ParseId(params[0])
	if err != nil {
		return err.Error()
	}

	if err := c.product.DeleteProductById(ctx, id); err != nil {
		return err.Error()
	}

	return fmt.Sprintf("product with id %v was deleted", id)
}
