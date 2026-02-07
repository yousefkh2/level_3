package config

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubernetesClient() (*kubernetes.Clientset, error) {
	// try in-cluster config first
	config, err := rest.InClusterConfig()
	if err != nil {
		// Not in cluster, try kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", "../../week_3/infrastructure/kubeconfig.yaml")
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
