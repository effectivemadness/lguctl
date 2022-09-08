package cmd

import (
	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/loadtest"
)

// Assume role with setup
func NewLoadtestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "loadtest",
		Short: "Start or stop loadtest environment",
	}

	cmd.AddCommand(loadtest.NewCmdLoadtestStatus())
	cmd.AddCommand(loadtest.NewCmdLoadtestStart())
	cmd.AddCommand(loadtest.NewCmdLoadtestStop())
	return cmd
}
