package main

import (
	"log"

	"argc.in/kay/cmd"
	_ "argc.in/kay/driver/leveldb"
	_ "argc.in/kay/driver/pebble"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
