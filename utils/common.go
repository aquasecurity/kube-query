package utils

import (
	"github.com/kolide/osquery-go"
	"k8s.io/client-go/kubernetes"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"k8s.io/client-go/tools/clientcmd"
)

// CreateKubeClient Creates the kubernetes client using the kubeconfig path
func CreateKubeClient(kubeconfig string) (kubernetes.Interface, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	// create the clientset
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// CreateMetricsClient creates a metrics client to use for fetching metrics data of cluster
func CreateMetricsClient(kubeconfig string) (*metrics.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil{
        return nil, err
    }

    mc, err := metrics.NewForConfig(config)
    if err != nil {
        return nil, err
	}
	return mc, nil
}

// CreateOsQueryExtension Generates and registers an osquery extension
// using given osquery socket path
func CreateOsQueryExtension(name, socket string) (server *osquery.ExtensionManagerServer, err error) {
	server, err = osquery.NewExtensionManagerServer(name, socket)
	return
}
