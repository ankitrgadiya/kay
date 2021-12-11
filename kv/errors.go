package kv

import "github.com/pkg/errors"

var (
	ErrSectionNotFound = errors.New("database not found in the config")
	ErrDriverNotFound  = errors.New("database driver not found")
)
