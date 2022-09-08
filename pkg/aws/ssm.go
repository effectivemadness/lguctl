package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func GetSSMClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *ssm.SSM {
	if creds == nil {
		return ssm.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return ssm.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}

// Create start session
func (c *Client) CreateSession(input *ssm.StartSessionInput) (*ssm.StartSessionOutput, string, error) {
	subctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	sess, err := c.SSMClient.StartSessionWithContext(subctx, input)
	if err != nil {
		return nil, "", err
	}
	return sess, c.SSMClient.Endpoint, nil
}

// Delete Session
func (c *Client) DeleteSession(sessionID string) error {
	subctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if _, err := c.SSMClient.TerminateSessionWithContext(subctx, &ssm.TerminateSessionInput{SessionId: &sessionID}); err != nil {
		return err
	}
	return nil
}

func (c *Client) SendCommand(input *ssm.SendCommandInput) (*ssm.SendCommandOutput, error) {
	sess, err := c.SSMClient.SendCommand(input)
	if err != nil {
		return nil, err
	}
	return sess, nil
}
