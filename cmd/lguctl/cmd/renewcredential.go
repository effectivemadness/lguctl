package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/constants"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

// renew credentials
func NewRenewCredentialsCommand() *cobra.Command {
	return builder.NewCmd("renew-credential").
		WithDescription("recreates aws credential of profile").
		SetFlags().
		RunWithNoArgs(funcRenewCredentials)
}

// funcRenewCredentials
func funcRenewCredentials(ctx context.Context, out io.Writer) error {
	return executor.RunExecutor(ctx, constants.SkipExpiredCheck, func(executor executor.Executor) error {
		return executor.Runner.RenewCredentials(out)
	})
}
