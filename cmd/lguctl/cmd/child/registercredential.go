package child

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

// NewCmdRegisterCredential is a command for registering aws credentials to mac OS keychain
func NewCmdRegisterCredential() *cobra.Command {
	return builder.NewCmd("register-credential").
		WithDescription("Register AWS credentials to macOS keychain").
		SetFlags().
		RunWithNoArgs(funcRegisterCredential)
}

// funcRegisterCredential is main function for NewCmdRegisterCredential
func funcRegisterCredential(ctx context.Context, out io.Writer) error {
	return executor.RunExecutorWithoutCheckingConfig(ctx, func(executor executor.Executor) error {
		return executor.Runner.RegisterCredentialToKeyChain(out)
	})
}
