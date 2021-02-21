package batch

import (
	"k8s.io/client-go/kubernetes"
	batchv1 "k8s.io/api/batch/v1"
	"sigs.k8s.io/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"encoding/json"
)

type job struct {
	client 					*kubernetes.Clientset
	err 					error
}

type Job interface {
	SetErr(err error)
	Create(input string) (err error)
	Delete(namespace string, name string) (err error)
	Update(input string) (err error)
	Get(namespace string, name string) (ret string, err error)
	GetStatus(namespace string, jobName string) (status *batchv1.JobStatus, err error)
}

func NewForClient(client *kubernetes.Clientset) *job {
	return &job{
		client: client,
	}
}

func (c *job) SetErr(err error)  {
	c.err = err
}

func (c *job) Create(input string) (err error) {
	if c.err != nil {
		return c.err
	}
	job := new(batchv1.Job)
	err = yaml.Unmarshal([]byte(input), job)
	if err != nil {
		return
	}
	// output for debug
	_, err = yaml.Marshal(job)
	if err != nil {
		return
	}
	namespace := job.Namespace
	_, err = c.client.BatchV1().Jobs(namespace).Create(job)
	if err != nil {
		return
	}
	return
}

func (c *job) Delete(namespace string, name string) (err error) {
	if c.err != nil {
		return c.err
	}
	return c.client.BatchV1().Jobs(namespace).Delete(name, &metav1.DeleteOptions{})
}

func (c *job) Update(input string) (err error) {
	if c.err != nil {
		return c.err
	}
	job := new(batchv1.Job)
	err = yaml.Unmarshal([]byte(input), job)
	if err != nil {
		return
	}
	// output for debug
	_, err = yaml.Marshal(job)
	if err != nil {
		return
	}
	namespace := job.Namespace
	_, err = c.client.BatchV1().Jobs(namespace).Update(job)
	if err != nil {
		return
	}
	return
}

func (c *job) Get(namespace string, name string) (ret string, err error) {
	if c.err != nil {
		return "", c.err
	}
	job, err := c.client.BatchV1().Jobs(namespace).Get(name, metav1.GetOptions{})
	d, err := json.Marshal(job)
	if err != nil {
		return "", err
	}
	ret = string(d)
	return
}

func (c *job) GetStatus(namespace string, jobName string) (status *batchv1.JobStatus, err error) {
	if c.err != nil {
		return nil, c.err
	}
	job, err := c.client.BatchV1().Jobs(namespace).Get(jobName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &job.Status, nil
}