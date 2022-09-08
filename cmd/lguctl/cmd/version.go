package cmd

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/u-cto-devops/lguctl/cmd/lguctl/cmd/builder"
	"github.com/u-cto-devops/lguctl/pkg/version"
)

// Get lguctl version
func NewVersionCommand() *cobra.Command {
	return builder.NewCmd("version").
		WithDescription("Print the version information").
		SetAliases([]string{"v"}).
		RunWithNoArgs(funcVersion)
}

// funcVersion
func funcVersion(_ context.Context, _ io.Writer) error {
	return version.Controller{}.Print(version.Get())
}
