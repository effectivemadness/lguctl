package teleport

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gravitational/teleport"
	"github.com/gravitational/teleport/api/client"
	"github.com/gravitational/teleport/api/types"

	"github.com/gravitational/trace"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{
	trace.Component: teleport.ComponentTSH,
})

const (
	clusterEnvVar = "TELEPORT_SITE"
)

type Teleport struct {
	Config CLIConf
}

func NewTeleport() *Teleport {
	return &Teleport{
		Config: CreateConfig(),
	}
}

func (t *Teleport) ConnectToNode() error {
	onSSH(&t.Config)
	return nil
}

func (t *Teleport) DescribeStatus() error {
	return onStatus(&t.Config)
}

func (t *Teleport) Login() error {
	return onLogin(&t.Config)
}

func (t *Teleport) ListNode() error {
	return onListNodes(&t.Config)
}

func makeClient(ctx context.Context) (*client.Client, error) {
	clt, err := client.New(ctx, client.Config{
		// Multiple Credentials can be provided to attempt to authenticate
		// the client. At least one Credentials object must be provided.
		Credentials: []client.Credentials{
			client.LoadProfile("", ""),
		},
		// set to true if your web proxy doesn't have HTTP/TLS certificate
		// configured yet (never use this in production).
		InsecureAddressDiscovery: false,
	})
	if err != nil {
		log.Fatal(err)
	}
	return clt, err
}

// SelectNode select node
func (t *Teleport) SelectNode() (string, error) {
	nodes, err := getNodeList(&t.Config)
	if err != nil {
		return "", err
	}

	var options []string
	for _, n := range nodes {
		options = append(options, n.GetHostname())
	}

	var selected string
	prompt := &survey.Select{
		Message: "Choose node: ",
		Options: options,
	}
	survey.AskOne(prompt, &selected)

	if len(selected) == 0 {
		return "", fmt.Errorf("selecting node has been canceled")
	}

	return selected, nil
}

// getNodeList retrieves all nodes of cluster
func getNodeList(cf *CLIConf) ([]types.Server, error) {
	ctx := context.Background()
	clt, err := makeClient(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, trace.Wrap(err)
	}
	defer clt.Close()

	// Get list of all nodes in backend and sort by "Node Name".
	nodes, err := clt.GetNodes(ctx, "default")

	if err != nil {
		return nil, trace.Wrap(err)
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].GetHostname() < nodes[j].GetHostname()
	})

	return nodes, nil
}

func onListNodes(cf *CLIConf) error {
	nodes, err := getNodeList(cf)
	if err != nil {
		return trace.Wrap(err)
	}

	for _, node := range nodes {
		fmt.Println(node.GetHostname())
	}

	return nil
}

func onLogin(cf *CLIConf) error {
	cmd := exec.Command("tsh", "-d", "login", "--proxy", cf.Proxy, "--insecure")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}
	onStatus(cf)
	return nil
}

func onStatus(cf *CLIConf) error {
	cmd := exec.Command("tsh", "status")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Print(string(stdout))
	return nil
}

func onSSH(cf *CLIConf) error {
	cmd := exec.Command("tsh", "ssh", cf.UserHost)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run() // add error checking
	return nil
}

// // onListNodes executes 'tsh ls' command.
// func onListNodes(cf *CLIConf) error {
// 	nodes, err := getNodeList(cf)
// 	if err != nil {
// 		return err
// 	}

// 	printNodesAsText(nodes, false)
// 	return nil
// }

// func onLogin(cf *CLIConf) error {
// 	var (
// 		err error
// 		tc  *client.TeleportClient
// 		key *client.Key
// 	)

// 	forceLogin := viper.GetBool("force")

// 	// populate cluster name from environment variables
// 	// only if not set by argument (that does not support env variables)
// 	clusterName := os.Getenv(clusterEnvVar)
// 	if cf.SiteName == "" {
// 		cf.SiteName = clusterName
// 	}

// 	if cf.IdentityFileIn != "" {
// 		utils.FatalError(trace.BadParameter("-i flag cannot be used here"))
// 	}

