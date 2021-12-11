package kv

import "sync"

var (
	mu      sync.Mutex
	drivers = make(map[string]Driver)
)

// Register stores the association of the name with the Driver implementation
// internally. This function must be called for each Driver.
func Register(name string, driver Driver) {
	mu.Lock()
	if _, has := drivers[name]; !has {
		drivers[name] = driver
	}
	mu.Unlock()
}

func getDriver(name string) Driver {
	mu.Lock()
	defer mu.Unlock()
	return drivers[name]
}
