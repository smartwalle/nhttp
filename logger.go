package http4go

type Logger interface {
	SetPrefix(prefix string)
	Prefix() string
	Println(args ...interface{})
	Printf(format string, args ...interface{})
	Output(calldepth int, s string) error
}
