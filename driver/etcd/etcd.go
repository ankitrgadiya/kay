package etcd

import (
	"context"

	"github.com/pkg/errors"
	clientv3 "go.etcd.io/etcd/client/v3"

	"argc.in/kay/kv"
)

func init() {
	kv.Register("etcd", kv.DriverFunc(func() kv.KeyValue {
		return new(impl)
	}))
}

type impl struct {
	Endpoints []string `ini:"endpoints"`
	Username  string   `ini:"username,omitempty"`
	Password  string   `ini:"password,omitempty"`

	db *clientv3.Client
}

func (i *impl) Init() error {
	db, err := clientv3.New(i.clientConfig())
	if err != nil {
		return err
	}

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
	resp, err := i.db.Get(context.Background(), key)
	if err != nil {
		return nil, errors.Wrapf(err, "getting key: %s", key)
	}

	if len(resp.Kvs) == 0 {
		return nil, errors.Errorf("key not found: %s", key)
	}

	kv := resp.Kvs[0]
	if kv == nil {
		return nil, errors.Errorf("invalid response for key: %s", key)
	}

	return kv.Value, nil
}

func (i *impl) Set(key string, value []byte) error {
	_, err := i.db.Put(context.Background(), key, string(value))
	if err != nil {
		return errors.Wrapf(err, "setting key: %s", key)
	}

	return nil
}

func (i *impl) Delete(key string) error {
	_, err := i.db.Delete(context.Background(), key)
	if err != nil {
		return errors.Wrapf(err, "deleting key: %s", key)
	}

	return nil
}

func (i *impl) Watch(ctx context.Context, key string) <-chan kv.Event {
	eventChan := make(chan kv.Event)

	go func() {
		watchChan := i.db.Watch(context.Background(), key)
		for {
			select {
			case resp, ok := <-watchChan:
				if !ok {
					return
				}

				processEvents(&resp, eventChan)
			case <-ctx.Done():
				return
			}
		}
	}()

	return eventChan
}

func (i *impl) clientConfig() clientv3.Config {
	return clientv3.Config{
		Endpoints: i.Endpoints,
		Username:  i.Username,
		Password:  i.Password,
	}
}

func processEvents(resp *clientv3.WatchResponse, eventChan chan kv.Event) {
	for _, e := range resp.Events {
		eventChan <- kv.Event{
			Key:   string(e.Kv.Key),
			Value: e.Kv.Value,
		}
	}
}
