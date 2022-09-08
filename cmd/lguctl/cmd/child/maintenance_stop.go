package child

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/executor"
)

func NewCmdMaintenanceStop() *cobra.Command {
	return builder.NewCmd("stop").
		WithDescription("stop to setup maintenance").
		SetFlags().
		RunWithArgsAndCmd(funcStopMaintenance)
}

// funcStartMaintenance
func funcStopMaintenance(ctx context.Context, out io.Writer, cmd *cobra.Command, args []string) error {
	return executor.RunExecutorConfigReadOnly(ctx, func(executor executor.Executor) error {
		return executor.Runner.StopMaintenance(out, args)
	})
}
