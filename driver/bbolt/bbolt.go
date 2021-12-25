// Package bbolt is the Driver implementation for Ben Johnson's BoltDB.
// It uses the Etcd's fork of the original project.
package bbolt

import (
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"

	"argc.in/kay/kv"
)

func init() {
	kv.Register("bbolt", kv.DriverFunc(func() kv.KeyValue {
		return new(impl)
	}))
}

type impl struct {
	Path   string `ini:"path"`
	Bucket string `ini:"bucket"`

	db *bbolt.DB
}

func (i *impl) Init() error {
	db, err := bbolt.Open(i.Path, 0600, nil)
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

	err := i.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(i.Bucket))
		if err != nil {
			return err
		}

		value = b.Get([]byte(key))

		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "getting key: %s", key)
	}

	if value == nil {
		return nil, errors.Errorf("key not found: %s", key)
	}

	return value, nil
}

func (i *impl) Set(key string, value []byte) error {
	if err := i.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(i.Bucket))
		if err != nil {
			return err
		}

		if err := b.Put([]byte(key), value); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.Wrapf(err, "setting key: %s", key)
	}

	return nil
}

func (i *impl) Delete(key string) error {
	if err := i.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(i.Bucket))
		if err != nil {
			return err
		}

		if err := b.Delete([]byte(key)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.Wrapf(err, "deleting key: %s", key)
	}

	return nil
}

func (i *impl) List(prefix string) (kv.Iterator, error) {
	tx, err := i.db.Begin(false)
	if err != nil {
		return nil, errors.Wrap(err, "opening read-only transaction")
	}

	b := tx.Bucket([]byte(i.Bucket))
	if b == nil {
		return nil, errors.Errorf("getting bucket %s: not found", i.Bucket)
	}

	return &iter{
		tx:     tx,
		cur:    b.Cursor(),
		prefix: prefix,
	}, nil
}
