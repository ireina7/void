package conf

type Effect interface {
	Version() string
	Host() string
	Port() string
}
