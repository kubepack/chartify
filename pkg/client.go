
package pkg

import (
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
	client"k8s.io/kubernetes/pkg/client/unversioned"
)

func NewKubeClient(context string) (*client.Client, error) {
	config := GetConfig(context)
	factory := cmdutil.NewFactory(config)
	kubeClient, err := factory.Client()
	return kubeClient, err

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

