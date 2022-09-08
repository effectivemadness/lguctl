package provider

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/sirupsen/logrus"
)

type DefaultProvider struct {
	Account         string
	AccessKeyID     string
	SecretAccessKey string
}

// NewDefaultProvider returns credentials
func NewDefaultProvider(account, accessKeyID, secretAccessKey string) *DefaultProvider {
	return &DefaultProvider{
		Account:         account,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
}

func (d *DefaultProvider) IsExpired() bool {
	return false
}

func (d *DefaultProvider) Retrieve() (val credentials.Value, err error) {
	logrus.Debugf("Looking up credentials of '%s'", d.Account)
	return credentials.Value{
		AccessKeyID:     d.AccessKeyID,
		SecretAccessKey: d.SecretAccessKey,
	}, err
}
