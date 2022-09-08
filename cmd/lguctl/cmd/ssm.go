package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

// NewSSMCommand is for connecting to node through teleport
func NewSSMCommand() *cobra.Command {
	return builder.NewCmd("ssm").
		WithDescription("connects to the node with amazon ssm").
		SetFlags().
		RunWithNoArgs(funcSSMConnect)
}

// funcSSMConnect
func funcSSMConnect(ctx context.Context, out io.Writer) error {
	return executor.RunExecutorWithoutCheckingConfig(ctx, func(executor executor.Executor) error {
		return executor.Runner.SSMToInstance(out)
	})
}
