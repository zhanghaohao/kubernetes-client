package secret

import (
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
	"encoding/json"
)

type secret struct {
	client 					*kubernetes.Clientset
	err 					error
}

type Secret interface {
	SetErr(err error)
	Create(input string) (err error)
	Delete(namespace string, name string) (err error)
	Update(input string) (err error)
	Get(namespace string, name string) (ret string, err error)
}

func NewForClient(client *kubernetes.Clientset) *secret {
	return &secret{
		client: client,
	}
}

func (c *secret) SetErr(err error)  {
	c.err = err
}

func (c *secret) Create(input string) (err error) {
	if c.err != nil {
		return c.err
	}
	secret := new(corev1.Secret)
	err = yaml.Unmarshal([]byte(input), secret)
	if err != nil {
		return
	}
	// output for debug
	_, err = yaml.Marshal(secret)
	if err != nil {
		return
	}
	namespace := secret.Namespace
	_, err = c.client.CoreV1().Secrets(namespace).Create(secret)
	if err != nil {
		return
	}
	return
}

func (c *secret) Delete(namespace string, name string) (err error) {
	if c.err != nil {
		return c.err
	}
	err = c.client.CoreV1().Secrets(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return
	}
	return
}

func (c *secret) Update(input string) (err error) {
	if c.err != nil {
		return c.err
	}
	secret := new(corev1.Secret)
	err = yaml.Unmarshal([]byte(input), secret)
	if err != nil {
		return
	}
	// output for debug
	_, err = yaml.Marshal(secret)
	if err != nil {
		return
	}
	namespace := secret.Namespace
	_, err = c.client.CoreV1().Secrets(namespace).Update(secret)
	if err != nil {
		return
	}
	return
}

func (c *secret) Get(namespace string, name string) (ret string, err error) {
	if c.err != nil {
		return "", c.err
	}
	secret, err := c.client.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	d, err := json.Marshal(secret)
	if err != nil {
		return "", err
	}
	ret = string(d)
	return
}