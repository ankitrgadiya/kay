package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"argc.in/kay/config"
)

const version = "v0.1.1"

var (
	confPath string
	conf     config.File
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if err := initConfigE(); err != nil {
		log.Fatalf("failed to load config %s: %+v", confPath, err)
	}

}

func initConfigE() error {
	var err error

	confDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	defaultPath := filepath.Join(confDir, "kay", "kay.conf")

	if confPath == "" {
		confPath = defaultPath
	}

	conf, err = config.NewFile(confPath)
	if err != nil {
		return err
	}

	return nil
}
