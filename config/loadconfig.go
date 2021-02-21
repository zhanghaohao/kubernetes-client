package config

import (
	"os"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	"fmt"
)

func readKubeConfigFile(kubeConfigPath string) (kubeConfig *rest.Config, err error) {
	kubeConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		err = fmt.Errorf("Read kubeconfig Error:\n %s", err)
		return
	}
	return
}

func BuildKubernetesClientFromKubeConfig(kubeConfig *rest.Config) (client *kubernetes.Clientset, err error) {
	client, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		err = fmt.Errorf("build kubernetes client error:\n %s", err)
		return
	}
	return
}

func BuildKubernetesClientFromKubeConfigFile(kubeConfigPath string) (client *kubernetes.Clientset, err error) {
	kubeConfig, err := readKubeConfigFile(kubeConfigPath)
	if err != nil {
		return
	}
	client, err = BuildKubernetesClientFromKubeConfig(kubeConfig)
	if err != nil {
		return
	}
	return
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

