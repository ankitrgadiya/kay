package cmd

import (
	"io"
	"os"

	"argc.in/kay/kv"
)

func openDatabase(name string) (kv.KeyValue, io.Closer, error) {
	s, err := conf.Section(name)
	if err != nil {
		return nil, nil, err
	}

	keyvalue, closer, err := kv.Open(s)
	if err != nil {
		return nil, nil, err
	}

	return keyvalue, closer, nil
}

func isInteractive(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}

	stat, err := f.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) == os.ModeCharDevice
}
