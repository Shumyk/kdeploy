package cmd

import (
	"os"
	"path/filepath"

	util "shumyk/kdeploy/cmd/util"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	clientSet *kubernetes.Clientset

	k8sNamespace        string
	k8sResourceFullName string
	k8sResourceType     string
)

func CreateClientConfigFromKubeConfig() clientcmd.ClientConfig {
	k8sConfigPath := filepath.Join(clientcmd.RecommendedConfigDir, clientcmd.RecommendedFileName)
	k8sConfigBytes, err := os.ReadFile(k8sConfigPath)
	util.ErrorCheck(err, "Failed to read Kube config")

	conf, err := clientcmd.NewClientConfigFromBytes(k8sConfigBytes)
	util.ErrorCheck(err, "Failed to create new k8s API client config")
	return conf
}

func LoadMetadata(config clientcmd.ClientConfig) {
	var err error
	k8sNamespace, _, err = config.Namespace()
	util.ErrorCheck(err, "Failed to resolve k8s namespace")

	k8sResourceFullName = ResolveResourceName()
	k8sResourceType = ResolveResourceType()

	util.PrintEnvironmentInfo(arg_k8sResourceFullName, k8sNamespace)
}

func ClientSet(config clientcmd.ClientConfig, ch chan<- bool) {
	configGetter := kubeConfigGetter(config)
	k8sRestConfig, err := clientcmd.BuildConfigFromKubeconfigGetter("", configGetter)
	util.ErrorCheck(err, "Building config from kube config getter failed")

	clientSet, err = kubernetes.NewForConfig(k8sRestConfig)
	util.ErrorCheck(err, "Creating Client Set failed")

	ch <- true
}

func kubeConfigGetter(c clientcmd.ClientConfig) clientcmd.KubeconfigGetter {
	return func() (*api.Config, error) {
		c, err := c.RawConfig()
		return &c, err
	}
}
