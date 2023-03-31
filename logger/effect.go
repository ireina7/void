package logger

type Effect interface {
	Info(msg string)
	Error(err error)
	Debug(s string)
	Fatal(err error)
}
