package cmd

import (
	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/child"
)

// Command related to assume role
func NewAssumeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assume",
		Short: "do work about assume role",
	}

	cmd.AddCommand(child.NewCmdAssumeList())
	return cmd
}
