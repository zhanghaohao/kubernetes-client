package configmap

import (
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
	"encoding/json"
)

type configMap struct {
	client 					*kubernetes.Clientset
	err 					error
}

type ConfigMap interface {
	SetErr(err error)
	Create(input string) (err error)
	Delete(namespace string, name string) (err error)
	Update(input string) (err error)
	Get(namespace string, name string) (ret string, err error)
}

func NewForClient(client *kubernetes.Clientset) *configMap {
	return &configMap{
		client: client,
	}
}

func (c *configMap) SetErr(err error)  {
	c.err = err
}

func (c *configMap) Create(input string) (err error) {
	if c.err != nil {
		return c.err
	}
	configMap := new(corev1.ConfigMap)
	err = yaml.Unmarshal([]byte(input), configMap)
	if err != nil {
		return
	}
	// output for debug
	_, err = yaml.Marshal(configMap)
	if err != nil {
		return
	}
	namespace := configMap.Namespace
	_, err = c.client.CoreV1().ConfigMaps(namespace).Create(configMap)
	if err != nil {
		return
	}
	return
}

func (c *configMap) Delete(namespace string, name string) (err error) {
	if c.err != nil {
		return c.err
	}
	err = c.client.CoreV1().ConfigMaps(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return
	}
	return
}

func (c *configMap) Update(input string) (err error) {
	if c.err != nil {
		return c.err
	}
	configMap := new(corev1.ConfigMap)
	err = yaml.Unmarshal([]byte(input), configMap)
	if err != nil {
		return
	}
	// output for debug
	_, err = yaml.Marshal(configMap)
	if err != nil {
		return
	}
	namespace := configMap.Namespace
	_, err = c.client.CoreV1().ConfigMaps(namespace).Update(configMap)
	if err != nil {
		return
	}
	return
}

func (c *configMap) Get(namespace string, name string) (ret string, err error) {
	if c.err != nil {
		return "", c.err
	}
	configMap, err := c.client.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	d, err := json.Marshal(configMap)
	if err != nil {
		return "", err
	}
	ret = string(d)
	return
}