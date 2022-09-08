package loadtest

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

func NewCmdLoadtestStop() *cobra.Command {
	return builder.NewCmd("stop").
		WithDescription("stop loadtest rds and asg").
		SetFlags().
		RunWithNoArgs(funcLoadtestStop)
}

func funcLoadtestStop(ctx context.Context, out io.Writer) error {
	return executor.RunExecutorConfigReadOnly(ctx, func(executor executor.Executor) error {
		return executor.Runner.StopLoadtest(out)
	})
}
