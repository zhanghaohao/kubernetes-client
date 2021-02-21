package namespace

import (
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type namespace struct {
	client 					*kubernetes.Clientset
	err 					error
}

type Namespace interface {
	SetErr(err error)
	/*
	create, delete and update in namespace is different from other resource objects
	  */
	Create(namespace string) (err error)
	Delete(namespace string) (err error)
	GetStatus(namespaceName string) (status string, err error)
}

func NewForClient(client *kubernetes.Clientset) *namespace {
	return &namespace{
		client: client,
	}
}

func (c *namespace) SetErr(err error)  {
	c.err = err
}

func (c *namespace) Create(namespaceName string) (err error) {
	if c.err != nil {
		return c.err
	}
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespaceName,
		},
	}
	_, err = c.client.CoreV1().Namespaces().Create(namespace)
	if err != nil {
		return
	}
	return
}

func (c *namespace) Delete(namespaceName string) (err error) {
	if c.err != nil {
		return c.err
	}
	err = c.client.CoreV1().Namespaces().Delete(namespaceName, &metav1.DeleteOptions{})
	if err != nil {
		return
	}
	return
}

func (c *namespace) GetStatus(namespaceName string) (status string, err error) {
	if c.err != nil {
		return "", c.err
	}
	namespace, err := c.client.CoreV1().Namespaces().Get(namespaceName, metav1.GetOptions{})
	if err != nil {
		return
	}
	status = string(namespace.Status.Phase)
	return
}