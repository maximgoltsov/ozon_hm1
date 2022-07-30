package help

import (
	"fmt"

	"github.com/pkg/errors"

	commandPkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command"
)

var ErrBadArgument = errors.New("invalid arguments")

func New(extendedMap map[string]string) commandPkg.Interface {
	if extendedMap == nil {
		extendedMap = map[string]string{}
	}

	return &command{
		extendedMap: extendedMap,
	}
}

type command struct {
	extendedMap map[string]string
}

func (c *command) Name() string {
	return "help"
}

func (c *command) Description() string {
	return "- list commands"
}

func (c *command) Process(_ string) string {
	result := fmt.Sprintf("/%s - %s\n", c.Name(), c.Description())

	for cmd, description := range c.extendedMap {
		result += fmt.Sprintf("/%s - %s\n", cmd, description)
	}

	return result
}