// 	switch cf.IdentityFormat {
// 	case identityfile.FormatFile, identityfile.FormatOpenSSH, identityfile.FormatKubernetes:
// 	default:
// 		utils.FatalError(trace.BadParameter("invalid identity format: %s", cf.IdentityFormat))
// 	}

// 	// Get the status of the active profile ~/.tsh/profile as well as the status
// 	// of any other proxies the user is logged into.
// 	profile, profiles, err := client.Status("", cf.Proxy)
// 	if err != nil {
// 		if !trace.IsNotFound(err) {
// 			utils.FatalError(err)
// 		}
// 	}

// 	// make the teleport client and retrieve the certificate from the proxy:
// 	tc, err = makeClient(cf, true)
// 	if err != nil {
// 		utils.FatalError(err)
// 	}

// 	// client is already logged in and profile is not expired
// 	if profile != nil && !profile.IsExpired(clockwork.NewRealClock()) && !forceLogin {
// 		switch {
// 		// in case if nothing is specified, print current status
// 		case cf.Proxy == "" && cf.SiteName == "" && cf.DesiredRoles == "" && cf.IdentityFileOut == "":
// 			printProfiles(cf.Debug, profile, profiles)
// 			return nil
// 		// in case if parameters match, print current status
// 		case host(cf.Proxy) == host(profile.ProxyURL.Host) && cf.DesiredRoles == "":
// 			if err := tc.SaveProfile(profile.ProxyURL.Host, ""); err != nil {
// 				utils.FatalError(err)
// 			}
// 			printCurrentProfile(cf.Debug, profile)
// 			return nil
// 		// proxy is unspecified or the same as the currently provided proxy,
// 		// but cluster is specified, treat this as selecting a new cluster
// 		// for the same proxy
// 		case (cf.Proxy == "" || host(cf.Proxy) == host(profile.ProxyURL.Host)) && cf.SiteName != "":
// 			// trigger reissue, preserving any active requests.
// 			err = tc.ReissueUserCerts(cf.Context, client.ReissueParams{
// 				AccessRequests: profile.ActiveRequests.AccessRequests,
// 				RouteToCluster: cf.SiteName,
// 			})
// 			if err != nil {
// 				utils.FatalError(err)
// 			}
// 			if err := tc.SaveProfile("", ""); err != nil {
// 				utils.FatalError(err)
// 			}
// 			if err := kubeconfig.UpdateWithClient("", tc); err != nil {
// 				utils.FatalError(err)
// 			}
// 			onStatus(cf)
// 			return nil
// 		// proxy is unspecified or the same as the currently provided proxy,
// 		// but desired roles are specified, treat this as a privilege escalation
// 		// request for the same login session.
// 		case (cf.Proxy == "" || host(cf.Proxy) == host(profile.ProxyURL.Host)) && cf.DesiredRoles != "" && cf.IdentityFileOut == "":
// 			executeAccessRequest(cf)
// 			return nil
// 		// otherwise just passthrough to standard login
// 		default:
// 		}
// 	}

// 	if cf.Username == "" {
// 		cf.Username = tc.Username
// 	}

// 	// -i flag specified? save the retreived cert into an identity file
// 	makeIdentityFile := cf.IdentityFileOut != ""
// 	activateKey := !makeIdentityFile

// 	key, err = tc.Login(cf.Context, activateKey)
// 	if err != nil {
// 		utils.FatalError(err)
// 	}

// 	if makeIdentityFile {
// 		if err := setupNoninteractiveClient(tc, key); err != nil {
// 			utils.FatalError(err)
// 		}
// 		// key.TrustedCA at this point only has the CA of the root cluster we
// 		// logged into. We need to fetch all the CAs for leaf clusters too, to
// 		// make them available in the identity file.
// 		authorities, err := tc.GetTrustedCA(cf.Context, key.ClusterName)
// 		if err != nil {
// 			utils.FatalError(err)
// 		}
// 		key.TrustedCA = auth.AuthoritiesToTrustedCerts(authorities)

