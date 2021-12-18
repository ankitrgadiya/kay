// Package leveldb is the Driver implementation for LevelDB embedded database.
// It uses the syndtr/goleveldb implementation.
package leveldb

import (
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"argc.in/kay/kv"
)

func init() {
	kv.Register("leveldb", kv.DriverFunc(func() kv.KeyValue {
		return new(impl)
	}))
}

type impl struct {
	Path string `ini:"path"`

	db *leveldb.DB
}

func (i *impl) Init() error {
	db, err := leveldb.OpenFile(i.Path, nil)
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
	val, err := i.db.Get([]byte(key), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "getting key: %s", key)
	}

	return val, nil
}

func (i *impl) Set(key string, value []byte) error {
	if err := i.db.Put([]byte(key), value, &opt.WriteOptions{Sync: true}); err != nil {
		return errors.Wrapf(err, "setting key: %s", key)
	}

	return nil
}

func (i *impl) Delete(key string) error {
	if err := i.db.Delete([]byte(key), &opt.WriteOptions{Sync: true}); err != nil {
		return errors.Wrapf(err, "deleting key: %s", key)
	}

	return nil
}
