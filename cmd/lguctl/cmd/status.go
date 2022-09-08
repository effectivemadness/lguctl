package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

// NewStatusCommand is for checking auth status of teleport cluster
func NewStatusCommand() *cobra.Command {
	return builder.NewCmd("status").
		WithDescription("status of teleport cluster").
		SetFlags().
		RunWithNoArgs(funcStatus)
}

// funcStatus
func funcStatus(ctx context.Context, out io.Writer) error {
	return executor.RunExecutorWithoutCheckingConfig(ctx, func(executor executor.Executor) error {
		return executor.Runner.GetClusterStatus(out)
	})
}
