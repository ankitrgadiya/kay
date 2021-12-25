package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"argc.in/kay/kv"
)

func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:          "list DATABASE [PREFIX]",
		Short:        "List all the key-value pairs from the Database",
		RunE:         runList,
		SilenceUsage: true,
		Args:         cobra.RangeArgs(1, 2),
	}
}

func runList(cmd *cobra.Command, args []string) error {
	name, prefix := args[0], ""
	if len(args) == 2 {
		prefix = args[1]
	}

	keyvalue, closer, err := openDatabase(name)
	if err != nil {
		return err
	}
	defer closer.Close()

	lister, ok := keyvalue.(kv.Lister)
	if !ok {
		return errors.New("driver does not implement list operation")
	}

	i, err := lister.List(prefix)
	if err != nil {
		return err
	}
	defer i.Close()

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 3, 3, 3, ' ', 0)

	fmt.Fprintf(w, "KEY\tVALUE\n")
	for {
		key, value, ok := i.Next()
		if !ok {
			break
		}

		fmt.Fprintf(w, "%s\t%s\n", key, string(value))
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}
