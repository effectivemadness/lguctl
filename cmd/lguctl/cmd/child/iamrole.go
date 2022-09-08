package child

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

func NewCmdIAMRoleArn() *cobra.Command {
	return builder.NewCmd("role-arn").
		WithDescription("Get IAM role arn assigned to User").
		SetFlags().
		RunWithArgsAndCmd(funcGetIAMRoleArn)
}

// funcGetIAMGroup
func funcGetIAMRoleArn(ctx context.Context, out io.Writer, cmd *cobra.Command, args []string) error {
	return executor.RunExecutorWithoutCheckingConfig(ctx, func(executor executor.Executor) error {
		return executor.Runner.GetIAMRoleArnForUser(out)
	})
}
