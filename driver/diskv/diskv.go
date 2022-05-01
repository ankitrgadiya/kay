// Package diskv is the Driver implementation for Peter Bourgon's diskv
// database. It uses the peterbourgon/diskv package.
package diskv

import (
	"github.com/peterbourgon/diskv/v3"
	"github.com/pkg/errors"

	"argc.in/kay/kv"
)

func init() {
	kv.Register("diskv", kv.DriverFunc(func() kv.KeyValue {
		return new(impl)
	}))
}

type impl struct {
	Path string `ini:"path"`

	db *diskv.Diskv
}

func (i *impl) Init() error {
	db := diskv.New(diskv.Options{
		BasePath:  i.Path,
		Transform: func(s string) []string { return []string{} },
	})

	i.db = db
	return nil
}

func (i *impl) Get(key string) ([]byte, error) {
	value, err := i.db.Read(key)
	if err != nil {
		return nil, errors.Wrapf(err, "getting key: %s", key)
	}

	return value, nil
}

func (i *impl) Set(key string, value []byte) error {
	if err := i.db.Write(key, value); err != nil {
		return errors.Wrapf(err, "setting key: %s", key)
	}

	return nil
}

func (i *impl) Delete(key string) error {
	if err := i.db.Erase(key); err != nil {
		return errors.Wrapf(err, "deleting key: %s", key)
	}

	return nil
}
