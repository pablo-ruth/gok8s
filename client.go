package gok8s

import (
	"errors"
	"io/ioutil"

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

// NewClientFromPod initialize a new Kubernetes client with service account defined in current pod
func NewClientFromPod() (*kubernetes.Clientset, *rest.Config, error) {
	token, err := ioutil.ReadFile("/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return nil, nil, errors.New("failed to read token file")
	}

	ca, err := ioutil.ReadFile("/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	if err != nil {
		return nil, nil, errors.New("failed to read ca file")
	}

	config := &rest.Config{
		Host:        "https://kubernetes.default.svc",
		BearerToken: string(token),
		TLSClientConfig: rest.TLSClientConfig{
			CAData:   ca,
			Insecure: false,
		},
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}
	return client, config, nil
}
