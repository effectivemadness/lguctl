package keymanager

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/u-cto-devops/lguctl/pkg/constants"
	"github.com/u-cto-devops/lguctl/pkg/tools"
)

type Credential struct {
	Credentials      credentials.Value
	AccessKeyID      string `json:"AccessKeyID"`
	SecretAccessKey  string `json:"SecretAccessKey"`
	SessionToken     string `json:"SessionToken"`
	Label            string `json:"label"`
	CreatedTime      time.Time
	ModificationTime time.Time
}

// GetCredentialsFromPrompt gets credentials from command line
func (k *KeyChain) GetCredentialsFromPrompt() (*credentials.Value, error) {
	var accessKeyID, secretAccessKeyID string
	var err error

	if accessKeyID, err = tools.Ask("> Access Key ID: ", !constants.Secret); err != nil {
		return nil, err
	}

	if secretAccessKeyID, err = tools.Ask("> Secret Access Key ID: ", constants.Secret); err != nil {
		return nil, err
	}

	creds := credentials.Value{AccessKeyID: accessKeyID, SecretAccessKey: secretAccessKeyID}

	return &creds, nil
}
