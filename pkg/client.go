package pkg

import (
	"fmt"

	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

func NewKubeClient(context string) (clientset.Interface, error) {
	config, err := GetConfig(context).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get kubernetes config for context '%s': %s", context, err)
	}
	return clientset.NewForConfig(config)
}

func GetConfig(context string) clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}

	if context != "" {
		overrides.CurrentContext = context
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}
