package config

import "gopkg.in/ini.v1"

type sectionImpl struct {
	*ini.Section
}

func (s *sectionImpl) DriverName() string {
	has := s.HasKey("driver")
	if !has {
		return ""
	}

	key, err := s.GetKey("driver")
	if err != nil {
		return ""
	}

	return key.MustString("")
}

func (s *sectionImpl) Unmarshal(v interface{}) error {
	return s.MapTo(v)
}
