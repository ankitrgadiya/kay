package badger

import (
	v3 "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"

	"argc.in/kay/kv"
)

func init() {
	kv.Register("badger", kv.DriverFunc(func() kv.KeyValue {
		return new(impl)
	}))
}

type impl struct {
	Path string `ini:"path"`

	db *v3.DB
}

func (i *impl) Init() error {
	db, err := v3.Open(v3.DefaultOptions(i.Path).WithLoggingLevel(v3.ERROR))
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
	var value []byte
	if err := i.db.View(func(tx *v3.Txn) error {
		item, err := tx.Get([]byte(key))
		if err != nil {
			return err
		}

		value, err = item.ValueCopy(value)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, errors.Wrapf(err, "getting key: %s", key)
	}
	return value, nil
}

func (i *impl) Set(key string, value []byte) error {
	if err := i.db.Update(func(tx *v3.Txn) error {
		if err := tx.Set([]byte(key), value); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.Wrapf(err, "setting key: %s", key)
	}

	return nil
}

func (i *impl) Delete(key string) error {
	if err := i.db.Update(func(tx *v3.Txn) error {
		if err := tx.Delete([]byte(key)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.Wrapf(err, "deleting key: %s", key)
	}

	return nil
}
