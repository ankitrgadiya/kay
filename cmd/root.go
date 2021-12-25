package cmd

import "github.com/spf13/cobra"

func NewCommand() *cobra.Command {
	c := &cobra.Command{
		Use:     "kay",
		Example: "kay get test key",
		Version: version,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	c.AddCommand(
		NewGetCommand(),
		NewSetCommand(),
		NewDelCommand(),
		NewWatchCommand(),
		NewListCommand(),
		NewDatabasesCommand())

	c.PersistentFlags().StringVar(&confPath, "config", "", "Path for config file")

	return c
}

func Execute() error {
	return NewCommand().Execute()
}
