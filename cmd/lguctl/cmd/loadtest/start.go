package loadtest

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

func NewCmdLoadtestStart() *cobra.Command {
	return builder.NewCmd("start").
		WithDescription("start loadtest rds and asg").
		SetFlags().
		RunWithNoArgs(funcLoadtestStart)
}

func funcLoadtestStart(ctx context.Context, out io.Writer) error {
	return executor.RunExecutorConfigReadOnly(ctx, func(executor executor.Executor) error {
		return executor.Runner.StartLoadtest(out)
	})
}
