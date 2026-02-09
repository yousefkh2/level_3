package config

import (
	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewKubernetesClient() (client.Client, error) {
	// try in-cluster config first
	config, err := rest.InClusterConfig()
	if err != nil {
		// Not in cluster, try kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", "../../week_3/infrastructure/kubeconfig.yaml")
		if err != nil {
			return nil, err
		}
	}

	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	return nil, err
	// }

	/*create a scheme (a scheme is like a registry of types) (tells the client what types exist)
	tells the client "here's how to convert between Go structs and JSON/YAML for the API."
	- the main function of the scheme is as a reference how to translate from that struct to json
	Without it, the client wouldn't know how to translate cnpgv1.Cluster into the JSON that K8s expects
	*/
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme) // teaches it about built-in K8s types
	_ = cnpgv1.AddToScheme(scheme)         // teaches it about CNPG types

	//create the controller-runtime client
	//one client to rule them all! it knows both std k8s resources AND cnpg resources
	k8sClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return k8sClient, nil
}
