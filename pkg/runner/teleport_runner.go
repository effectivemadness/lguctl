package runner

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/u-cto-devops/lguctl/pkg/constants"
	"github.com/u-cto-devops/lguctl/pkg/teleport"
	"github.com/u-cto-devops/lguctl/pkg/tools"
)

// var clusterList = map[string]string{
// 	"preprod": "teleport.relengd.spddevops.com:443",
// 	"prod":    "teleport.releng.spddevops.com:443",
// }

var err error

// ListInstances lists all instances in registered teleport cluster
func (r *Runner) ListInstances(out io.Writer) error {
	t := teleport.NewTeleport()

	return t.ListNode()
}

// SSHToInstance connects to a instance with teleport
func (r *Runner) SSHToInstance(out io.Writer, args []string) error {
	t := teleport.NewTeleport()
	var userHost string

	if len(args) == 0 {
		hostname, err := t.SelectNode()
		if err != nil {
			return err
		}
		userHost = makeUserHost(hostname)
	} else {
		userHost = args[0]
		if ok, err := isValidHost(userHost); !ok {
			return err
		}
	}

	t.Config.UserHost = userHost

	return t.ConnectToNode()
}

func (r *Runner) GetClusterStatus(out io.Writer) error {
	t := teleport.NewTeleport()
	return t.DescribeStatus()
}

// LoginToCluster login to cluster
func (r *Runner) LoginToCluster(out io.Writer, args []string) error {
	var key string

	if len(args) == 0 {
		key, err = AskClusterTarget(r.Config.Teleport)
		if err != nil {
			return err
		}
	} else {
		key = args[0]
		if len(r.Config.Alias) > 0 {
			if _, ok := r.Config.Alias[key]; ok {
				key = r.Config.Alias[key]
			}
		}
	}

	t := teleport.NewTeleport()
	t.Config.Proxy = key
	t.Config.AuthConnector = constants.DefaultAuthConnector

	return t.Login()
}

// AskAssumeTarget asks assume target
func AskClusterTarget(clusterList map[string]string) (string, error) {
	var selected string

	keys := tools.GetKeys(clusterList)

	sort.Strings(keys)

	prompt := &survey.Select{
		Message: "Choose cluster: ",
		Options: keys,
	}
	survey.AskOne(prompt, &selected)

	if len(selected) == 0 {
		return selected, errors.New("choosing cluster has been canceled")
	}
	clusterName := clusterList[selected]

	return clusterName, nil
}

// isValidHost validates userHost value
func isValidHost(userHost string) (bool, error) {
	values := strings.Split(userHost, "@")
	if len(values) != 2 {
		return false, fmt.Errorf("argument should be <user>@<host> format")
	}

	return true, nil
}

// makeUserHost creates user@host string
func makeUserHost(hostname string) string {
	return fmt.Sprintf("%s@%s", constants.DefaultSSHUser, hostname)
}
