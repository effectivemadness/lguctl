package runner

import (
	"fmt"
	"io"
	"strings"

	"github.com/u-cto-devops/lguctl/pkg/constants"
	"github.com/u-cto-devops/lguctl/pkg/ssm"
)

// SSMToInstance connects to a instance with amazon ssm
func (r *Runner) SSMToInstance(out io.Writer) error {
	agent := ssm.New(&r.AWSClient, r.Region)

	fmt.Println("Please choose the instance")
	userHost, err := agent.SelectNode()
	if err != nil {
		return err
	}

	return agent.StartSession(out, getInstanceID(userHost))
}

// getInstanceID returns instance ID only
func getInstanceID(userHost string) string {
	sp := strings.Split(userHost, constants.NameDelimiter)
	return sp[1]
}

func (r *Runner) DeployNewArtifact(out io.Writer, arg []string) error {
	agent := ssm.New(&r.AWSClient, r.Region)

	fmt.Println("Please choose the instance")
	userHost, err := agent.SelectNode()
	if err != nil {
		return err
	}

	return agent.DeployArtifact(out, getInstanceID(userHost), arg[0])
}
