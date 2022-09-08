package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

// Assume role with setup
func NewWhoCommand() *cobra.Command {
	return builder.NewCmd("who").
		WithDescription("check the account information of current shell").
		SetFlags().
		RunWithNoArgs(funcWho)
}

// funcWho
func funcWho(ctx context.Context, out io.Writer) error {
	return executor.RunExecutorWithoutCheckingConfig(ctx, func(executor executor.Executor) error {
		return executor.Runner.Who(out)
	})
}