// 		filesWritten, err := identityfile.Write(cf.IdentityFileOut, key, cf.IdentityFormat, tc.KubeClusterAddr())
// 		if err != nil {
// 			utils.FatalError(err)
// 		}
// 		fmt.Printf("\nThe certificate has been written to %s\n", strings.Join(filesWritten, ","))
// 		return nil
// 	}

// 	// If the proxy is advertising that it supports Kubernetes, update kubeconfig.
// 	if tc.KubeProxyAddr != "" {
// 		if err := kubeconfig.UpdateWithClient("", tc); err != nil {
// 			utils.FatalError(err)
// 		}
// 	}

// 	// Regular login without -i flag.
// 	if err := tc.SaveProfile(key.ProxyHost, ""); err != nil {
// 		utils.FatalError(err)
// 	}

// 	// Print status to show information of the logged in user. Update the
// 	// command line flag (used to print status) for the proxy to make sure any
// 	// advertised settings are picked up.
// 	webProxyHost, _ := tc.WebProxyHostPort()
// 	cf.Proxy = webProxyHost
// 	if cf.DesiredRoles != "" {
// 		fmt.Println("") // visually separate onRequestExecute output
// 		executeAccessRequest(cf)
// 		return nil
// 	}

// 	return onStatus(cf)
// }

// // onStatus command shows which proxy the user is logged into and metadata
// // about the certificate.
// func onStatus(cf *CLIConf) error {
// 	// Get the status of the active profile as well as the status
// 	// of any other proxies the user is logged into.
// 	profile, profiles, err := client.Status("", cf.Proxy)
// 	if err != nil {
// 		if trace.IsNotFound(err) {
// 			fmt.Printf("Not logged in.\n")
// 			return nil
// 		}
// 		return trace.Wrap(err)
// 	}

// 	printProfiles(cf.Debug, profile, profiles)
// 	return nil
// }

// onSSH executes 'tsh ssh' command
// func onSSH(cf *CLIConf) error {
// 	tc, err := makeClient(cf, false)
// 	if err != nil {
// 		return trace.Wrap(err)
// 	}

// 	tc.Stdin = os.Stdin
// 	err = client.RetryWithRelogin(cf.Context, tc, func() error {
// 		return tc.SSH(cf.Context, cf.RemoteCommand, cf.LocalExec)
// 	})

// 	if err != nil {
// 		if strings.Contains(utils.UserMessageFromError(err), teleport.NodeIsAmbiguous) {
// 			allNodes, err := tc.ListAllNodes(cf.Context)
// 			if err != nil {
// 				return trace.Wrap(err)
// 			}
// 			var nodes []services.Server
// 			for _, node := range allNodes {
// 				if node.GetHostname() == tc.Host {
// 					nodes = append(nodes, node)
// 				}
// 			}
// 			fmt.Fprintf(os.Stderr, "error: ambiguous host could match multiple nodes\n\n")
// 			printNodesAsText(nodes, true)
// 			fmt.Fprintf(os.Stderr, "Hint: try addressing the node by unique id (ex: tsh ssh user@node-id)\n")
// 			fmt.Fprintf(os.Stderr, "Hint: use 'tsh ls -v' to list all nodes with their unique ids\n")
// 			fmt.Fprintf(os.Stderr, "\n")
// 			os.Exit(1)
// 		}
// 		// exit with the same exit status as the failed command:
// 		if tc.ExitStatus != 0 {
// 			fmt.Fprintln(os.Stderr, utils.UserMessageFromError(err))
// 			os.Exit(tc.ExitStatus)
// 		} else {
// 			return trace.Wrap(err)
// 		}
// 	}
// 	return nil
// }

// func printNodesAsText(nodes []types.Server, verbose bool) {
// 	// Reusable function to get addr or tunnel for each node
// 	getAddr := func(n types.Server) string {
// 		if n.GetUseTunnel() {
// 			return "⟵ Tunnel"
// 		}
// 		return n.GetAddr()
// 	}
// }

