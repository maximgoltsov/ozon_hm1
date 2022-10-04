package parser

import (
	"strconv"

	"github.com/pkg/errors"
)

var BadArgument = errors.New("bad arguments")
var ErrInvalidIdParser = errors.Wrapf(BadArgument, "passed invalid id")

func ParseId(data string) (uint64, error) {
	id, err := strconv.ParseUint(data, 10, 64)
	if err != nil {
		return 0, ErrInvalidIdParser
	}

	return id, nil
}

func ParseTitleAndPrice(data []string) (string, uint64, error) {
	title := data[0]

	price, err := strconv.ParseUint(data[1], 10, 64)
	if err != nil {
		return "", 0, errors.Wrapf(BadArgument, "passed invalid price")
	}

	return title, price, nil
}
