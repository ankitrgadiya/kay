package pebble

import (
	"github.com/cockroachdb/pebble"
	"github.com/pkg/errors"
)

type iter struct {
	cur    *pebble.Iterator
	seeked bool
}

func (i *iter) Next() (string, []byte, bool) {
	valid := i.cur.Valid()
	if !valid {
		return "", nil, false
	}
	defer i.cur.Next()

	key, value := i.cur.Key(), i.cur.Value()
	return string(key), value, true
}

func (i *iter) Close() error {
	if err := i.cur.Close(); err != nil {
		return errors.Wrap(err, "closing iterator")
	}

	return nil
}
