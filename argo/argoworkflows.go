package argo

import (
	"fmt"
	"github.com/magefile/mage/mg"
)

type ArgoWorkflows mg.Namespace

type ArgoWorkflowsConfig struct {
	Namespace       string
	Version         string
	PortForwardPort string
}

var (
	ArgoWFConfig = ArgoWorkflowsConfig{
		Namespace:       "argo",
		Version:         "v3.5.10", // use `stable` for the latest version
		PortForwardPort: "2746",
	}
)

// Install Install argocd workflows
func (ArgoWorkflows) Install() error {
	// Create the ArgoCD namespace
	output, err := createNamespace(ArgoWFConfig.Namespace)
	if err != nil {
		return fmt.Errorf("unable to create argocd namespace. ERROR: %s", err)
	}
	fmt.Println(output)

	// Deploy Argo on the cluster
	output, err = run(fmt.Sprintf("kubectl apply -n %s -f https://github.com/argoproj/argo-workflows/releases/download/%s/install.yaml", ArgoWFConfig.Namespace, ArgoWFConfig.Version))
	if err != nil {
		return fmt.Errorf("unable to deploy argocd. ERROR: %s", err)
	}
	fmt.Println(output)

	return nil
}

// ArgoServer Start the Argo Server
func (ArgoWorkflows) ArgoServer() error {
	// Port forward the argo-server
	_, err := run(fmt.Sprintf("argo server --namespace %s --auth-mode=server", ArgoWFConfig.Namespace))
	if err != nil {
		return fmt.Errorf("unable to start the argo-server. ERROR: %s", err)
	}

	return nil
}
