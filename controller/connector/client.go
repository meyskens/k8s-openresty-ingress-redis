package connector

import (
	"sync"

	core_v1 "k8s.io/api/core/v1"

	"k8s.io/api/networking/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset *kubernetes.Clientset

	// IngressChangeChan is called when an ingress change is detected
	IngressChangeChan chan struct{}
	ingresses         map[string]*v1beta1.Ingress
	ingressesMutex    sync.RWMutex

	// ServicesChangeChan is called when a service change is detected
	ServicesChangeChan chan struct{}
	services           map[string]*core_v1.Service
	servicesMutex      sync.RWMutex

	// SecretChangeChan is called when a secret change is detected
	SecretChangeChan chan struct{}
	secrets          map[string]*core_v1.Secret
	secretsMutex     sync.RWMutex
}

// NewClient generates a client with the right configuration
func NewClient() (*Client, error) {
	clientset, err := getInClusterClientset()
	if err != nil {
		//return nil, err
		clientset, err = getLocalClientSet("kind-kind") // TODO: change me
		if err != nil {
			return nil, err
		}
	}
	client := Client{
		clientset: clientset,
		ingresses: map[string]*v1beta1.Ingress{},
		services:  map[string]*core_v1.Service{},
		secrets:   map[string]*core_v1.Secret{},
	}
	return &client, nil
}

func getInClusterClientset() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func getLocalClientSet(context string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		},
	).ClientConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
