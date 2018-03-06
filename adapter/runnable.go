package adapter

type Listener interface {
	Listen() error
	Stop() error
}
