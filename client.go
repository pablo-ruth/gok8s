package gok8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewClient initialize a new Kubernetes client
func NewClient(host, token string, ca []byte, skipTLSVerify bool) (*kubernetes.Clientset, *rest.Config, error) {
	config := &rest.Config{
		Host:        host,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData:   ca,
			Insecure: skipTLSVerify,
		},
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return client, config, nil
}
