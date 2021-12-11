package config

type Section interface {
	DriverName() string
	Unmarshal(v interface{}) error
}

type File interface {
	AllSections() (map[string]Section, error)
	Section(name string) (Section, error)
}
