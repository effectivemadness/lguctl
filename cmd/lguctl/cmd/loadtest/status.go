package loadtest

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

func NewCmdLoadtestStatus() *cobra.Command {
	return builder.NewCmd("status").
		WithDescription("get status loadtest rds and asg").
		SetFlags().
		RunWithNoArgs(funcLoadtestStatus)
}

func funcLoadtestStatus(ctx context.Context, out io.Writer) error {
	return executor.RunExecutorConfigReadOnly(ctx, func(executor executor.Executor) error {
		return executor.Runner.GetLoadtestStatus(out)
	})
}
