package main // import "argc.in/kay"

import (
	"log"

	"argc.in/kay/cmd"
	_ "argc.in/kay/driver/boltdb"
	_ "argc.in/kay/driver/leveldb"
	_ "argc.in/kay/driver/pebble"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
