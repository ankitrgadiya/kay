package redis

import (
	"context"

	v8 "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"argc.in/kay/kv"
)

func init() {
	kv.Register("redis", kv.DriverFunc(func() kv.KeyValue {
		return new(impl)
	}))
}

type impl struct {
	Network  string `ini:"network,omitempty"`
	Address  string `ini:"address"`
	Username string `ini:"username,omitempty"`
	Password string `ini:"password,omitempty"`

	db *v8.Client
}

func (i *impl) Init() error {
	db := v8.NewClient(&v8.Options{
		Network:  i.Network,
		Addr:     i.Address,
		Username: i.Username,
		Password: i.Password,
	})

	i.db = db
	return nil
}

func (i *impl) Close() error {
	if err := i.db.Close(); err != nil {
		return errors.Wrap(err, "closing")
	}

	return nil
}

func (i *impl) Get(key string) ([]byte, error) {
	value, err := i.db.Get(context.Background(), key).Bytes()
	if err != nil {
		return nil, errors.Wrapf(err, "getting key: %s", key)
	}

	return value, nil
}

func (i *impl) Set(key string, value []byte) error {
	if err := i.db.Set(context.Background(), key, value, 0).Err(); err != nil {
		return errors.Wrapf(err, "setting key: %s", key)
	}

	return nil
}

func (i *impl) Delete(key string) error {
	if err := i.db.Del(context.Background(), key).Err(); err != nil {
		return errors.Wrapf(err, "deleting key: %s", key)
	}

	return nil
}
