package command

type Interface interface {
	Name() string
	Description() string
	Process(args string) string
}
