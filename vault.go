package gok8s

import (
	"errors"
	"fmt"

	"github.com/pablo-ruth/govault"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// New k8s client with config from Vault
func NewClientFromVault(addr, path, roleID, secretID string) (*kubernetes.Clientset, *rest.Config, error) {

	// Init Vault client
	client := govault.Client{
		Address:       addr,
		TLSSkipVerify: true,
	}

	// Try to login to Vault with AppRole
	err := client.AppRoleLogin(roleID, secretID)
	if err != nil {
		return nil, nil, fmt.Errorf("Approle login failed: %s", err)
	}
	defer client.Logout()

	// Read k8s config in specified path
	entry, err := client.Read(path, 200)
	if err != nil {
		return nil, nil, err
	}

	// Get CA to connect to k8s
	certificateAuthorityData, ok := entry["certificate-authority-data"]
	if !ok {
		return nil, nil, errors.New("Missing certificate-authority-data")
	}

	// Get Token to connect to k8s
	token, ok := entry["token"]
	if !ok {
		return nil, nil, errors.New("Missing token")
	}

	// Get server url to connect to k8s
	server, ok := entry["server"]
	if !ok {
		return nil, nil, errors.New("Missing server")
	}

	// Init a new k8s client with previously imported config
	k8sClient, config, err := NewClient(server.(string), token.(string), []byte(certificateAuthorityData.(string)), false)
	if err != nil {
		return nil, nil, err
	}

	return k8sClient, config, nil
}
