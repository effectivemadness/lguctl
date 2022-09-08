package child

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

func NewCmdMaintenanceStart() *cobra.Command {
	return builder.NewCmd("start").
		WithDescription("start to setup maintenance").
		SetFlags().
		RunWithArgsAndCmd(funcStartMaintenance)
}

// funcStartMaintenance
func funcStartMaintenance(ctx context.Context, out io.Writer, cmd *cobra.Command, args []string) error {
	return executor.RunExecutorConfigReadOnly(ctx, func(executor executor.Executor) error {
		return executor.Runner.StartMaintenance(out, args)
	})
}
