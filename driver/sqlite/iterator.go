package sqlite

import (
	"database/sql"

	"github.com/pkg/errors"
)

type iter struct {
	rows *sql.Rows
}

func (i *iter) Next() (string, []byte, bool) {
	has := i.rows.Next()
	if !has {
		return "", nil, false
	}

	var key, value string

	if err := i.rows.Scan(&key, &value); err != nil {
		return "", nil, false
	}

	return key, []byte(value), true
}

func (i *iter) Close() error {
	if err := i.rows.Close(); err != nil {
		return errors.Wrap(err, "closing iterator")
	}

	return nil
}
