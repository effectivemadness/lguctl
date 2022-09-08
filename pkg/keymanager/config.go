package keymanager

import "github.com/u-cto-devops/lguctl/pkg/constants"

type Config struct {
	Path        string
	ServiceName string
	Account     string
	Label       string
	Description string
}

// DefaultConfig returns default configuration for keychain
func DefaultConfig() Config {
	return Config{
		Path:        constants.DefaultKeyChainPath,
		ServiceName: constants.ServiceName,
		Account:     constants.DefaultKeyChainAccount,
	}
}
