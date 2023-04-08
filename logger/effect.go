package logger

type Effect interface {
	Info(msg string)
	LogError(err error)
	Debug(s string)
	Fatal(err error)
}
