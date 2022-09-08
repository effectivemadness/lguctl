package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

// Initialize lguctl configuration
func NewInitCommand() *cobra.Command {
	return builder.NewCmd("init").
		WithDescription("initialize lguctl command line tool").
		RunWithNoArgs(funcInit)
}

// funcInit
func funcInit(ctx context.Context, _ io.Writer) error {
	return executor.RunExecutorWithoutCheckingConfig(ctx, func(executor executor.Executor) error {
		if err := executor.Runner.InitConfiguration(); err != nil {
			return err
		}

		return nil
	})
}
