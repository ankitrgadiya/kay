package kv

// Driver defines an interface for initializing new KeyValue implementation. The
// databases will have corresponding Driver implementations.
//
// To use a Driver with this package, it must first be registered through the
// Register function.
type Driver interface {
	New() KeyValue
}

// DriverFunc defines a function type that implements the Driver interface.
type DriverFunc func() KeyValue

func (df DriverFunc) New() KeyValue { return df() }
