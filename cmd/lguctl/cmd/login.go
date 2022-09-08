package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

// NewLoginCommand is for logging in to teleport clusters
func NewLoginCommand() *cobra.Command {
	return builder.NewCmd("login").
		WithDescription("login to teleport cluster").
		SetFlags().
		RunWithArgsAndCmd(funcLogin)
}

// funcLogin
func funcLogin(ctx context.Context, out io.Writer, cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return cmd.Help()
	}

	return executor.RunExecutorConfigReadOnly(ctx, func(executor executor.Executor) error {
		return executor.Runner.LoginToCluster(out, args)
	})
}
