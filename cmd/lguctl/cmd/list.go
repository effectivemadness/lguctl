package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

// NewListCommand is for listing the instances in teleport cluster
func NewListCommand() *cobra.Command {
	return builder.NewCmd("list").
		WithDescription("list of registered instances in current logged-in teleport cluster").
		SetAliases([]string{"ls"}).
		SetFlags().
		RunWithNoArgs(funcList)
}

// funcList
func funcList(ctx context.Context, out io.Writer) error {
	return executor.RunExecutorWithoutCheckingConfig(ctx, func(executor executor.Executor) error {
		return executor.Runner.ListInstances(out)
	})
}
