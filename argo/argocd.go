package argo

import (
	"fmt"
	"github.com/magefile/mage/mg"
)

type ArgoCD mg.Namespace

type ArgoCdConfig struct {
	Namespace       string
	Version         string
	PortForwardPort string
	SSHKeyPath      string
}

var (
	ArgoCDConfig = ArgoCdConfig{
		Namespace:       "argocd",
		Version:         "v2.11.3", // use `stable` for the latest version
		SSHKeyPath:      "~/.ssh/id_rsa",
		PortForwardPort: "8080",
	}
)

// Install Creates the argocd namespace and installs argocd
func (a ArgoCD) Install() error {
	// Create the ArgoCD namespace
	output, err := createNamespace(ArgoCDConfig.Namespace)
	if err != nil {
		return fmt.Errorf("unable to create argocd namespace. ERROR: %s", err)
	}
	fmt.Println(output)

	// Deploy Argo on the cluster
	output, err = run(fmt.Sprintf("kubectl apply -n %s -f https://raw.githubusercontent.com/argoproj/argo-cd/%s/manifests/install.yaml", ArgoCDConfig.Namespace, ArgoCDConfig.Version))
	if err != nil {
		return fmt.Errorf("unable to deploy argocd. ERROR: %s", err)
	}
	fmt.Println(output)

	return nil
}

// PortForward Port-forward the argocd gitops server
func (a ArgoCD) PortForward() error {
	fmt.Println(fmt.Sprintf("Argo can be accessed at:\nhttps://localhost:%s", ArgoCDConfig.PortForwardPort))
	// Port forward the argo-server
	_, err := run(fmt.Sprintf("kubectl port-forward svc/argocd-server -n %s %s:443", ArgoCDConfig.Namespace, ArgoCDConfig.PortForwardPort))
	if err != nil {
		return fmt.Errorf("unable to port-forward svc/argocd-server. ERROR: %s", err)
	}

	return nil
}

// GetAdminPassword Get the initial ArgoCD admin password
func (a ArgoCD) GetAdminPassword() error {
	// Fetching admin password
	output, err := run(fmt.Sprintf("argocd admin initial-password -n %s | head -n 1", ArgoCDConfig.Namespace))
	if err != nil {
		return fmt.Errorf("unable fetch admin credentials. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// Login Login to argo via the cli (requires the argocd service to be accessible)
func (a ArgoCD) Login() error {
	// Fetching admin password
	adminPass, err := run(fmt.Sprintf("argocd admin initial-password -n %s | head -n 1", ArgoCDConfig.Namespace))
	if err != nil {
		return fmt.Errorf("unable fetch admin credentials. ERROR: %s", err)
	}

	// Running argocd login using admin pass
	output, err := run(fmt.Sprintf("argocd login --username admin --password %s --insecure localhost:%s", adminPass, ArgoCDConfig.PortForwardPort))
	if err != nil {
		return fmt.Errorf("unable to login. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// AddHostSSHCert Add host ssh cert - Expected args: (hostname)
func (a ArgoCD) AddHostSSHCert(hostname string) error {
	// Default port
	defaultPort := "22"

	return a.AddHostSSHCertWithPort(hostname, defaultPort)
}

// AddHostSSHCertWithPort Add host ssh cert - Expected args: (hostname)
func (a ArgoCD) AddHostSSHCertWithPort(hostname, port string) error {
	mg.Deps(a.Login)

	// Add github ssh cert
	output, err := run(fmt.Sprintf("ssh-keyscan -p %s %s | argocd cert add-ssh --batch", port, hostname))
	if err != nil {
		return fmt.Errorf("unable add %s ssh cert. ERROR: %s", hostname, err)
	}

	fmt.Println(output)

	return nil
}

// AddGithubSSHCert Add github ssh cert
func (a ArgoCD) AddGithubSSHCert() error {
	mg.Deps(a.Login)

	// Add github ssh cert
	return a.AddHostSSHCert("github.com")
}

// AddRepoSSHCreds Add Argo repo credentials - Expected args: (repoURL sshKeyPath)
func (a ArgoCD) AddRepoSSHCreds(repoURL, sshKeyPath string) error {
	mg.Deps(a.Login)

	// Add repocreds
	output, err := run(fmt.Sprintf("argocd repocreds add %s --ssh-private-key-path %s", repoURL, sshKeyPath))
	if err != nil {
		return fmt.Errorf("unable add repocreds. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// AddGithubSSHCreds Add Argo repo credentials
func (a ArgoCD) AddGithubSSHCreds() error {
	githubSSHUrl := "git@github.com"

	// Add repocreds
	return a.AddRepoSSHCreds(githubSSHUrl, ArgoCDConfig.SSHKeyPath)
}

// AddHTTPRepo Add HTTP repository to Argo - Expected args: (repoURL argoURL)
func (a ArgoCD) AddHTTPRepo(repoURL, argoURL string) error {
	// Add new repo
	output, err := run(fmt.Sprintf("argocd repo add %s --server %s", repoURL, argoURL))
	if err != nil {
		return fmt.Errorf("unable add repository. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// AddRepoSSH Add SSH repository to Argo - Expected args: (repoURL sshKeyPath argoURL)
func (a ArgoCD) AddRepoSSH(repoURL, sshKeyPath, argoURL string) error {
	// Add new repo
	output, err := run(fmt.Sprintf("argocd repo add %s --ssh-private-key-path %s --server %s", repoURL, sshKeyPath, argoURL))
	if err != nil {
		return fmt.Errorf("unable add repository. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// CreateAppCLI Created new application via argocd cli - Expected args: (appName path repoURL  namespace)
func (a ArgoCD) CreateAppCLI(appName, path, repoURL, namespace string) error {
	// Add new app via argocd cli
	output, err := run(fmt.Sprintf("argocd app create %s --repo %s --path %s --dest-server https://kubernetes.default.svc --dest-namespace %s", appName, repoURL, path, namespace))

	if err != nil {
		return fmt.Errorf("unable add application. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}

// CreateAppManifest Created new application via manifest  - Expected args: (manifestPath)
func (a ArgoCD) CreateAppManifest(manifestPath string) error {
	// Add new app via manifest
	output, err := run(fmt.Sprintf("kubectl apply -f %s", manifestPath))
	if err != nil {
		return fmt.Errorf("unable add application. ERROR: %s", err)
	}

	fmt.Println(output)

	return nil
}
