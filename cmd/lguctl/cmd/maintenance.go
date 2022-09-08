package cmd

import (
	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/child"
)

// Assume role with setup
func NewMaintenanceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "maintenance",
		Short: "Setup maintenance mode to services",
	}

	cmd.AddCommand(child.NewCmdMaintenanceStart())
	cmd.AddCommand(child.NewCmdMaintenanceStop())
	return cmd
}
