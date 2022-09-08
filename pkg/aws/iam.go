package aws

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/sirupsen/logrus"

	"github.com/u-cto-devops/lguctl/pkg/color"
	"github.com/u-cto-devops/lguctl/pkg/tools"
)

type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

type StatementEntry struct {
	Sid      string
	Effect   string
	Action   interface{}
	Resource interface{}
}

func GetIAMClientFn(sess client.ConfigProvider, creds *credentials.Credentials) *iam.IAM {
	if creds == nil {
		return iam.New(sess)
	}
	return iam.New(sess, aws.NewConfig().WithCredentials(creds))
}

// CreateNewCredentials creates new ACCESS_KEY, SECRET_ACCESS_KEY
func (c Client) CreateNewCredentials(name string) (*iam.AccessKey, error) {
	input := &iam.CreateAccessKeyInput{
		UserName: aws.String(name),
	}

	result, err := c.IAMClient.CreateAccessKey(input)
	if err != nil {
		return nil, err
	}

	return result.AccessKey, nil
}

// DeleteAccessKey deletes access key of user
func (c Client) DeleteAccessKey(accessKey, userName string) error {
	input := &iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(accessKey),
		UserName:    aws.String(userName),
	}

	_, err := c.IAMClient.DeleteAccessKey(input)
	if err != nil {
		return err
	}

	return nil
}

// CheckAccessKeyExpired check whether access key is expired or not
func (c Client) CheckAccessKeyExpired(name, accessKeyID string) error {
	keys, err := c.GetAccessKeyList(name)
	if err != nil {
		return err
	}

	for _, key := range keys {
		if *key.AccessKeyId == accessKeyID {
			if tools.IsExpired(*key.CreateDate, 24*180*time.Hour) {
				return errors.New("your access key is expired. please renew by running `lguctl renew-credential`")
			}

			logrus.Debugf("your access key is not expired")
			return nil
		}
	}
	return errors.New("your access key configuration is wrong")
}

// GetAccessKeyList lists all access key of user
func (c Client) GetAccessKeyList(name string) ([]*iam.AccessKeyMetadata, error) {
	if err := tools.ClearOsEnv(); err != nil {
		return nil, nil
	}
	svc := iam.New(GetAwsSession())
	input := &iam.ListAccessKeysInput{
		UserName: aws.String(name),
	}

	result, err := svc.ListAccessKeys(input)
	if err != nil {
		return nil, err
	}

	return result.AccessKeyMetadata, nil
}

// GetListGroupForUser lists all group of user
func (c Client) GetListGroupForUser(out io.Writer, userName string) error {
	svc := iam.New(GetAwsSession())

	groupNames, err := c.getGroupName(svc, userName)

	if err != nil {
		return err
	}

	for _, groupName := range groupNames {
		color.Green.Fprintf(out, groupName)
	}

	return nil
}

func (c Client) getGroupName(svc *iam.IAM, name string) ([]string, error) {
	input := &iam.ListGroupsForUserInput{
		UserName: aws.String(name),
	}

	result, err := svc.ListGroupsForUser(input)

	if err != nil {
		return nil, err
	}

	groupNames := make([]string, 0)

	for _, groupName := range result.Groups {
		groupNames = append(groupNames, *groupName.GroupName)
	}

	return groupNames, nil
}

// GetListPolicyAttachedGroup all group of user
func (c Client) GetListPolicyAttachedGroup(out io.Writer, userName string) error {
	svc := iam.New(GetAwsSession())

	groupNames, err := c.getGroupName(svc, userName)

	if err != nil {
		return err
	}

	var groupPolicyList []*iam.AttachedPolicy
	for _, groupName := range groupNames {
		policyList, err := c.getPoliciesList(svc, groupName)
		if err != nil {
			return err
		}
		groupPolicyList = append(groupPolicyList, policyList.AttachedPolicies...)
	}

	groupPolicyList = makeSliceUnique(groupPolicyList)

	for _, policy := range groupPolicyList {
		color.Green.Fprintf(out, *policy.PolicyName)
	}

	return nil
}

func (c Client) getPoliciesList(svc *iam.IAM, groupName string) (*iam.ListAttachedGroupPoliciesOutput, error) {
	input := &iam.ListAttachedGroupPoliciesInput{
		GroupName: aws.String(groupName),
	}

	result, err := svc.ListAttachedGroupPolicies(input)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) GetListRoleArn(out io.Writer, userName string) error {
	svc := iam.New(GetAwsSession())

	groupNames, err := c.getGroupName(svc, userName)

	if err != nil {
		return err
	}

	var groupPolicyList []*iam.AttachedPolicy
	for _, groupName := range groupNames {
		policyList, err := c.getPoliciesList(svc, groupName)
		if err != nil {
			return err
		}
		groupPolicyList = append(groupPolicyList, policyList.AttachedPolicies...)
	}

	groupPolicyList = makeSliceUnique(groupPolicyList)

	for _, policy := range groupPolicyList {
		if !strings.Contains(*policy.PolicyName, "Billing") && !strings.Contains(*policy.PolicyName, "SelfControlMFA") {
			input := &iam.GetPolicyVersionInput{
				PolicyArn: aws.String(*policy.PolicyArn),
				VersionId: aws.String("v1"),
			}
			result, err := svc.GetPolicyVersion(input)

			if err != nil {
				return err
			}
			decodedValue, err := url.QueryUnescape(aws.StringValue(result.PolicyVersion.Document))
			if err != nil {
				return err
			}
			err = convertPolicy(out, decodedValue)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func convertPolicy(out io.Writer, data string) error {
	var doc PolicyDocument
	err := json.Unmarshal([]byte(data), &doc)

	if err != nil {
		return err
	}
	//find out the type of resource string or []string
	for _, statement := range doc.Statement {
		if reflect.TypeOf(statement.Resource).Name() != "string" {
			// we will convert the []interface to []string
			x := statement.Resource.([]interface{})
			y := make([]string, len(x))
			for i := 0; i < len(x); i++ {
				y[i] = x[i].(string)
			}
			statement.Resource = y
		}
	}

	assumeRole := fmt.Sprintf("%v", doc.Statement[0].Resource)
	roleName := strings.Split(assumeRole, "/")
	var b bytes.Buffer
	b.WriteString(roleName[1])
	color.Green.Fprintf(out, "%s : %s", b.String(), assumeRole)
	return nil
}

func makeSliceUnique(policyList []*iam.AttachedPolicy) []*iam.AttachedPolicy {
	keys := make(map[string]struct{})
	result := make([]*iam.AttachedPolicy, 0)
	for _, val := range policyList {
		if _, ok := keys[*val.PolicyName]; ok {
			continue
		} else {
			keys[*val.PolicyName] = struct{}{}
			result = append(result, val)
		}
	}
	return result
}
