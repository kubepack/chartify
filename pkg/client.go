package pkg

import (
	"fmt"

	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

func NewKubeClient() (clientset.Interface, error) {
	config, err := GetConfig().ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("Could not get kubernetes config: %s", err)
	}
	return clientset.NewForConfig(config)
}

func GetConfig() clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}
