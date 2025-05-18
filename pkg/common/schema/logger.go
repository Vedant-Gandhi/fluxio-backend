package schema

type Logger interface {
	Info(msg string, data ...interface{})
	Warn(msg string)
	Debug(msg string, data ...interface{})
	Error(msg string, err error)
}
