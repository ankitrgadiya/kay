package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
)

func NewFile(path string) (File, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	if err := ensureFile(path); err != nil {
		return nil, err
	}

	f, err := ini.Load(path)
	if err != nil {
		return nil, err
	}

	return &fileImpl{File: f}, nil
}

type fileImpl struct {
	*ini.File
}

func (f *fileImpl) AllSections() (map[string]Section, error) {
	list := f.File.Sections()
	sections := make(map[string]Section, len(list))

	for _, each := range list {
		if each.Name() == ini.DEFAULT_SECTION {
			continue
		}
		sections[each.Name()] = &sectionImpl{Section: each}
	}

	return sections, nil
}

func (f *fileImpl) Section(name string) (Section, error) {
	has := f.HasSection(name)
	if !has {
		return nil, errors.Errorf("section %s not found", name)
	}

	section, err := f.GetSection(name)
	if err != nil {
		return nil, err
	}

	return &sectionImpl{Section: section}, nil
}

func ensureFile(path string) error {
	_, err := os.Stat(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err == nil {
		return nil
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}
