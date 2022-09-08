//go:build !darwin
// +build !darwin

package keymanager

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/u-cto-devops/lguctl/pkg/schema"
	"io"
)

type KeyChain struct {
	Config
}

// NewKeyChain creates a keychain
func NewKeyChain() (*KeyChain, error) {
	return nil, fmt.Errorf("in order to support MacOS Key Chain, you need to build on MacOS")
}

// Generate creates new executable credentials
func (k *KeyChain) Generate(profile string, config *schema.Config) (*credentials.Credentials, error) {
	return nil, fmt.Errorf("in order to support MacOS Key Chain, you need to build on MacOS")
}

// Status shows the current status of keychain list
func (k *KeyChain) Status(out io.Writer) error {
	return fmt.Errorf("in order to support MacOS Key Chain, you need to build on MacOS")
}

// Add adds credentials to keychain
func (k *KeyChain) Add() error {
	return fmt.Errorf("in order to support MacOS Key Chain, you need to build on MacOS")
}
