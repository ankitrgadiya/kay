package kv

import (
	"context"
	"io"
)

// KeyValue is the most basic interface that each Database MUST implement.
type KeyValue interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

// Initializer defines the interface that a KeyValue implementation can
// implement if explicit initialization is required.
//
// This method will be called after the configuration is marshalled into the
// implementation struct.
type Initializer interface {
	Init() error
}

// Deleter defines an extension interface for KeyValue implementations that
// defines the Delete operation on the Database.
//
// This interface is used in the "del" command.
type Deleter interface {
	KeyValue
	Delete(key string) error
}

// Event defines the struct used for Watcher to receive Events from the Driver.
type Event struct {
	Key   string
	Value []byte
}

// Watcher defines an extension interface for KeyValue implementations that can
// watch for changes in Keys.
//
// This interface is used in the "watch" command.
type Watcher interface {
	KeyValue
	Watch(ctx context.Context, key string) <-chan Event
}

// Lister defines an extension inteface for KeyValue implementations that can
// list the keys in the Database.
type Lister interface {
	KeyValue
	List(prefix string) (Iterator, error)
}

// Searcher defines an extension inteface for KeyValue implementations that can
// search for keys in the Database.
type Searcher interface {
	KeyValue
	Search(term string) (Iterator, error)
}

// Iterator defines the interface for iterating over the key-values of the
// Database.
type Iterator interface {
	io.Closer
	Next() (string, []byte, bool)
}
