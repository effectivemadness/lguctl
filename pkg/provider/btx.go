package provider

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	lguctlaws "github.com/u-cto-devops/lguctl/pkg/aws"
	"github.com/u-cto-devops/lguctl/pkg/constants"
	"github.com/u-cto-devops/lguctl/pkg/schema"
	"github.com/u-cto-devops/lguctl/pkg/tools"
)

type lguctlProvider struct {
	StsClient       *sts.STS
	RoleARN         string
	RoleSessionName string
	ExternalID      string
	Duration        time.Duration
	ExpiryWindow    time.Duration
	credentials.Expiry
}

func NewAssumeProvider(account string, config *schema.Config, creds *credentials.Credentials) *lguctlProvider {
	region := viper.GetString("region")
	sess, err := lguctlaws.GetAwsSessionWithConfig(region)
	if err != nil {
		return nil
	}
	sess.Config.Credentials = creds

	var arn string
	if len(config.Alias) > 0 {
		arn = config.AssumeRoles[config.Alias[account]]
	}

	if len(arn) == 0 {
		arn = config.AssumeRoles[account]
	}

	return &lguctlProvider{
		StsClient:       sts.New(sess),
		RoleARN:         arn,
		RoleSessionName: config.Name,
		Duration:        time.Duration(int64(config.Duration)) * time.Second,
		ExpiryWindow:    constants.DefaultExpirationWindow,
		Expiry:          credentials.Expiry{},
	}
}

// Retrieve can retrieve sessionName
func (b *lguctlProvider) Retrieve() (credentials.Value, error) {
	role, err := b.assumeRole()
	if err != nil {
		return credentials.Value{}, err
	}

	b.SetExpiration(*role.Expiration, b.ExpiryWindow)
	return credentials.Value{
		AccessKeyID:     *role.AccessKeyId,
		SecretAccessKey: *role.SecretAccessKey,
		SessionToken:    *role.SessionToken,
	}, nil
}

// assumeRole assumes role
func (b *lguctlProvider) assumeRole() (*sts.Credentials, error) {
	var err error

	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(b.RoleARN),
		RoleSessionName: aws.String(b.roleSessionName()),
		DurationSeconds: aws.Int64(int64(b.Duration.Seconds())),
	}

	if b.ExternalID != "" {
		input.ExternalId = aws.String(b.ExternalID)
	}

	logrus.Debugf("Using STS endpoint %s", b.StsClient.Endpoint)

	resp, err := b.StsClient.AssumeRole(input)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("Generated credentials %s using AssumeRole, expires in %s", tools.FormatKeyForDisplay(*resp.Credentials.AccessKeyId), time.Until(*resp.Credentials.Expiration).String())

	return resp.Credentials, nil
}

// roleSessionName mean session name of assume role
func (b *lguctlProvider) roleSessionName() string {
	if len(b.RoleSessionName) == 0 {
		// Try to work out a role name that will hopefully end up unique.
		return fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	}

	return b.RoleSessionName
}
