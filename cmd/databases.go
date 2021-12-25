package cmd

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"argc.in/kay/config"
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

	return displayDatabases(cmd.OutOrStdout(), sections)
}

func displayDatabases(w io.Writer, sections map[string]config.Section) error {
	tw := tabwriter.NewWriter(w, 3, 3, 3, ' ', 0)

	if isInteractive(w) {
		fmt.Fprintf(w, "NAME\tDRIVER\n")
	}

	for name, section := range sections {
		fmt.Fprintf(w, "%s\t%s\n", name, section.DriverName())
	}

	return tw.Flush()
}
