package teleport

import (
	"fmt"
	"strings"

	"github.com/gravitational/teleport/lib/utils"
	"github.com/gravitational/trace"
)

// AllOptions is a listing of all known OpenSSH options.
var AllOptions = map[string]map[string]bool{
	"AddKeysToAgent":                   {"yes": true},
	"AddressFamily":                    {},
	"BatchMode":                        {},
	"BindAddress":                      {},
	"CanonicalDomains":                 {},
	"CanonicalizeFallbackLocal":        {},
	"CanonicalizeHostname":             {},
	"CanonicalizeMaxDots":              {},
	"CanonicalizePermittedCNAMEs":      {},
	"CertificateFile":                  {},
	"ChallengeResponseAuthentication":  {},
	"CheckHostIP":                      {},
	"Cipher":                           {},
	"Ciphers":                          {},
	"ClearAllForwardings":              {},
	"Compression":                      {},
	"CompressionLevel":                 {},
	"ConnectionAttempts":               {},
	"ConnectTimeout":                   {},
	"ControlMaster":                    {},
	"ControlPath":                      {},
	"ControlPersist":                   {},
	"DynamicForward":                   {},
	"EscapeChar":                       {},
	"ExitOnForwardFailure":             {},
	"FingerprintHash":                  {},
	"ForwardAgent":                     {"yes": true, "no": true},
	"ForwardX11":                       {},
	"ForwardX11Timeout":                {},
	"ForwardX11Trusted":                {},
	"GatewayPorts":                     {},
	"GlobalKnownHostsFile":             {},
	"GSSAPIAuthentication":             {},
	"GSSAPIDelegateCredentials":        {},
	"HashKnownHosts":                   {},
	"Host":                             {},
	"HostbasedAuthentication":          {},
	"HostbasedKeyTypes":                {},
	"HostKeyAlgorithms":                {},
	"HostKeyAlias":                     {},
	"HostName":                         {},
	"IdentityFile":                     {},
	"IdentitiesOnly":                   {},
	"IPQoS":                            {},
	"KbdInteractiveAuthentication":     {},
	"KbdInteractiveDevices":            {},
	"KexAlgorithms":                    {},
	"LocalCommand":                     {},
	"LocalForward":                     {},
	"LogLevel":                         {},
	"MACs":                             {},
	"Match":                            {},
	"NoHostAuthenticationForLocalhost": {},
	"NumberOfPasswordPrompts":          {},
	"PasswordAuthentication":           {},
	"PermitLocalCommand":               {},
	"PKCS11Provider":                   {},
	"Port":                             {},
	"PreferredAuthentications":         {},
	"Protocol":                         {},
	"ProxyCommand":                     {},
	"ProxyUseFdpass":                   {},
	"PubkeyAcceptedKeyTypes":           {},
	"PubkeyAuthentication":             {},
	"RekeyLimit":                       {},
	"RemoteForward":                    {},
	"RequestTTY":                       {"yes": true, "no": true},
	"RhostsRSAAuthentication":          {},
	"RSAAuthentication":                {},
	"SendEnv":                          {},
	"ServerAliveInterval":              {},
	"ServerAliveCountMax":              {},
	"StreamLocalBindMask":              {},
	"StreamLocalBindUnlink":            {},
	"StrictHostKeyChecking":            {"yes": true, "no": true},
	"TCPKeepAlive":                     {},
	"Tunnel":                           {},
	"TunnelDevice":                     {},
	"UpdateHostKeys":                   {},
	"UsePrivilegedPort":                {},
	"User":                             {},
	"UserKnownHostsFile":               {},
	"VerifyHostKeyDNS":                 {},
	"VisualHostKey":                    {},
	"XAuthLocation":                    {},
}

// Options holds parsed values of OpenSSH options.
type Options struct {
	// AddKeysToAgent specifies whether keys should be automatically added to a
	// running SSH agent. Supported options values are "yes".
	AddKeysToAgent bool

	// ForwardAgent specifies whether the connection to the authentication
	// agent will be forwarded to the remote machine. Supported option values
	// are "yes" and "no".
	ForwardAgent bool

	// RequestTTY specifies whether to request a pseudo-tty for the session.
	// Supported option values are "yes" and "no".
	RequestTTY bool

	// StrictHostKeyChecking is used control if tsh will automatically add host
	// keys to the ~/.tsh/known_hosts file. Supported option values are "yes"
	// and "no".
	StrictHostKeyChecking bool
}

func parseOptions(opts []string) (Options, error) {
	// By default, Teleport prefers strict host key checking and adding keys
	// to system SSH agent.
	options := Options{
		StrictHostKeyChecking: true,
		AddKeysToAgent:        true,
	}

	for _, o := range opts {
		key, value, err := splitOption(o)
		if err != nil {
			return Options{}, trace.Wrap(err)
		}

		supportedValues, ok := AllOptions[key]
		if !ok {
			return Options{}, trace.BadParameter("unsupported option key: %v", key)
		}

		if len(supportedValues) == 0 {
			fmt.Printf("WARNING: Option '%v' is not supported.\n", key)
			continue
		}

		_, ok = supportedValues[value]
		if !ok {
			return Options{}, trace.BadParameter("unsupported option value: %v", value)
		}

		switch key {
		case "AddKeysToAgent":
			options.AddKeysToAgent = utils.AsBool(value)
		case "ForwardAgent":
			options.ForwardAgent = utils.AsBool(value)
		case "RequestTTY":
			options.RequestTTY = utils.AsBool(value)
		case "StrictHostKeyChecking":
			options.StrictHostKeyChecking = utils.AsBool(value)
		}
	}

	return options, nil
}

func splitOption(option string) (string, string, error) {
	parts := strings.FieldsFunc(option, fieldsFunc)

	if len(parts) != 2 {
		return "", "", trace.BadParameter("invalid format for option")
	}

	return parts[0], parts[1], nil
}

// fieldsFunc splits key-value pairs off ' ' and '='.
func fieldsFunc(c rune) bool {
	switch {
	case c == ' ':
		return true
	case c == '=':
		return true
	default:
		return false
	}
}
