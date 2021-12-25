package bbolt

import (
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
)

type iter struct {
	tx     *bbolt.Tx
	cur    *bbolt.Cursor
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
