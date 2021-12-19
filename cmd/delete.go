package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"argc.in/kay/kv"
)

func NewDelCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "del DATABASE KEY VALUE",
		Short: "Deletes the key from database",
		RunE:  runDel,
		Args:  cobra.ExactArgs(2),
	}
}

func runDel(cmd *cobra.Command, args []string) error {
	name, key := args[0], args[1]

	keyvalue, closer, err := openDatabase(name)
	if err != nil {
		return err
	}
	defer closer.Close()

	deleter, ok := keyvalue.(kv.Deleter)
	if !ok {
		return errors.New("driver does not implement delete operation")
	}

	if err := deleter.Delete(key); err != nil {
		return err
	}

	return nil
}
