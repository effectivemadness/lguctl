package child

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

// NewCmdStatusCredential is a command for retrieving aws credentials from mac OS keychain
func NewCmdStatusCredential() *cobra.Command {
	return builder.NewCmd("status").
		WithDescription("Show status of AWS credentials from macOS keychain").
		SetFlags().
		RunWithNoArgs(funcStatusCredential)
}

// funcStatusCredential is main function for NewCmdStatusCredential
func funcStatusCredential(ctx context.Context, out io.Writer) error {
	return executor.RunExecutorWithoutCheckingConfig(ctx, func(executor executor.Executor) error {
		if err := executor.Runner.StatusCredentialFromKeyChain(out); err != nil {
			logrus.Errorf(err.Error())
		}
		return nil
	})
}
