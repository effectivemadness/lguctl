package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/u-cto-devops/lguctl/pkg/constants"
)

func GetEC2ClientFn(sess client.ConfigProvider, region string, creds *credentials.Credentials) *ec2.EC2 {
	if creds == nil {
		return ec2.New(sess, &aws.Config{Region: aws.String(region)})
	}
	return ec2.New(sess, &aws.Config{Region: aws.String(region), Credentials: creds})
}

// GetInstanceList returns list of instances
func (c *Client) GetInstanceList(ret []*ec2.Instance, nextToken *string) ([]*ec2.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		NextToken: nextToken,
	}

	result, err := c.EC2Client.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	for _, reservation := range result.Reservations {
		ret = append(ret, reservation.Instances...)
	}

	if result.NextToken != nil {
		return c.GetInstanceList(ret, result.NextToken)
	}

	return ret, nil
}

// GetInstanceIds retrieves only ids from list
func (c *Client) GetInstanceIds(instanceList []*ec2.Instance) ([]string, error) {
	var ids []string
	for _, instance := range instanceList {
		ids = append(ids, *instance.InstanceId)
	}

	return ids, nil
}

// GetInstanceListOnlyIds retrieves only ids from list
func (c *Client) GetInstanceListOnlyIds(ret []string, nextToken *string) ([]string, error) {
	input := &ec2.DescribeInstancesInput{
		NextToken: nextToken,
	}

	result, err := c.EC2Client.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	var nameTag string
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			nameTag = "Unknown"
			for _, tag := range instance.Tags {
				if strings.ToLower(*tag.Key) == "name" {
					nameTag = *tag.Value
				}
			}
			ret = append(ret, makeInstanceName(*instance.InstanceId, nameTag))
		}
	}

	if result.NextToken != nil {
		return c.GetInstanceListOnlyIds(ret, result.NextToken)
	}

	return ret, nil
}

// makeInstanceName makes instance name
func makeInstanceName(id, name string) string {
	return fmt.Sprintf("%s%s%s", name, constants.NameDelimiter, id)
}
