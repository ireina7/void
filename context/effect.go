package context

type Effect interface {
	Error() error
	// FoundError(err error)
}
