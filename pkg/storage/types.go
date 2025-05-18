package storage

type Operations[T any] interface {
	Store(k string, v T) error
	Get(k string) (T, error)
	Delete(k string) error
	Close() error
}
