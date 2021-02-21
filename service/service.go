package service

import (
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
	"encoding/json"
	"sigs.k8s.io/yaml"
)

type service struct {
	client 					*kubernetes.Clientset
	err 					error
}

type Service interface {
	SetErr(err error)
	Create(input string) (err error)
	Update(input string) (err error)
	Delete(namespace string, serviceName string) (err error)
	Get(namespace string, name string) (ret string, err error)
}

func NewForClient(client *kubernetes.Clientset) *service {
	return &service{
		client: client,
	}
}

func (c *service) SetErr(err error) {
	c.err = err
}

func (c *service) Get(namespace string, serviceName string) (serviceStr string, err error) {
	if c.err != nil {
		return "", c.err
	}
	service, err := c.client.CoreV1().Services(namespace).Get(serviceName, metav1.GetOptions{})
	if err != nil {
		return
	}
	serviceJsonBytes, err := json.Marshal(service)
	if err != nil {
		return
	}
	serviceStr = string(serviceJsonBytes)
	return
}

func (c *service) Create(input string) (err error) {
	if c.err != nil {
		return c.err
	}
	service := new(v1.Service)
	err = yaml.Unmarshal([]byte(input), service)
	if err != nil {
		return
	}
	// output for debug
	_, err = yaml.Marshal(service)
	if err != nil {
		return
	}
	namespace := service.Namespace
	_, err = c.client.CoreV1().Services(namespace).Create(service)
	if err != nil {
		return
	}
	return
}

func (c *service) Update(input string) (err error) {
	if c.err != nil {
		return c.err
	}
	service := new(v1.Service)
	err = yaml.Unmarshal([]byte(input), service)
	if err != nil {
		return
	}
	// output for debug
	_, err = yaml.Marshal(service)
	if err != nil {
		return
	}
	namespace := service.Namespace
	_, err = c.client.CoreV1().Services(namespace).Update(service)
	if err != nil {
		return
	}
	return
}

func (c *service) Delete(namespace string, serviceName string) (err error) {
	if c.err != nil {
		return c.err
	}
	err = c.client.CoreV1().Services(namespace).Delete(serviceName, &metav1.DeleteOptions{})
	if err != nil {
		return
	}
	return
}