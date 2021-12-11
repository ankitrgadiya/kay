package kv

type Initializer interface {
	Init() error
}

type KeyValue interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

type Deleter interface {
	Delete(key string) error
}
