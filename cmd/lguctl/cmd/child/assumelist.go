package child

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/constants"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

//Create RDS Token for IAM Authentication
func NewCmdAssumeList() *cobra.Command {
	return builder.NewCmd("list").
		WithDescription("List all accounts for assume role").
		SetFlags().
		RunWithNoArgs(funcAssumeList)
}

// Function for list command
func funcAssumeList(ctx context.Context, out io.Writer) error {
	return executor.RunExecutor(ctx, constants.NeedExpiredCheck, func(executor executor.Executor) error {
		if err := executor.Runner.PrintAssumeList(out); err != nil {
			logrus.Errorf(err.Error())
		}
		return nil
	})
}
