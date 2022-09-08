package ssm

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/service/ssm"

	"github.com/u-cto-devops/lguctl/pkg/aws"
	"github.com/u-cto-devops/lguctl/pkg/color"
	"github.com/u-cto-devops/lguctl/pkg/constants"
)

type Agent struct {
	Region string
	AWS    *aws.Client
}

func New(awsClient *aws.Client, region string) *Agent {
	return &Agent{
		Region: region,
		AWS:    awsClient,
	}
}

var (
	command      string = "/opt/fast_deployment.sh"
	documentName string = "AWS-RunShellScript"
)

// SelectNode chooses one instance to connect
func (a *Agent) SelectNode() (string, error) {
	nodeList, err := a.GetNodeList()
	if err != nil {
		return constants.EmptyString, err
	}

	var selected string
	prompt := &survey.Select{
		Message: "Choose instance: ",
		Options: nodeList,
	}
	survey.AskOne(prompt, &selected)

	return selected, nil
}

// GetNodeList returns instance list in the aws account
func (a *Agent) GetNodeList() ([]string, error) {
	var instanceList []string
	var err error
	if instanceList, err = a.AWS.GetInstanceListOnlyIds(instanceList, nil); err != nil {
		return nil, err
	}

	return instanceList, nil
}

// StartSession starts session
func (a *Agent) StartSession(out io.Writer, node string) error {
	color.Green.Fprintf(out, "Selected Instance: %s", node)
	params := &ssm.StartSessionInput{
		Target: &node,
	}

	session, endpoint, err := a.AWS.CreateSession(params)
	if err != nil {
		return err
	}

	jsonSession, err := json.Marshal(session)
	if err != nil {
		return err
	}

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return err
	}

	// call session-manager-plugin
	if err := callSubprocess("session-manager-plugin", string(jsonSession), a.Region, "StartSession", "default", string(jsonParams), endpoint); err != nil {
		return fmt.Errorf("error occurred in session-manager")
	}

	color.Yellow.Fprintf(out, "Delete session: %s", *session.SessionId)
	if err := a.AWS.DeleteSession(*session.SessionId); err != nil {
		return err
	}

	return nil
}

// callSubprocess runs subprocess with command
func callSubprocess(process string, args ...string) error {
	call := exec.Command(process, args...)
	call.Stderr = os.Stderr
	call.Stdout = os.Stdout
	call.Stdin = os.Stdin

	// ignore signal(sigint)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-sigs:
			case <-done:
				break
			}
		}
	}()
	defer close(done)

	// run subprocess
	if err := call.Run(); err != nil {
		return err
	}
	return nil
}

// StartSession starts session
func (a *Agent) DeployArtifact(out io.Writer, node string, branch string) error {
	color.Green.Fprintf(out, "Selected Instance: %s\n\n", node)
	params := &ssm.SendCommandInput{
		InstanceIds:  []*string{&node},
		DocumentName: &documentName,
		Parameters: map[string][]*string{
			"commands": {&command, &branch},
		},
	}

	result, err := a.AWS.SendCommand(params)
	if err != nil {
		return err
	}

	color.Green.Fprintf(out, result.Command.String())
	return nil
}
