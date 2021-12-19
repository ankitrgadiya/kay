package kv

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
