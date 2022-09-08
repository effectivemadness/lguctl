package keymanager

import (
	"fmt"
	"io"
	"runtime"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/u-cto-devops/lguctl/pkg/constants"
	"github.com/u-cto-devops/lguctl/pkg/schema"
)

type KeyManager interface {
	Generate(string, *schema.Config) (*credentials.Credentials, error)
	Status(io.Writer) error
	Add() error
}

// GenerateKeyManager returns key manager according to key type
func GenerateKeyManager(keyType string) (KeyManager, error) {
	if keyType == constants.DefaultKeyType && runtime.GOOS == "darwin" {
		return NewKeyChain()
	}

	return nil, fmt.Errorf("no manager exist")
}