// 	var t asciitable.Table
// 	switch verbose {
// 	// In verbose mode, print everything on a single line and include the Node
// 	// ID (UUID). Useful for machines that need to parse the output of "tsh ls".
// 	case true:
// 		t = asciitable.MakeTable([]string{"Node Name", "Node ID", "Address", "Labels"})
// 		for _, n := range nodes {
// 			t.AddRow([]string{
// 				n.GetHostname(), n.GetName(), getAddr(n), n.LabelsString(),
// 			})
// 		}
// 	// In normal mode chunk the labels and print two per line and allow multiple
// 	// lines per node.
// 	case false:
// 		t = asciitable.MakeTable([]string{"Node Name", "Address", "Labels"})
// 		for _, n := range nodes {
// 			labelChunks := chunkLabels(n.GetAllLabels(), 2)
// 			for i, v := range labelChunks {
// 				if i == 0 {
// 					t.AddRow([]string{n.GetHostname(), getAddr(n), strings.Join(v, ", ")})
// 				} else {
// 					t.AddRow([]string{"", "", strings.Join(v, ", ")})
// 				}
// 			}
// 		}
// 	}

// 	fmt.Println(t.AsBuffer().String())
// }

// // makeClient takes the command-line configuration and constructs & returns
// // a fully configured TeleportClient object
// func makeClient(cf *CLIConf, useProfileLogin bool) (*client.TeleportClient, error) {
// 	// Parse OpenSSH style options.
// 	options, err := parseOptions(cf.Options)
// 	if err != nil {
// 		return nil, trace.Wrap(err)
// 	}

// 	// apply defaults
// 	if cf.MinsToLive == 0 {
// 		cf.MinsToLive = int32(defaults.CertDuration / time.Minute)
// 	}

// 	// split login & host
// 	hostLogin := cf.NodeLogin
// 	var labels map[string]string
// 	if cf.UserHost != "" {
// 		parts := strings.Split(cf.UserHost, "@")
// 		partsLength := len(parts)
// 		if partsLength > 1 {
// 			hostLogin = strings.Join(parts[:partsLength-1], "@")
// 			cf.UserHost = parts[partsLength-1]
// 		}
// 		// see if remote host is specified as a set of labels
// 		if strings.Contains(cf.UserHost, "=") {
// 			labels, err = client.ParseLabelSpec(cf.UserHost)
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 	}
// 	fPorts, err := client.ParsePortForwardSpec(cf.LocalForwardPorts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dPorts, err := client.ParseDynamicPortForwardSpec(cf.DynamicForwardedPorts)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 1: start with the defaults
// 	c := client.MakeDefaultConfig()

// 	// ProxyJump is an alias of Proxy flag
// 	if cf.ProxyJump != "" {
// 		hosts, err := utils.ParseProxyJump(cf.ProxyJump)
// 		if err != nil {
// 			return nil, trace.Wrap(err)
// 		}
// 		c.JumpHosts = hosts
// 	}

// 	// load profile. if no --proxy is given the currently active profile is used, otherwise
// 	// fetch profile for exact proxy we are trying to connect to.
// 	err = c.LoadProfile("", cf.Proxy)
// 	if err != nil {
// 		fmt.Printf("WARNING: Failed to load tsh profile for %q: %v\n", cf.Proxy, err)
// 	}

// 	// 3: override with the CLI flags
// 	if cf.Namespace != "" {
// 		c.Namespace = cf.Namespace
// 	}
// 	if cf.Username != "" {
// 		c.Username = cf.Username
// 	}
// 	// if proxy is set, and proxy is not equal to profile's
// 	// loaded addresses, override the values
// 	if cf.Proxy != "" && c.WebProxyAddr == "" {
// 		err = c.ParseProxyHost(cf.Proxy)
// 		if err != nil {
// 			return nil, trace.Wrap(err)
// 		}
// 	}
// 	if len(fPorts) > 0 {
// 		c.LocalForwardPorts = fPorts
// 	}
// 	if len(dPorts) > 0 {
// 		c.DynamicForwardedPorts = dPorts
// 	}
// 	if cf.SiteName != "" {
// 		c.SiteName = cf.SiteName
// 	}

