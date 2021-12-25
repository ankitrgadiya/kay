package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type iter struct {
	tx     *bolt.Tx
	cur    *bolt.Cursor
	prefix string
	seeked bool
}

func (i *iter) seek() ([]byte, []byte) {
	i.seeked = true
	if i.prefix != "" {
		return i.cur.Seek([]byte(i.prefix))
	}

	return i.cur.First()
}

func (i *iter) Next() (string, []byte, bool) {
	var key, value []byte
	if !i.seeked {
		key, value = i.seek()
	} else {
		key, value = i.cur.Next()
	}

	if key == nil && value == nil {
		return "", nil, false
	}

	return string(key), value, true
}

func (i *iter) Close() error {
	if err := i.tx.Rollback(); err != nil {
		return errors.Wrap(err, "closing iterator (transaction)")
	}

	return nil
}
