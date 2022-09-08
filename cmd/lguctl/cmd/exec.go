package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/args"
	"github.com/u-cto-devops/lguctl/pkg/constants"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

// NewExecCommand is a command for executing aws command with private configuration
func NewExecCommand() *cobra.Command {
	return builder.NewCmd("exec").
		WithDescription("Execute aws command with private configuration").
		SetFlags().
		RunWithArgs(funcExec)
}

// funcExec
func funcExec(ctx context.Context, out io.Writer, arg []string) error {
	a, err := args.Parse(arg)
	if err != nil {
		return err
	}

	return executor.RunExecutor(ctx, constants.NeedExpiredCheck, func(executor executor.Executor) error {
		return executor.Runner.Exec(out, a)
	})
}
