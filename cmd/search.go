package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"argc.in/kay/kv"
)

func NewSearchCommand() *cobra.Command {
	return &cobra.Command{
		Use:          "search DATABASE [TERM]",
		Short:        "Search for the term in all keys",
		RunE:         runSearch,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(2),
	}
}

func runSearch(cmd *cobra.Command, args []string) error {
	name, term := args[0], args[1]

	keyvalue, closer, err := openDatabase(name)
	if err != nil {
		return err
	}
	defer closer.Close()

	searcher, ok := keyvalue.(kv.Searcher)
	if !ok {
		return errors.New("driver does not implement search operation")
	}

	i, err := searcher.Search(term)
	if err != nil {
		return err
	}
	defer i.Close()

	return displayList(cmd.OutOrStdout(), i)
}
