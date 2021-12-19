package main // import "argc.in/kay"

import (
	"os"

	"argc.in/kay/cmd"
	_ "argc.in/kay/driver/bbolt"
	_ "argc.in/kay/driver/boltdb"
	_ "argc.in/kay/driver/diskv"
	_ "argc.in/kay/driver/etcd"
	_ "argc.in/kay/driver/leveldb"
	_ "argc.in/kay/driver/pebble"
	_ "argc.in/kay/driver/redis"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
