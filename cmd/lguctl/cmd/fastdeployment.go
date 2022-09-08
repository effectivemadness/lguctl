package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/color"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

// NewFastDeploymentCommand is for fastly deploy new branch artifact to the node
func NewFastDeploymentCommand() *cobra.Command {
	return builder.NewCmd("fast-deployment").
		WithDescription("deploy new branch artifact to the node").
		SetFlags().
		RunWithArgs(funcFastDeployment)
}

// funcFastDeployment./
func funcFastDeployment(ctx context.Context, out io.Writer, args []string) error {
	if len(args) != 1 {
		color.Red.Fprintf(out, "Input branch is only one.\n\nusage: lguctl fast-deployment [branch]")
		return nil
	}
	return executor.RunExecutorWithoutCheckingConfig(ctx, func(executor executor.Executor) error {
		return executor.Runner.DeployNewArtifact(out, args)
	})
}
