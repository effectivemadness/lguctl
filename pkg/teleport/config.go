package teleport

import (
	"context"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gravitational/teleport/lib/client/identityfile"
	"github.com/gravitational/teleport/lib/utils"
	"github.com/sirupsen/logrus"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// CLIConf stores command line arguments and flags:
type CLIConf struct {
	// UserHost contains "[login]@hostname" argument to SSH command
	UserHost string
	// Commands to execute on a remote host
	RemoteCommand []string
	// DesiredRoles indicates one or more roles which should be requested.
	DesiredRoles string
	// RequestReason indicates the reason for an access request.
	RequestReason string
	// SuggestedReviewers is a list of suggested request reviewers.
	SuggestedReviewers string
	// NoWait can be used with an access request to exit without waiting for a request resolution.
	NoWait bool
	// RequestedResourceIDs is a list of resources to request access to.
	RequestedResourceIDs []string
	// RequestID is an access request ID
	RequestID string
	// RequestIDs is a list of access request IDs
	RequestIDs []string
	// ReviewReason indicates the reason for an access review.
	ReviewReason string
	// ReviewableRequests indicates that only requests which can be reviewed should
	// be listed.
	ReviewableRequests bool
	// SuggestedRequests indicates that only requests which suggest the current user
	// as a reviewer should be listed.
	SuggestedRequests bool
	// MyRequests indicates that only requests created by the current user
	// should be listed.
	MyRequests bool
	// Approve/Deny indicates the desired review kind.
	Approve, Deny bool
	// ResourceKind is the resource kind to search for
	ResourceKind string
	// Username is the Teleport user's username (to login into proxies)
	Username string
	// ExplicitUsername is true if Username was initially set by the end-user
	// (for example, using command-line flags).
	ExplicitUsername bool
	// Proxy keeps the hostname:port of the SSH proxy to use
	Proxy string
	// TTL defines how long a session must be active (in minutes)
	MinsToLive int32
	// SSH Port on a remote SSH host
	NodePort int32
	// Login on a remote SSH host
	NodeLogin string
	// InsecureSkipVerify bypasses verification of HTTPS certificate when talking to web proxy
	InsecureSkipVerify bool
	// SessionID identifies the session tsh is operating on.
	// For `tsh join`, it is the ID of the session to join.
	// For `tsh play`, it is either the ID of the session to play,
	// or the path to a local session file which has already been
	// downloaded.
	SessionID string
	// Src:dest parameter for SCP
	CopySpec []string
	// -r flag for scp
	RecursiveCopy bool
	// -L flag for ssh. Local port forwarding like 'ssh -L 80:remote.host:80 -L 443:remote.host:443'
	LocalForwardPorts []string
	// DynamicForwardedPorts is port forwarding using SOCKS5. It is similar to
	// "ssh -D 8080 example.com".
	DynamicForwardedPorts []string
	// ForwardAgent agent to target node. Equivalent of -A for OpenSSH.
	ForwardAgent bool
	// ProxyJump is an optional -J flag pointing to the list of jumphosts,
	// it is an equivalent of --proxy flag in tsh interpretation
	ProxyJump string
	// --local flag for ssh
	LocalExec bool
	// SiteName specifies remote site go login to
	SiteName string
	// KubernetesCluster specifies the kubernetes cluster to login to.
	KubernetesCluster string
	// DaemonAddr is the daemon listening address.
	DaemonAddr string
	// DaemonCertsDir is the directory containing certs used to create secure gRPC connection with daemon service
	DaemonCertsDir string
	// DatabaseService specifies the database proxy server to log into.
	DatabaseService string
	// DatabaseUser specifies database user to embed in the certificate.
	DatabaseUser string
	// DatabaseName specifies database name to embed in the certificate.
	DatabaseName string
	// AppName specifies proxied application name.
	AppName string
	// Interactive, when set to true, launches remote command with the terminal attached
	Interactive bool
	// Quiet mode, -q command (disables progress printing)
	Quiet bool
	// Namespace is used to select cluster namespace
	Namespace string
	// NoCache is used to turn off client cache for nodes discovery
	NoCache bool
	// BenchDuration is a duration for the benchmark
	BenchDuration time.Duration
	// BenchRate is a requests per second rate to maintain
	BenchRate int
	// BenchInteractive indicates that we should create interactive session
	BenchInteractive bool
	// BenchExport exports the latency profile
	BenchExport bool
	// BenchExportPath saves the latency profile in provided path
	BenchExportPath string
	// BenchTicks ticks per half distance
	BenchTicks int32
	// BenchValueScale value at which to scale the values recorded
	BenchValueScale float64
	// Context is a context to control execution
	Context context.Context
	// IdentityFileIn is an argument to -i flag (path to the private key+cert file)
	IdentityFileIn string
	// Compatibility flags, --compat, specifies OpenSSH compatibility flags.
	Compatibility string
	// CertificateFormat defines the format of the user SSH certificate.
	CertificateFormat string
	// IdentityFileOut is an argument to -out flag
	IdentityFileOut string
	// IdentityFormat (used for --format flag for 'tsh login') defines which
	// format to use with --out to store a freshly retrieved certificate
	IdentityFormat identityfile.Format
	// IdentityOverwrite when true will overwrite any existing identity file at
	// IdentityFileOut. When false, user will be prompted before overwriting
	// any files.
	IdentityOverwrite bool

	// BindAddr is an address in the form of host:port to bind to
	// during `tsh login` command
	BindAddr string

	// AuthConnector is the name of the connector to use.
	AuthConnector string

	// MFAMode is the preferred mode for MFA/Passwordless assertions.
	MFAMode string

	// SkipVersionCheck skips version checking for client and server
	SkipVersionCheck bool

	// Options is a list of OpenSSH options in the format used in the
	// configuration file.
	Options []string

	// Verbose is used to print extra output.
	Verbose bool

	// Format is used to change the format of output
	Format string

	// SearchKeywords is a list of search keywords to match against resource field values.
	SearchKeywords string

	// PredicateExpression defines boolean conditions that will be matched against the resource.
	PredicateExpression string

	// NoRemoteExec will not execute a remote command after connecting to a host,
	// will block instead. Useful when port forwarding. Equivalent of -N for OpenSSH.
	NoRemoteExec bool

	// X11ForwardingUntrusted will set up untrusted X11 forwarding for the session ('ssh -X')
	X11ForwardingUntrusted bool

	// X11Forwarding will set up trusted X11 forwarding for the session ('ssh -Y')
	X11ForwardingTrusted bool

	// X11ForwardingTimeout can optionally set to set a timeout for untrusted X11 forwarding.
	X11ForwardingTimeout time.Duration

	// Debug sends debug logs to stdout.
	Debug bool

	// Browser can be used to pass the name of a browser to override the system default
	// (not currently implemented), or set to 'none' to suppress browser opening entirely.
	Browser string

	// UseLocalSSHAgent set to false will prevent this client from attempting to
	// connect to the local ssh-agent (or similar) socket at $SSH_AUTH_SOCK.
	//
	// Deprecated in favor of `AddKeysToAgent`.
	UseLocalSSHAgent bool

	// AddKeysToAgent specifies the behavior of how certs are handled.
	AddKeysToAgent string

	// EnableEscapeSequences will scan stdin for SSH escape sequences during
	// command/shell execution. This also requires stdin to be an interactive
	// terminal.
	EnableEscapeSequences bool

	// PreserveAttrs preserves access/modification times from the original file.
	PreserveAttrs bool

	// executablePath is the absolute path to the current executable.
	executablePath string

	// unsetEnvironment unsets Teleport related environment variables.
	unsetEnvironment bool

	// overrideStdout allows to switch standard output source for resource command. Used in tests.
	overrideStdout io.Writer
	// overrideStderr allows to switch standard error source for resource command. Used in tests.
	overrideStderr io.Writer

	// mockSSOLogin used in tests to override sso login handler in teleport client.
	// mockSSOLogin client.SSOLoginFunc

	// HomePath is where tsh stores profiles
	HomePath string

	// GlobalTshConfigPath is a path to global TSH config. Can be overridden with TELEPORT_GLOBAL_TSH_CONFIG.
	GlobalTshConfigPath string

	// LocalProxyPort is a port used by local proxy listener.
	LocalProxyPort string
	// LocalProxyCertFile is the client certificate used by local proxy.
	LocalProxyCertFile string
	// LocalProxyKeyFile is the client key used by local proxy.
	LocalProxyKeyFile string
	// LocalProxyTunnel specifies whether local proxy will open auth'd tunnel.
	LocalProxyTunnel bool

	// AWSRole is Amazon Role ARN or role name that will be used for AWS CLI access.
	AWSRole string
	// AWSCommandArgs contains arguments that will be forwarded to AWS CLI binary.
	AWSCommandArgs []string
	// AWSEndpointURLMode is an AWS proxy mode that serves an AWS endpoint URL
	// proxy instead of an HTTPS proxy.
	AWSEndpointURLMode bool

	// Reason is the reason for starting an ssh or kube session.
	Reason string

	// Invited is a list of invited users to an ssh or kube session.
	Invited []string

	// JoinMode is the participant mode someone is joining a session as.
	JoinMode string

	// displayParticipantRequirements is set if verbose participant requirement information should be printed for moderated sessions.
	displayParticipantRequirements bool

	// TshConfig is the loaded tsh configuration file ~/.tsh/config/config.yaml.
	// TshConfig TshConfig

	// ListAll specifies if an ls command should return results from all clusters and proxies.
	ListAll bool
	// SampleTraces indicates whether traces should be sampled.
	SampleTraces bool

	// TracingProvider is the provider to use to create tracers, from which spans can be created.
	TracingProvider oteltrace.TracerProvider

	// disableAccessRequest disables automatic resource access requests.
	disableAccessRequest bool

	// FromUTC is the start time to use for the range of sessions listed by the session recordings listing command
	FromUTC string

	// ToUTC is the start time to use for the range of sessions listed by the session recordings listing command
	ToUTC string

	// maxRecordingsToShow is the maximum number of session recordings to show per page of results
	maxRecordingsToShow int

	// recordingsSince is a duration which sets the time into the past in which to list session recordings
	recordingsSince string

	// command is the selected command (and subcommands) parsed from command
	// line args. Note that this command does not contain the binary (e.g. tsh).
	command string

	// cmdRunner is a custom function to execute provided exec.Cmd. Mainly used
	// in testing.
	cmdRunner func(*exec.Cmd) error
}

// CreateConfig creates a configuration
func CreateConfig() CLIConf {
	var cf CLIConf
	utils.InitLogger(utils.LoggingForCLI, logrus.WarnLevel)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		exitSignals := make(chan os.Signal, 1)
		signal.Notify(exitSignals, syscall.SIGTERM, syscall.SIGINT)

		sig := <-exitSignals
		log.Debugf("signal: %v", sig)
		cancel()
	}()
	cf.Context = ctx

	defaultIdentityFormat := identityfile.DefaultFormat
	cf.IdentityFormat = defaultIdentityFormat

	return cf
}
