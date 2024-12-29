package transport

//go:generate mockgen -package mocks -source=interface.go -destination=mocks/mocks.go *

type IConnWrapper interface {
	WriteMessage(msg string) error
	ReadMessage() (string, error)
	RemoteAddr() string
	Close() error
}