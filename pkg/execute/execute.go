package execute

import (
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/sirupsen/logrus"

	"github.com/u-cto-devops/lguctl/pkg/args"
	"github.com/u-cto-devops/lguctl/pkg/tools"
)

// Execute runs based on credentials
func Execute(creds *credentials.Credentials, out io.Writer, a *args.Argument) error {
	runnableCmd, err := exec.LookPath(a.Command)
	if err != nil {
		return err
	}
	values, err := creds.Get()
	if err != nil {
		return err
	}

	env, err := DefaultEnv()
	if err != nil {
		return err
	}

	logrus.Debugln("Setting subprocess env: AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY")
	env.Set("AWS_ACCESS_KEY_ID", values.AccessKeyID)
	env.Set("AWS_SECRET_ACCESS_KEY", values.SecretAccessKey)

	if values.SessionToken != "" {
		logrus.Debugln("Setting subprocess env: AWS_SESSION_TOKEN, AWS_SECURITY_TOKEN")
		env.Set("AWS_SESSION_TOKEN", values.SessionToken)
		env.Set("AWS_SECURITY_TOKEN", values.SessionToken)
	}

	if expiration, err := creds.ExpiresAt(); err == nil {
		logrus.Debugln("Setting subprocess env: AWS_SESSION_EXPIRATION")
		env.Set("AWS_SESSION_EXPIRATION", expiration.UTC().Format(time.RFC3339))
	}

	logrus.Debugf("Found executable %s", runnableCmd)

	argv := make([]string, 0, 1+len(a.Args))
	argv = append(argv, a.Command)
	argv = append(argv, a.Args...)

	return syscall.Exec(runnableCmd, argv, env)
}

// DefaultEnv returns default environment
func DefaultEnv() (Environ, error) {
	if err := tools.ClearOsEnv(); err != nil {
		return nil, err
	}

	env := Environ(os.Environ())

	return env, nil
}
