package cmd

import (
	"os"
	"path/filepath"

	. "shumyk/kdeploy/cmd/util"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	clientSet *kubernetes.Clientset

	namespace       string
	microservice    string
	k8sResource     string
	k8sResourceName string
)

func CreateClientConfigFromKubeConfig() clientcmd.ClientConfig {
	k8sConfigPath := filepath.Join(clientcmd.RecommendedConfigDir, clientcmd.RecommendedFileName)
	k8sConfigBytes, err := os.ReadFile(k8sConfigPath)
	ErrorCheck(err, "Couldn't read kube config")

	conf, err := clientcmd.NewClientConfigFromBytes(k8sConfigBytes)
	ErrorCheck(err, "Failed creating new API server client")
	return conf
}

func LoadMetadata(config clientcmd.ClientConfig) {
	var err error
	namespace, _, err = config.Namespace()
	ErrorCheck(err, "Resolving namespace failed")

	k8sResourceName = namespace + "-" + microservice
	k8sResource = ResolveResourceType(microservice)

	PrintEnvironmentInfo(microservice, namespace)
}

func ClientSet(config clientcmd.ClientConfig, ch chan<- bool) {
	configGetter := kubeConfigGetter(config)
	k8sRestConfig, err := clientcmd.BuildConfigFromKubeconfigGetter("", configGetter)
	ErrorCheck(err, "Building config from kube config getter failed")

	clientSet, err = kubernetes.NewForConfig(k8sRestConfig)
	ErrorCheck(err, "Creating Client Set failed")

	ch <- true
}

func kubeConfigGetter(c clientcmd.ClientConfig) clientcmd.KubeconfigGetter {
	return func() (*api.Config, error) {
		c, err := c.RawConfig()
		return &c, err
	}
}