// 	// Currently not supported by teleport client go
// 	// if cf.KubernetesCluster != "" {
// 	//	c.KubernetesCluster = cf.KubernetesCluster
// 	//}
// 	// if cf.DatabaseService != "" {
// 	//	c.DatabaseService = cf.DatabaseService
// 	//}

// 	// if host logins stored in profiles must be ignored...
// 	if !useProfileLogin {
// 		c.HostLogin = ""
// 	}
// 	if hostLogin != "" {
// 		c.HostLogin = hostLogin
// 	}
// 	c.Host = cf.UserHost
// 	c.HostPort = int(cf.NodePort)
// 	c.Labels = labels
// 	c.KeyTTL = time.Minute * time.Duration(cf.MinsToLive)
// 	c.InsecureSkipVerify = cf.InsecureSkipVerify

// 	// If a TTY was requested, make sure to allocate it. Note this applies to
// 	// "exec" command because a shell always has a TTY allocated.
// 	if cf.Interactive || options.RequestTTY {
// 		c.Interactive = true
// 	}

// 	if !cf.NoCache {
// 		c.CachePolicy = &client.CachePolicy{}
// 	}

// 	// check version compatibility of the server and client
// 	c.CheckVersions = !cf.SkipVersionCheck

// 	// parse compatibility parameter
// 	certificateFormat, err := parseCertificateCompatibilityFlag(cf.Compatibility, cf.CertificateFormat)
// 	if err != nil {
// 		return nil, trace.Wrap(err)
// 	}
// 	c.CertificateFormat = certificateFormat

// 	// copy the authentication connector over
// 	if cf.AuthConnector != "" {
// 		c.AuthConnector = cf.AuthConnector
// 	}

// 	// If agent forwarding was specified on the command line enable it.
// 	if cf.ForwardAgent || options.ForwardAgent {
// 		c.ForwardAgent = true
// 	}

// 	// If the caller does not want to check host keys, pass in a insecure host
// 	// key checker.
// 	if !options.StrictHostKeyChecking {
// 		c.HostKeyCallback = client.InsecureSkipHostKeyChecking
// 	}
// 	c.BindAddr = cf.BindAddr

// 	// Don't execute remote command, used when port forwarding.
// 	c.NoRemoteExec = cf.NoRemoteExec

// 	// Allow the default browser used to open tsh login links to be overridden
// 	// (not currently implemented) or set to 'none' to suppress browser opening entirely.
// 	c.Browser = cf.Browser

// 	// c.AddKeysToAgent = cf.AddKeysToAgent
// 	// if !cf.UseLocalSSHAgent {
// 	//	c.AddKeysToAgent = client.AddKeysToAgentNo
// 	// }

// 	c.EnableEscapeSequences = cf.EnableEscapeSequences

// 	// pass along mock sso login if provided (only used in tests)
// 	// c.MockSSOLogin = cf.mockSSOLogin

// 	tc, err := client.NewClient(c)
// 	if err != nil {
// 		return nil, trace.Wrap(err)
// 	}
// 	// If identity file was provided, we skip loading the local profile info
// 	// (above). This profile info provides the proxy-advertised listening
// 	// addresses.
// 	// To compensate, when using an identity file, explicitly fetch these
// 	// addresses from the proxy (this is what Ping does).
// 	if cf.IdentityFileIn != "" {
// 		log.Debug("Pinging the proxy to fetch listening addresses for non-web ports.")
// 		if _, err := tc.Ping(cf.Context); err != nil {
// 			return nil, trace.Wrap(err)
// 		}
// 	}
// 	return tc, nil
// }

// func parseCertificateCompatibilityFlag(compatibility string, certificateFormat string) (string, error) {
// 	switch {
// 	// if nothing is passed in, the role will decide
// 	case compatibility == "" && certificateFormat == "":
// 		return teleport.CertificateFormatUnspecified, nil
// 	// supporting the old --compat format for backward compatibility
// 	case compatibility != "" && certificateFormat == "":
// 		return utils.CheckCertificateFormatFlag(compatibility)
// 	// new documented flag --cert-format
// 	case compatibility == "" && certificateFormat != "":
// 		return utils.CheckCertificateFormatFlag(certificateFormat)
// 	// can not use both
// 	default:
// 		return "", trace.BadParameter("--compat or --cert-format must be specified")
// 	}
// }

