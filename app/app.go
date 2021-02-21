package app

import (
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/apps/v1"
	"sigs.k8s.io/yaml"
	"encoding/json"
)

type deployment struct {
	client 					*kubernetes.Clientset
	err 					error
}

type DeploymentStatus struct {
	Replicas 						int
	UpdatedReplicas 				int
	ReadyReplicas 					int
	AvailableReplicas 				int
	UnavailableReplicas				int
	ConditionType 					string
	ConditionStatus 				string
	LastUpdateTime					string
	Reason 							string
	Message 						string
}

type Deployment interface {
	SetErr(err error)
	Create(input string) (err error)
	Delete(namespace string, name string) (err error)
	Update(input string) (err error)
	Get(namespace string, name string) (ret string, err error)
	GetStatus(namespace string, deploymentName string) (deploymentStatus *DeploymentStatus, err error)
}

func NewDeploymentForClient(client *kubernetes.Clientset) *deployment {
	return &deployment{
		client: client,
	}
}

func (c *deployment) SetErr(err error)  {
	c.err = err
}

// trigger deployment include job change and rollback
func (c *deployment) Trigger(namespace string, deploymentName string, imageName string, imageTag string) (err error) {
	if c.err != nil {
		return c.err
	}
	// get deployment
	deployment, err := c.client.AppsV1().Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
	if err != nil {
		return
	}
	// update job
	image := imageName + ":" + imageTag
	deployment.Spec.Template.Spec.Containers[0].Image = image
	_, err = c.client.AppsV1().Deployments(namespace).Update(deployment)
	if err != nil {
		return
	}
	return
}

func (c *deployment) GetStatus(namespace string, deploymentName string) (status *DeploymentStatus, err error) {
	if c.err != nil {
		return nil, c.err
	}
	deployment, err := c.client.AppsV1().Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
	if err != nil {
		return
	}
	status.Replicas = int(deployment.Status.Replicas)
	status.AvailableReplicas = int(deployment.Status.AvailableReplicas)
	status.UnavailableReplicas = int(deployment.Status.UnavailableReplicas)
	status.UpdatedReplicas = int(deployment.Status.UpdatedReplicas)
	status.ReadyReplicas = int(deployment.Status.ReadyReplicas)
	conditions := deployment.Status.Conditions
	status.ConditionType = string(conditions[len(conditions)-1].Type)
	status.ConditionStatus = string(conditions[len(conditions)-1].Status)
	status.LastUpdateTime = conditions[len(conditions)-1].LastUpdateTime.Format("2006-01-02 15:04:05")
	status.Reason = conditions[len(conditions)-1].Reason
	status.Message = conditions[len(conditions)-1].Message
	return
}

func (c *deployment) Create(input string) (err error) {
	if c.err != nil {
		return c.err
	}
	deployment := new(v1.Deployment)
	err = yaml.Unmarshal([]byte(input), deployment)
	if err != nil {
		return
	}
	// output for debug
	_, err = yaml.Marshal(deployment)
	if err != nil {
		return
	}
	namespace := deployment.Namespace
	_, err = c.client.AppsV1().Deployments(namespace).Create(deployment)
	if err != nil {
		return
	}
	return
}

func (c *deployment) Delete(namespace string, deploymentName string) (err error) {
	if c.err != nil {
		return c.err
	}
	err = c.client.AppsV1().Deployments(namespace).Delete(deploymentName, &metav1.DeleteOptions{})
	if err != nil {
		return
	}
	return
}

func (c *deployment) Update(input string) (err error) {
	if c.err != nil {
		return c.err
	}
	deployment := new(v1.Deployment)
	err = yaml.Unmarshal([]byte(input), deployment)
	if err != nil {
		return
	}
	// output for debug
	_, err = yaml.Marshal(deployment)
	if err != nil {
		return
	}
	namespace := deployment.Namespace
	_, err = c.client.AppsV1().Deployments(namespace).Update(deployment)
	if err != nil {
		return
	}
	return
}

func (c *deployment) Get(namespace string, name string) (ret string, err error) {
	if c.err != nil {
		return "", c.err
	}
	deployment, err := c.client.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
	d, err := json.Marshal(deployment)
	if err != nil {
		return "", err
	}
	ret = string(d)
	return
}