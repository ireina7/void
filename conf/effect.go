package conf

import "fmt"

type Effect interface {
	Version() string
	Host() string
	Port() string
}

func Addr(conf Effect) string {
	return fmt.Sprintf("%s:%s", conf.Host(), conf.Port())
}
