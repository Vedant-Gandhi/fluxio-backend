package schema

type Logger interface {
	Info(msg string, data ...interface{})
	Warn(msg string)
	Debug(msg string, data ...interface{})
	Error(msg string, err error)
	With(key string, value interface{}) LoggerChain
	WithField(fields map[string]interface{}) LoggerChain
}

type LoggerChain interface {
	Logger
}
