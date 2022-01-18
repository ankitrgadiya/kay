package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"argc.in/kay/kv"
)

func init() {
	kv.Register("sqlite3", kv.DriverFunc(func() kv.KeyValue {
		return new(impl)
	}))
}

type impl struct {
	DSN         string `ini:"dsn"`
	Table       string `ini:"table"`
	KeyColumn   string `ini:"keyColumn"`
	ValueColumn string `ini:"valueColumn"`

	conn *sql.DB
}

func (i *impl) Init() error {
	conn, err := sql.Open("sqlite3", i.DSN)
	if err != nil {
		return errors.Wrapf(err, "opening %s", i.DSN)
	}

	i.conn = conn
	return nil
}

func (i *impl) Close() error {
	if err := i.conn.Close(); err != nil {
		return errors.Wrapf(err, "closing %s", i.DSN)
	}

	return nil
}

func (i *impl) Get(key string) ([]byte, error) {
	query := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ? LIMIT 1`, i.ValueColumn, i.Table, i.KeyColumn)

	var value string
	if err := i.conn.QueryRowContext(context.Background(), query, key).Scan(&value); err != nil {
		return nil, errors.Wrapf(err, "getting key: %s", key)
	}

	return []byte(value), nil
}

func (i *impl) Set(key string, value []byte) error {
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s)
                          VALUES(:1, :2)
                          ON CONFLICT(%s)
                          DO UPDATE SET %s=excluded.%s`, i.Table, i.KeyColumn, i.ValueColumn, i.KeyColumn, i.KeyColumn, i.KeyColumn)

	if _, err := i.conn.ExecContext(context.Background(), query, key, string(value)); err != nil {
		return errors.Wrapf(err, "setting key: %s", key)
	}

	return nil
}

func (i *impl) Delete(key string) error {
	query := fmt.Sprintf(`DELETE FROM %s
                          WHERE %s = ?`, i.Table, i.KeyColumn)

	if _, err := i.conn.ExecContext(context.Background(), query, key); err != nil {
		return errors.Wrapf(err, "deleting key: %s", key)
	}

	return nil
}

func (i *impl) List(prefix string) (kv.Iterator, error) {
	query := fmt.Sprintf(`SELECT %s, %s FROM %s
                          WHERE %s LIKE ?`, i.KeyColumn, i.ValueColumn, i.Table, i.KeyColumn)

	rows, err := i.conn.QueryContext(context.Background(), query, prefix+"%")
	if err != nil {
		return nil, errors.Wrapf(err, "listing keys with prefix: %s", prefix)
	}

	return &iter{rows: rows}, nil
}

func (i *impl) Search(term string) (kv.Iterator, error) {
	query := fmt.Sprintf(`SELECT %s, %s FROM %s
                          WHERE %s LIKE ?`, i.KeyColumn, i.ValueColumn, i.Table, i.KeyColumn)

	rows, err := i.conn.QueryContext(context.Background(), query, fmt.Sprintf("%%%s%%", term))
	if err != nil {
		return nil, errors.Wrapf(err, "listing keys with prefix: %s", term)
	}

	return &iter{rows: rows}, nil
}