// // chunkLabels breaks labels into sized chunks. Used to improve readability
// // of "tsh ls".
// func chunkLabels(labels map[string]string, chunkSize int) [][]string {
// 	// First sort labels so they always occur in the same order.
// 	sorted := make([]string, 0, len(labels))
// 	for k, v := range labels {
// 		sorted = append(sorted, fmt.Sprintf("%v=%v", k, v))
// 	}
// 	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

// 	// Then chunk labels into sized chunks.
// 	var chunks [][]string
// 	for chunkSize < len(sorted) {
// 		sorted, chunks = sorted[chunkSize:], append(chunks, sorted[0:chunkSize:chunkSize])
// 	}
// 	chunks = append(chunks, sorted)

// 	return chunks
// }

// func printProfiles(debug bool, profile *client.ProfileStatus, profiles []*client.ProfileStatus) {
// 	// Print the active profile.
// 	if profile != nil {
// 		printStatus(debug, profile, true)
// 	}

// 	// Print all other profiles.
// 	for _, p := range profiles {
// 		printStatus(debug, p, false)
// 	}
// }

// func printCurrentProfile(debug bool, profile *client.ProfileStatus) {
// 	if profile != nil {
// 		printStatus(debug, profile, true)
// 	}
// }

// // printStatus prints the status of the profile.
// func printStatus(debug bool, p *client.ProfileStatus, isActive bool) {
// 	var count int
// 	var prefix string
// 	if isActive {
// 		prefix = "> "
// 	} else {
// 		prefix = "  "
// 	}
// 	duration := time.Until(p.ValidUntil)
// 	humanDuration := "EXPIRED"
// 	if duration.Nanoseconds() > 0 {
// 		humanDuration = fmt.Sprintf("valid for %v", duration.Round(time.Minute))
// 	}

// 	fmt.Printf("%vProfile URL:        %v\n", prefix, p.ProxyURL.String())
// 	fmt.Printf("  Logged in as:       %v\n", p.Username)
// 	if p.Cluster != "" {
// 		fmt.Printf("  Cluster:            %v\n", p.Cluster)
// 	}
// 	fmt.Printf("  Roles:              %v*\n", strings.Join(p.Roles, ", "))
// 	if debug {
// 		for k, v := range p.Traits {
// 			if count == 0 {
// 				fmt.Printf("  Traits:             %v: %v\n", k, v)
// 			} else {
// 				fmt.Printf("                      %v: %v\n", k, v)
// 			}
// 			count++
// 		}
// 	}
// 	fmt.Printf("  Logins:             %v\n", strings.Join(p.Logins, ", "))
// 	// if p.KubeEnabled {
// 	//	fmt.Printf("  Kubernetes:         enabled\n")
// 	//	if p.KubeCluster != "" {
// 	//		fmt.Printf("  Kubernetes cluster: %q\n", p.KubeCluster)
// 	//	}
// 	//	if len(p.KubeUsers) > 0 {
// 	//		fmt.Printf("  Kubernetes users:   %v\n", strings.Join(p.KubeUsers, ", "))
// 	//	}
// 	//	if len(p.KubeGroups) > 0 {
// 	//		fmt.Printf("  Kubernetes groups:  %v\n", strings.Join(p.KubeGroups, ", "))
// 	//	}
// 	// } else {
// 	//	fmt.Printf("  Kubernetes:         disabled\n")
// 	// }
// 	// if len(p.Databases) != 0 {
// 	//	fmt.Printf("  Databases:          %v\n", strings.Join(p.DatabaseServices(), ", "))
// 	// }
// 	fmt.Printf("  Valid until:        %v [%v]\n", p.ValidUntil, humanDuration)
// 	fmt.Printf("  Extensions:         %v\n", strings.Join(p.Extensions, ", "))

// 	fmt.Printf("\n")
// }

// // host is a utility function that extracts
// // host from the host:port pair, in case of any error
// // returns the original value
// func host(in string) string {
// 	out, err := utils.Host(in)
// 	if err != nil {
// 		return in
// 	}
// 	return out
// }

