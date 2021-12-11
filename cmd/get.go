package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "get DATABASE KEY",
		Short: "Gets the key from database",
		RunE:  runGet,
		Args:  cobra.ExactArgs(2),
	}
}

func runGet(cmd *cobra.Command, args []string) error {
	name, key := args[0], args[1]

	keyvalue, err := openDatabase(name)
	if err != nil {
		return err
	}

	value, err := keyvalue.Get(key)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprint(cmd.OutOrStdout(), string(value)); err != nil {
		return err
	}

	return nil
}
