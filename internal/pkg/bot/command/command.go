package command

import "context"

type Interface interface {
	Name() string
	Description() string
	Process(ctx context.Context, args string) string
}
