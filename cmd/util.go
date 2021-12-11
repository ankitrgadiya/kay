package cmd

import "argc.in/kay/kv"

func openDatabase(name string) (kv.KeyValue, error) {
	s, err := conf.Section(name)
	if err != nil {
		return nil, err
	}

	keyvalue, err := kv.Open(name, s)
	if err != nil {
		return nil, err
	}

	return keyvalue, nil
}
