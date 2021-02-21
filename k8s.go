package k8s

import (
	"github.com/zhanghaohao/kubernetes-client/service"
	k8sconfig "github.com/zhanghaohao/kubernetes-client/config"
	"github.com/zhanghaohao/kubernetes-client/pod"
	"github.com/zhanghaohao/kubernetes-client/namespace"
	"github.com/zhanghaohao/kubernetes-client/event"
	"github.com/zhanghaohao/kubernetes-client/app"
	"github.com/zhanghaohao/kubernetes-client/configmap"
	"github.com/zhanghaohao/kubernetes-client/secret"
	"k8s.io/client-go/kubernetes"
	"github.com/zhanghaohao/kubernetes-client/batch"
	"fmt"
	"k8s.io/client-go/rest"
)

const (
	KubernetesDeployment ResourceObjectType = "deployment"
	KubernetesService ResourceObjectType = "service"
	KubernetesJob 	ResourceObjectType = "job"
	KubernetesConfigMap ResourceObjectType = "configMap"
	KubernetesEvent ResourceObjectType = "event"
	KubernetesPod 	ResourceObjectType = "pod"
	KubernetesSecret ResourceObjectType = "secret"

	defaultClusterName = "default"
	defaultKubeConfigPath = "/root/kubeconfig"
)

type ResourceObjectType string

type ResourceObjectRegister map[ResourceObjectType]ResourceObject

type KubeConfigPath struct {
	ClusterName 			string
	Path 					string
}

type RestKubeConfig struct {
	ClusterName 			string
	KubeConfig				*rest.Config
}

type k8sClient struct {
	client 					*kubernetes.Clientset
	err 					error
}

type k8sClients struct {
	clients 				map[string]*kubernetes.Clientset
}

type K8SClients interface {
	GetClient(clusterName string) K8SClient
	Load(clusterName string, client *kubernetes.Clientset)
}

type K8SClient interface {
	CommonResourceObject(resourceObjectType ResourceObjectType) ResourceObject
	Service() service.Service
	Pod() pod.Pod
	Namespace() namespace.Namespace
	Event() event.Event
	Deployment() app.Deployment
	ConfigMap() configmap.ConfigMap
	Secret() secret.Secret
	Job() batch.Job
}

type ResourceObject interface {
	SetErr(err error)
	Create(input string) (err error)
	Delete(namespace string, name string) (err error)
	Update(input string) (err error)
	Get(namespace string, name string) (ret string, err error)
}

func (c ResourceObjectType) String() string {
	return string(c)
}

func NewFromDefaultKubeconfigPaths() (clients K8SClients, err error) {
	paths := []KubeConfigPath{
		{
			ClusterName: defaultClusterName,
			Path: defaultKubeConfigPath,
		},
	}
	clients, err = NewFromKubeconfigPaths(paths)
	if err != nil {
		return
	}
	return
}

func NewFromKubeconfigPaths(paths []KubeConfigPath) (clients K8SClients, err error) {
	/*
	build k8s clients from kubeConfig file paths
	 */
	if len(paths) == 0 {
		err := fmt.Errorf("empty kubeconfig paths provided")
		return nil, err
	}
	for _, path := range paths {
		client, err := k8sconfig.BuildKubernetesClientFromKubeConfigFile(path.Path)
		if err != nil {
			return nil, err
		}
		clients.Load(path.ClusterName, client)
	}
	return
}

func NewFromKubeconfigs(configs []RestKubeConfig) (clients K8SClients, err error) {
	/*
	build k8s clients from rest kubeConfigs
	 */
	for _, config := range configs {
		client, err := k8sconfig.BuildKubernetesClientFromKubeConfig(config.KubeConfig)
		if err != nil {
			return nil, err
		}
		clients.Load(config.ClusterName, client)
	}
	return
}

func (k *k8sClients) GetClient(clusterName string) K8SClient {
	r := new(k8sClient)
	client, ok := k.clients[clusterName]
	if !ok {
		r.err = fmt.Errorf("invalid clusterName %s", clusterName)
		return r
	}
	r.client = client
	return r
}

func (k *k8sClients) Load(clusterName string, client *kubernetes.Clientset) {
	k.clients[clusterName] = client
}

func (k *k8sClient) register() (r ResourceObjectRegister) {
	r[KubernetesDeployment] = app.NewDeploymentForClient(k.client)
	r[KubernetesService] = service.NewForClient(k.client)
	r[KubernetesJob] = batch.NewForClient(k.client)
	r[KubernetesConfigMap] = configmap.NewForClient(k.client)
	r[KubernetesEvent] = event.NewForClient(k.client)
	r[KubernetesPod] = pod.NewForClient(k.client)
	r[KubernetesSecret] = secret.NewForClient(k.client)
	return
}

func (k *k8sClient) CommonResourceObject(resourceObjectType ResourceObjectType) ResourceObject {
	r := k.register()
	o, ok := r[resourceObjectType]
	if !ok {
		err := fmt.Errorf("invalid resourceObjectType %s", resourceObjectType)
		o.SetErr(err)
	}
	return o
}

func (k *k8sClient) Service() service.Service {
	r := service.NewForClient(k.client)
	if k.err != nil {
		r.SetErr(k.err)
	}
	return r
}

func (k *k8sClient) Pod() pod.Pod {
	r := pod.NewForClient(k.client)
	if k.err != nil {
		r.SetErr(k.err)
	}
	return r
}

func (k *k8sClient) Namespace() namespace.Namespace {
	r := namespace.NewForClient(k.client)
	if k.err != nil {
		r.SetErr(k.err)
	}
	return r
}

func (k *k8sClient) Event() event.Event {
	r := event.NewForClient(k.client)
	if k.err != nil {
		r.SetErr(k.err)
	}
	return r
}

func (k *k8sClient) Deployment() app.Deployment {
	r := app.NewDeploymentForClient(k.client)
	if k.err != nil {
		r.SetErr(k.err)
	}
	return r
}

func (k *k8sClient) ConfigMap() configmap.ConfigMap {
	r := configmap.NewForClient(k.client)
	if k.err != nil {
		r.SetErr(k.err)
	}
	return r
}

func (k *k8sClient) Secret() secret.Secret {
	r := secret.NewForClient(k.client)
	if k.err != nil {
		r.SetErr(k.err)
	}
	return r
}

func (k *k8sClient) Job() batch.Job {
	r := batch.NewForClient(k.client)
	if k.err != nil {
		r.SetErr(k.err)
	}
	return r
}
