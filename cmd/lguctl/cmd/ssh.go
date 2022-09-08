package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

// NewSSHCommand is for connecting to node through teleport
func NewSSHCommand() *cobra.Command {
	return builder.NewCmd("ssh").
		WithDescription("connecting to node through teleport").
		SetFlags().
		RunWithArgsAndCmd(funcConnect)
}

// funcConnect
func funcConnect(ctx context.Context, out io.Writer, cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("lguctl ssh can get at most one argument")
	}

	return executor.RunExecutorWithoutCheckingConfig(ctx, func(executor executor.Executor) error {
		return executor.Runner.SSHToInstance(out, args)
	})
}
