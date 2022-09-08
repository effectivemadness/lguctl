package cmd

import (
	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/child"
)

// Get information with lguctl
func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get token or information with lguctl",
	}

	cmd.AddCommand(child.NewCmdRDSToken())
	cmd.AddCommand(child.NewCmdIAMRoleArn())
	cmd.AddCommand(child.NewCmdIAMGroup())
	cmd.AddCommand(child.NewCmdIAMPolicy())
	return cmd
}