// // setupNoninteractiveClient sets up existing client to use
// // non-interactive authentication methods
// func setupNoninteractiveClient(tc *client.TeleportClient, key *client.Key) error {
// 	certUsername, err := key.CertUsername()
// 	if err != nil {
// 		return trace.Wrap(err)
// 	}
// 	tc.Username = certUsername

// 	// Extract and set the HostLogin to be the first principal. It doesn't
// 	// matter what the value is, but some valid principal has to be set
// 	// otherwise the certificate won't be validated.
// 	certPrincipals, err := key.CertPrincipals()
// 	if err != nil {
// 		return trace.Wrap(err)
// 	}
// 	if len(certPrincipals) == 0 {
// 		return trace.BadParameter("no principals found")
// 	}
// 	tc.HostLogin = certPrincipals[0]

// 	identityAuth, err := authFromIdentity(key)
// 	if err != nil {
// 		return trace.Wrap(err)
// 	}
// 	tc.TLS, err = key.ClientTLSConfig()
// 	if err != nil {
// 		return trace.Wrap(err)
// 	}
// 	tc.AuthMethods = []ssh.AuthMethod{identityAuth}
// 	tc.Interactive = false
// 	tc.SkipLocalAuth = true

// 	// When user logs in for the first time without a CA in ~/.tsh/known_hosts,
// 	// and specifies the -out flag, we need to avoid writing anything to
// 	// ~/.tsh/ but still validate the proxy cert. Because the existing
// 	// client.Client methods have a side-effect of persisting the CA on disk,
// 	// we do all of this by hand.
// 	//
// 	// Wrap tc.HostKeyCallback with a another checker. This outer checker uses
// 	// key.TrustedCA to validate the remote host cert first, before falling
// 	// back to the original HostKeyCallback.
// 	oldHostKeyCallback := tc.HostKeyCallback
// 	tc.HostKeyCallback = func(hostname string, remote net.Addr, hostKey ssh.PublicKey) error {
// 		checker := ssh.CertChecker{
// 			// ssh.CertChecker will parse hostKey, extract public key of the
// 			// signer (CA) and call IsHostAuthority. IsHostAuthority in turn
// 			// has to match hostCAKey to any known trusted CA.
// 			IsHostAuthority: func(hostCAKey ssh.PublicKey, address string) bool {
// 				for _, ca := range key.TrustedCA {
// 					caKeys, err := ca.SSHCertPublicKeys()
// 					if err != nil {
// 						return false
// 					}
// 					for _, caKey := range caKeys {
// 						if sshutils.KeysEqual(caKey, hostCAKey) {
// 							return true
// 						}
// 					}
// 				}
// 				return false
// 			},
// 		}
// 		err := checker.CheckHostKey(hostname, remote, hostKey)
// 		if err != nil && oldHostKeyCallback != nil {
// 			errOld := oldHostKeyCallback(hostname, remote, hostKey)
// 			if errOld != nil {
// 				return trace.NewAggregate(err, errOld)
// 			}
// 		}
// 		return nil
// 	}
// 	return nil
// }

// // authFromIdentity returns a standard ssh.Authmethod for a given identity file
// func authFromIdentity(k *client.Key) (ssh.AuthMethod, error) {
// 	signer, err := sshutils.NewSigner(k.Priv, k.Cert)
// 	if err != nil {
// 		return nil, trace.Wrap(err)
// 	}
// 	return client.NewAuthMethodForCert(signer), nil
// }

// func executeAccessRequest(cf *CLIConf) {
// 	if cf.DesiredRoles == "" {
// 		utils.FatalError(trace.BadParameter("one or more roles must be specified"))
// 	}
// 	roles := strings.Split(cf.DesiredRoles, ",")
// 	tc, err := makeClient(cf, true)
// 	if err != nil {
// 		utils.FatalError(err)
// 	}
// 	if cf.Username == "" {
// 		cf.Username = tc.Username
// 	}
// 	req, err := services.NewAccessRequest(cf.Username, roles...)
// 	if err != nil {
// 		utils.FatalError(err)
// 	}
// 	fmt.Fprintf(os.Stderr, "Seeking request approval... (id: %s)\n", req.GetName())
// 	if err := getRequestApproval(cf, tc, req); err != nil {
// 		utils.FatalError(err)
// 	}
// 	fmt.Fprintf(os.Stderr, "Approval received, getting updated certificates...\n\n")
// 	if err := reissueWithRequests(cf, tc, req.GetName()); err != nil {
// 		utils.FatalError(err)
// 	}
// 	onStatus(cf)
// }

