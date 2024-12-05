package tui

type View interface {
	Hydrate(data ...interface{}) error
}
