package cmd

import (
	"context"
	"fmt"
	"io"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"argc.in/kay/kv"
)

func NewWatchCommand() *cobra.Command {
	return &cobra.Command{
		Use:          "watch DATABASE KEY",
		Short:        "Watch the changes on key",
		RunE:         runWatch,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(2),
	}
}

func runWatch(cmd *cobra.Command, args []string) error {
	name, key := args[0], args[1]

	keyvalue, closer, err := openDatabase(name)
	if err != nil {
		return err
	}
	defer closer.Close()

	watcher, ok := keyvalue.(kv.Watcher)
	if !ok {
		return errors.New("driver does not implement watch operation")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	displayEvents(ctx, cmd.OutOrStdout(), watcher.Watch(ctx, key))
	return nil
}

func displayEvents(ctx context.Context, w io.Writer, watchChan <-chan kv.Event) {
	for {
		select {
		case e, ok := <-watchChan:
			if !ok {
				return
			}

			fmt.Fprintf(w, "%s\t%s\n", e.Key, e.Value)
		case <-ctx.Done():
			return
		}
	}
}
