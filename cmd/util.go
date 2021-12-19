package cmd

import (
	"io"

	"argc.in/kay/kv"
)

func openDatabase(name string) (kv.KeyValue, io.Closer, error) {
	s, err := conf.Section(name)
	if err != nil {
		return nil, nil, err
	}

	keyvalue, closer, err := kv.Open(s)
	if err != nil {
		return nil, nil, err
	}

	return keyvalue, closer, nil
}