// // getRequestApproval registers an access request with the auth server and waits for it to be approved.
// func getRequestApproval(cf *CLIConf, tc *client.TeleportClient, req services.AccessRequest) error {
// 	// set up request watcher before submitting the request to the admin server
// 	// in order to avoid potential race.
// 	filter := services.AccessRequestFilter{
// 		User: tc.Username,
// 	}
// 	watcher, err := tc.NewWatcher(cf.Context, services.Watch{
// 		Name: "await-request-approval",
// 		Kinds: []services.WatchKind{
// 			{
// 				Kind:   services.KindAccessRequest,
// 				Filter: filter.IntoMap(),
// 			},
// 		},
// 	})
// 	if err != nil {
// 		return trace.Wrap(err)
// 	}
// 	defer watcher.Close()
// 	if err := tc.CreateAccessRequest(cf.Context, req); err != nil {
// 		utils.FatalError(err)
// 	}
// Loop:
// 	for {
// 		select {
// 		case event := <-watcher.Events():
// 			switch event.Type {
// 			case backend.OpInit:
// 				log.Infof("Access-request watcher initialized...")
// 				continue Loop
// 			case backend.OpPut:
// 				r, ok := event.Resource.(*services.AccessRequestV3)
// 				if !ok {
// 					return trace.BadParameter("unexpected resource type %T", event.Resource)
// 				}
// 				if r.GetName() != req.GetName() || r.GetState().IsPending() {
// 					log.Infof("Skipping put event id=%s,state=%s.", r.GetName(), r.GetState())
// 					continue Loop
// 				}
// 				if !r.GetState().IsApproved() {
// 					return trace.Errorf("request %s has been set to %s", r.GetName(), r.GetState().String())
// 				}
// 				return nil
// 			case backend.OpDelete:
// 				if event.Resource.GetName() != req.GetName() {
// 					log.Infof("Skipping delete event id=%s", event.Resource.GetName())
// 					continue Loop
// 				}
// 				return trace.Errorf("request %s has expired or been deleted...", event.Resource.GetName())
// 			default:
// 				log.Warnf("Skipping unknown event type %s", event.Type)
// 			}
// 		case <-watcher.Done():
// 			utils.FatalError(watcher.Error())
// 		}
// 	}
// }

// // reissueWithRequests handles a certificate reissue, applying new requests by ID,
// // and saving the updated profile.
// func reissueWiåthRequests(cf *CLIConf, tc *client.TeleportClient, reqIDs ...string) error {
// 	profile, _, err := client.Status("", cf.Proxy)
// 	if err != nil {
// 		return trace.Wrap(err)
// 	}
// 	params := client.ReissueParams{
// 		AccessRequests: reqIDs,
// 		RouteToCluster: cf.SiteName,
// 	}
// 	// if the certificate already had active requests, add them to our inputs parameters.
// 	if len(profile.ActiveRequests.AccessRequests) > 0 {
// 		params.AccessRequests = append(params.AccessRequests, profile.ActiveRequests.AccessRequests...)
// 	}
// 	if params.RouteToCluster == "" {
// 		params.RouteToCluster = profile.Cluster
// 	}
// 	if err := tc.ReissueUserCerts(cf.Context, params); err != nil {
// 		return trace.Wrap(err)
// 	}
// 	if err := tc.SaveProfile("", ""); err != nil {
// 		return trace.Wrap(err)
// 	}
// 	if err := kubeconfig.UpdateWithClient("", tc); err != nil {
// 		return trace.Wrap(err)
// 	}
// 	return nil
// }
