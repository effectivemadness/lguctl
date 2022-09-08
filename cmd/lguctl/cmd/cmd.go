package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/u-cto-devops/lguctl/pkg/constants"
	"github.com/u-cto-devops/lguctl/pkg/tools"
	"github.com/u-cto-devops/lguctl/pkg/version"
)

var (
	cfgFile string
	v       string
)

// Get root command
func NewRootCommand(out, stderr io.Writer) *cobra.Command {
	cobra.OnInitialize(initConfig)
	rootCmd := &cobra.Command{
		Use:           "lguctl",
		Short:         "Private command line tool for u-cto-devops",
		Long:          `Private command line tool for u-cto-devops`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			cmd.Root().SetOutput(out)

			// Setup logs
			if err := tools.SetUpLogs(stderr, v); err != nil {
				return err
			}

			version := version.Get()

			logrus.Infof("lguctl %+v", version)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Group by commands
	groups := templates.CommandGroups{
		{
			Message: "managing configuration of lguctl",
			Commands: []*cobra.Command{
				NewInitCommand(),
			},
		},
		{
			Message: "commands related to aws IAM credentials",
			Commands: []*cobra.Command{
				NewRenewCredentialsCommand(),
			},
		},
		{
			Message: "commands for controlling assume role",
			Commands: []*cobra.Command{
				NewSetupCommand(),
				NewWhoCommand(),
				NewExecCommand(),
			},
		},
		{
			Message: "commands for retrieving information of AWS resources.",
			Commands: []*cobra.Command{
				NewCmdDescribeWebACL(),
				NewCmdHasIP(),
			},
		},
		{
			Message: "commands for setting maintenance mode to services.",
			Commands: []*cobra.Command{
				NewMaintenanceCommand(),
			},
		},
		{
			Message: "commands for start or stop loadtest environment.",
			Commands: []*cobra.Command{
				NewLoadtestCommand(),
			},
		},
		{
			Message: "commands for teleport",
			Commands: []*cobra.Command{
				NewListCommand(),
				NewSSHCommand(),
				NewLoginCommand(),
				NewStatusCommand(),
			},
		},
	}

	groups.Add(rootCmd)

	rootCmd.AddCommand(NewGetCommand())
	rootCmd.AddCommand(NewAssumeCommand())
	rootCmd.AddCommand(NewCmdCompletion())
	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.AddCommand(NewEcrLoginCommand())
	rootCmd.AddCommand(NewSSMCommand())
	rootCmd.AddCommand(NewFastDeploymentCommand())

	if checkMacOS() == nil {
		rootCmd.AddCommand(NewKeyChainCommand())
	}

	rootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", constants.DefaultLogLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")

	templates.ActsAsRootCommand(rootCmd, nil, groups...)

	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match
}
