package kv

import (
	"github.com/pkg/errors"

	"argc.in/kay/config"
)

func Open(name string, s config.Section) (KeyValue, error) {
	keyvalue, err := openConn(name, s)
	if err != nil {
		return nil, err
	}

	if initializer, ok := keyvalue.(Initializer); ok {
		if err := initializer.Init(); err != nil {
			return nil, err
		}
	}

	return keyvalue, nil
}

func openConn(name string, s config.Section) (KeyValue, error) {
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
