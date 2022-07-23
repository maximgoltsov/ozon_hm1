package handlers

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func parseAddArguments(data string) (string, uint64, error) {
	title, err := parseTitle(data, 0)
	if err != nil {
		return "", 0, err
	}

	price, err := parsePrice(data, 1)
	if err != nil {
		return "", 0, err
	}

	return title, price, nil
}

func parseDeleteArguments(data string) (uint64, error) {
	id, err := parseId(data, 0)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func parseEditArguments(data string) (uint64, string, uint64, error) {
	id, err := parseId(data, 0)
	if err != nil {
		return 0, "", 0, err
	}

	title, err := parseTitle(data, 1)
	if err != nil {
		return 0, "", 0, err
	}

	price, err := parsePrice(data, 2)
	if err != nil {
		return 0, "", 0, err
	}

	return id, title, price, nil
}

func parseId(data string, position int) (uint64, error) {
	params := strings.Split(data, " ")

	if len(params) < position {
		return 0, errors.Wrapf(BadArgument, "%d items: <%v>", len(params), params)
	}

	id, err := strconv.ParseUint(params[position], 10, 64)
	if err != nil {
		return 0, errors.Wrapf(BadArgument, "passed invalid id")
	}

	return id, nil
}

func parseTitle(data string, position int) (string, error) {
	params := strings.Split(data, " ")

	if len(params) < position {
		return "", errors.Wrapf(BadArgument, "%d items: <%v>", len(params), params)
	}

	return params[position], nil
}

func parsePrice(data string, position int) (uint64, error) {
	params := strings.Split(data, " ")

	if len(params) < position {
		return 0, errors.Wrapf(BadArgument, "%d items: <%v>", len(params), params)
	}

	price, err := strconv.ParseUint(params[position], 10, 64)
	if err != nil {
		return 0, errors.Wrapf(BadArgument, "passed invalid price")
	}

	return price, nil
}
