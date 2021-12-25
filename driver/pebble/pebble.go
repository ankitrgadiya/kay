// Package pebble is the Driver implementation for CockroachDB's Pure Go RocksDB
// implementation called Pebble.
package pebble

import (
	"github.com/cockroachdb/pebble"
	"github.com/pkg/errors"

	"argc.in/kay/kv"
)

func init() {
	kv.Register("pebble", kv.DriverFunc(func() kv.KeyValue {
		return new(impl)
	}))
}

type impl struct {
	Path string `ini:"path"`

	db *pebble.DB
}

func (i *impl) Init() error {
	db, err := pebble.Open(i.Path, nil)
	if err != nil {
		return errors.Wrapf(err, "opening %s", i.Path)
	}

	i.db = db
	return nil
}

func (i *impl) Close() error {
	if err := i.db.Close(); err != nil {
		return errors.Wrapf(err, "closing %s", i.Path)
	}

	return nil
}

func (i *impl) Get(key string) ([]byte, error) {
	val, closer, err := i.db.Get([]byte(key))
	if err != nil {
		return nil, errors.Wrapf(err, "getting key: %s", key)
	}
	defer closer.Close()

	buf := make([]byte, len(val))
	copy(buf, val)

	return buf, nil
}

func (i *impl) Set(key string, value []byte) error {
	if err := i.db.Set([]byte(key), value, pebble.Sync); err != nil {
		return errors.Wrapf(err, "setting key: %s", key)
	}

	return nil
}

func (i *impl) Delete(key string) error {
	if err := i.db.Delete([]byte(key), pebble.Sync); err != nil {
		return errors.Wrapf(err, "deleting key: %s", key)
	}

	return nil
}

func (i *impl) List(prefix string) (kv.Iterator, error) {
	var opt *pebble.IterOptions
	if prefix != "" {
		opt = &pebble.IterOptions{
			LowerBound: []byte(prefix),
		}
	}

	cur := i.db.NewIter(opt)
	if ok := cur.First(); !ok {
		return nil, errors.Errorf("getting iterator: invalid")
	}

	return &iter{cur: cur}, nil
}
