package pod

import (
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"encoding/json"
)

type pod struct {
	client 					*kubernetes.Clientset
	err 					error
}

type PodInfo struct {
	PodName 				string
	Status 					string
	HostIP 					string
	PodIP 					string
	StartTime 				string
	Containers 				[]ContainerInfo
}

type ContainerInfo struct {
	Name 					string
	State 					string
	Ready 					bool
	RestartCount			int
	Image  					string
	ImageID 				string
	ContainerID 			string
}

type Pod interface {
	SetErr(err error)
	Create(input string) (err error)
	Delete(namespace string, name string) (err error)
	Update(input string) (err error)
	Get(namespace string, name string) (ret string, err error)
	ListPods(namespace string) (podList []PodInfo, err error)
	GetLogs(namespace string, podName string) (logs string, err error)
}

func NewForClient(client *kubernetes.Clientset) *pod {
	return &pod{
		client: client,
	}
}

func (c *pod) SetErr(err error)  {
	c.err = err
}

func (c *pod) Create(input string) (err error) {
	//todo
	return
}

func (c *pod) Update(input string) (err error) {
	//todo
	return
}

func (c *pod) Delete(namespace string, name string) (err error) {
	//todo
	return
}

func (c *pod) Get(namespace string, name string) (ret string, err error) {
	if c.err != nil {
		return "", c.err
	}
	pod, err := c.client.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
	d, err := json.Marshal(pod)
	if err != nil {
		return "", err
	}
	ret = string(d)
	return
}

func (c *pod) ListPods(namespace string) (podList []PodInfo, err error) {
	if c.err != nil {
		return nil, c.err
	}
	pods, err := c.client.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		return
	}
	//body, err := json.Marshal(pods)
	//if err != nil {
	//	logger.Error.Println(err)
	//	return
	//}
	//var out bytes.Buffer
	//err = json.Indent(&out, body, "", "    ")
	//if err != nil {
	//	logger.Error.Println(err)
	//	return
	//}
	//logger.Info.Println(out.String())
	for _, e := range pods.Items {
		var pod PodInfo
		pod.PodName = e.Name
		pod.Status = string(e.Status.Phase)
		pod.HostIP = e.Status.HostIP
		pod.PodIP = e.Status.PodIP
		pod.StartTime = e.Status.StartTime.Format("2006-01-02 15:04:05")
		for _, c := range e.Status.ContainerStatuses {
			var container ContainerInfo
			container.Name = c.Name
			container.State = c.State.String()
			container.Ready = c.Ready
			container.RestartCount = int(c.RestartCount)
			container.Image = c.Image
			container.ImageID = c.ImageID
			container.ContainerID = c.ContainerID
			pod.Containers = append(pod.Containers, container)
		}
		podList = append(podList, pod)
	}
	return
}

func (c *pod) GetLogs(namespace string, podName string) (logs string, err error) {
	if c.err != nil {
		return "", c.err
	}
	opts := &corev1.PodLogOptions{
		Timestamps: true,
	}
	resp, err := c.client.CoreV1().Pods(namespace).GetLogs(podName, opts).DoRaw()
	if err != nil {
		return
	}
	logs = string(resp)
	return
}