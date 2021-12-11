package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func NewDatabasesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "databases",
		Short: "Lists all the configured databases",
		RunE:  runDatabases,
		Args:  cobra.NoArgs,
	}
}

func runDatabases(cmd *cobra.Command, args []string) error {
	sections, err := conf.AllSections()
	if err != nil {
		return err
	}

	if len(sections) == 0 {
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 3, 3, 3, ' ', 0)

	fmt.Fprintf(w, "NAME\tDRIVER\n")
	for name, section := range sections {
		fmt.Fprintf(w, "%s\t%s\n", name, section.DriverName())
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}
