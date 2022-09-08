package child

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

func NewCmdIAMGroup() *cobra.Command {
	return builder.NewCmd("group").
		WithDescription("Get IAM group assigned to User").
		SetFlags().
		RunWithArgsAndCmd(funcGetIAMGroup)
}

// funcGetIAMGroup
func funcGetIAMGroup(ctx context.Context, out io.Writer, cmd *cobra.Command, args []string) error {
	return executor.RunExecutorWithoutCheckingConfig(ctx, func(executor executor.Executor) error {
		return executor.Runner.GetIAMGroupForUser(out)
	})
}
