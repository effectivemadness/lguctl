package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/child"
)

// NewKeyChainCommand is a parent command line for all keychain-related CLI
func NewKeyChainCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keychain",
		Short: "Commands related to keychain function",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return checkMacOS()
		},
	}

	cmd.AddCommand(child.NewCmdRegisterCredential())
	cmd.AddCommand(child.NewCmdStatusCredential())
	return cmd
}

// checkMacOS checks if OS is darwin because keychain exists only in MacOS
func checkMacOS() error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("this command only works with MacOS")
	}

	return nil
}
