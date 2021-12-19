package cmd

import (
	"github.com/spf13/cobra"
)

func NewSetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "set DATABASE KEY VALUE",
		Short: "Sets the key from database",
		RunE:  runSet,
		Args:  cobra.ExactArgs(3),
	}
}

func runSet(cmd *cobra.Command, args []string) error {
	name, key, value := args[0], args[1], args[2]

	keyvalue, closer, err := openDatabase(name)
	if err != nil {
		return err
	}
	defer closer.Close()

	if err := keyvalue.Set(key, []byte(value)); err != nil {
		return err
	}

	return nil
}
