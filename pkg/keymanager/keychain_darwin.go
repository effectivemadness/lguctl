package keymanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	kclib "github.com/keybase/go-keychain"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"

	"github.com/u-cto-devops/lguctl/pkg/color"
	"github.com/u-cto-devops/lguctl/pkg/constants"
	"github.com/u-cto-devops/lguctl/pkg/provider"
	"github.com/u-cto-devops/lguctl/pkg/schema"
	"github.com/u-cto-devops/lguctl/pkg/tools"
)

type KeyChain struct {
	Config
}

// NewKeyChain creates a keychain
func NewKeyChain() (*KeyChain, error) {
	kc := KeyChain{DefaultConfig()}

	return &kc, nil
}

// AddCredentials add credentials to keychain
func (k *KeyChain) AddCredentials(credential Credential) error {
	var kc kclib.Keychain

	// when we are setting a value, we create or open
	if len(k.Config.Path) > 0 {
		var err error
		kc, err = k.CreateOrOpenKeychain()
		if err != nil {
			return err
		}
	}

	bytes, err := json.Marshal(credential.Credentials)
	if err != nil {
		return err
	}

	item := k.GenerateBaseItem()
	item.SetLabel(k.Config.Label)
	item.SetDescription(k.Config.Description)
	item.SetAccessible(kclib.AccessibleWhenUnlocked)
	item.SetData(bytes)

	if len(k.Config.Path) > 0 {
		item.UseKeychain(kc)
	}

	err = kclib.AddItem(item)

	if err == kclib.ErrorDuplicateItem {
		logrus.Debugf("Item already exists, updating")
		err = k.UpdateCredential(kc, item)
	}

	if err != nil {
		return err
	}

	return nil
}

// CreateOrOpenKeychain creates or opens macOS Keychain
func (k *KeyChain) CreateOrOpenKeychain() (kclib.Keychain, error) {
	kc := kclib.NewWithPath(k.Config.Path)

	logrus.Debugln("Checking keychain status")
	err := kc.Status()
	if err == nil {
		logrus.Debugln("Keychain status returned nil, keychain exists")
		return kc, nil
	}

	logrus.Debugf("Keychain status returned error: %v", err)

	if err != kclib.ErrorNoSuchKeychain {
		return kclib.Keychain{}, err
	}

	logrus.Debugf("Creating keychain %s with prompt", k.Config.Path)
	return kclib.NewKeychainWithPrompt(k.Config.Path)
}

// UpdateCredential updates credentials
func (k *KeyChain) UpdateCredential(kc kclib.Keychain, kcItem kclib.Item) error {
	queryItem := k.GenerateBaseItem()
	queryItem.SetMatchLimit(kclib.MatchLimitOne)
	queryItem.SetReturnAttributes(true)

	if len(k.Config.Path) > 0 {
		queryItem.SetMatchSearchList(kc)
	}

	results, err := kclib.QueryItem(queryItem)
	if err != nil {
		return fmt.Errorf("failed to query keychain: %v", err)
	}

	if len(results) == 0 {
		return errors.New("no results")
	}

	if err := kclib.UpdateItem(queryItem, kcItem); err != nil {
		return fmt.Errorf("failed to update item in keychain: %v", err)
	}

	return nil
}

// GetCredential gets credential from keychain
func (k *KeyChain) GetCredential(key string) (Credential, error) {
	query := k.GenerateBaseItem()
	query.SetMatchLimit(kclib.MatchLimitOne)
	query.SetReturnAttributes(true)
	query.SetReturnData(true)

	if len(k.Config.Path) > 0 {
		// When we are querying, we don't create by default
		query.SetMatchSearchList(kclib.NewWithPath(k.Config.Path))
	}

	logrus.Debugf("Querying keychain for service=%q, account=%q, keychain=%q", k.Config.ServiceName, key, k.Config.Path)
	results, err := kclib.QueryItem(query)
	if err == kclib.ErrorItemNotFound || len(results) == 0 {
		logrus.Debugln("No results found")
		return Credential{}, kclib.ErrorItemNotFound
	}

	if err != nil {
		logrus.Debugf("error: %#v", err)
		return Credential{}, err
	}

	var creds Credential
	if err := json.Unmarshal(results[0].Data, &creds); err != nil {
		return creds, err
	}

	creds.Label = results[0].Label
	creds.CreatedTime = results[0].CreationDate
	creds.ModificationTime = results[0].ModificationDate

	return creds, nil
}

// CheckEmptyError checks if error occurred because there is no key in keychain
func (k *KeyChain) CheckEmptyError(err error) bool {
	return err == kclib.ErrorItemNotFound
}

// ShowCredentialsStatus shows credentials status
func (k *KeyChain) ShowCredentialsStatus(out io.Writer, credentials []Credential) {
	if len(credentials) == 0 {
		color.Red.Fprintln(out, "no credential exists")
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Account", "Access Key ID", "Label", "Created Date", "modified Date"})

	data := makeTableRecords(credentials)

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}

// makeTableRecords makes credential data to records
func makeTableRecords(credentials []Credential) [][]string {
	var data [][]string
	for _, cred := range credentials {
		data = append(data, []string{
			constants.DefaultKeyChainAccount, tools.FormatKeyForDisplay(cred.AccessKeyID), cred.Label, cred.CreatedTime.String(), cred.ModificationTime.String(),
		})
	}

	return data
}

// GenerateBaseItem generates base item for keychain function
func (k *KeyChain) GenerateBaseItem() kclib.Item {
	item := kclib.NewItem()
	item.SetSecClass(kclib.SecClassGenericPassword)
	item.SetService(k.Config.ServiceName)
	item.SetAccount(k.Config.Account)
	item.SetAccessible(kclib.AccessibleWhenUnlocked)

	return item
}

// Generate creates new executable credentials
func (k *KeyChain) Generate(profile string, config *schema.Config) (*credentials.Credentials, error) {
	provider, err := k.NewTemporaryProvider(profile, config)
	if err != nil {
		return nil, err
	}

	return credentials.NewCredentials(provider), nil
}

// Status shows the current status of keychain list
func (k *KeyChain) Status(out io.Writer) error {
	creds, err := k.GetCredential(constants.DefaultKeyChainAccount)
	if err != nil {
		if k.CheckEmptyError(err) {
			k.ShowCredentialsStatus(out, nil)
		}
		return err
	}

	credsList := []Credential{creds}
	k.ShowCredentialsStatus(out, credsList)

	return nil
}

// Add adds credentials to keychain
func (k *KeyChain) Add() error {
	creds, err := k.GetCredentialsFromPrompt()
	if err != nil {
		return err
	}

	if err := k.AddCredentials(Credential{Credentials: *creds}); err != nil {
		return err
	}

	return nil
}

// NewTemporaryProvider creates new temporary provider
func (k *KeyChain) NewTemporaryProvider(profile string, config *schema.Config) (credentials.Provider, error) {
	base, err := k.GetCredential(constants.DefaultKeyChainAccount)
	if err != nil {
		return nil, err
	}

	baseProvider := provider.NewDefaultProvider(profile, base.AccessKeyID, base.SecretAccessKey)

	if profile == constants.DefaultProfile {
		return baseProvider, nil
	}

	return provider.NewAssumeProvider(profile, config, credentials.NewCredentials(baseProvider)), nil
}
