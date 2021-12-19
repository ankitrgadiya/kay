package kv

import (
	"io"

	"github.com/pkg/errors"

	"argc.in/kay/config"
)

// Open opens a new database connection based on the configuration specified in
// the Section.
//
// The Section must specify a Driver that is registered in this package to
// successfully open the connection.
func Open(s config.Section) (KeyValue, io.Closer, error) {
	keyvalue, err := openConn(s)
	if err != nil {
		return nil, nil, err
	}

	if initializer, ok := keyvalue.(Initializer); ok {
		if err := initializer.Init(); err != nil {
			return nil, nil, err
		}
	}

	return keyvalue, getCloser(keyvalue), nil
}

func openConn(s config.Section) (KeyValue, error) {
	driver := getDriver(s.DriverName())
	if driver == nil {
		return nil, ErrDriverNotFound
	}

	keyvalue := driver.New()

	if err := s.Unmarshal(keyvalue); err != nil {
		return nil, errors.Wrap(err, "loading driver config")
	}

	return keyvalue, nil
}
