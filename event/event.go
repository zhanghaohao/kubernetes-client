package event

import (
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"encoding/json"
)

type event struct {
	client 					*kubernetes.Clientset
	err 					error
}

type EventFieldSelector struct {
	Name 					string
	Kind 					string
}

type EventInfo struct {
	Kind 					string
	Name 					string
	Reason 					string
	Message 				string
	LastTimestamp 			string
	Count 					int
	Type 					string
}

type Event interface {
	SetErr(err error)
	Create(input string) (err error)
	Delete(namespace string, name string) (err error)
	Update(input string) (err error)
	Get(namespace string, name string) (ret string, err error)
	List(namespace string, fieldSelector *EventFieldSelector) (eventList []EventInfo, err error)
}

func NewForClient(client *kubernetes.Clientset) *event {
	return &event{
		client: client,
	}
}

func (c *event) SetErr(err error)  {
	c.err = err
}

func (c *event) Create(input string) (err error) {
	//todo
	return
}

func (c *event) Update(input string) (err error) {
	//todo
	return
}

func (c *event) Delete(namespace string, name string) (err error) {
	//todo
	return
}

func (c *event) Get(namespace string, name string) (ret string, err error) {
	if c.err != nil {
		return "", c.err
	}
	event, err := c.client.CoreV1().Events(namespace).Get(name, metav1.GetOptions{})
	d, err := json.Marshal(event)
	if err != nil {
		return "", err
	}
	ret = string(d)
	return
}

func (c *event) List(namespace string, fieldSelector *EventFieldSelector) (eventList []EventInfo, err error) {
	if c.err != nil {
		return nil, c.err
	}
	/*
	filter by fieldSelector
	*/
	var filter string
	var opts metav1.ListOptions
	if fieldSelector == nil {
		filter = ""
	} else {
		var name, kind *string
		if len(fieldSelector.Name) == 0 {
			name = nil
		} else {
			name = &fieldSelector.Name
		}
		if len(fieldSelector.Kind) == 0 {
			kind = nil
		} else {
			kind = &fieldSelector.Kind
		}
		filter = c.client.CoreV1().Events(namespace).GetFieldSelector(name, &namespace, kind, nil).String()
		opts = metav1.ListOptions{
			FieldSelector: filter,
		}
	}
	//logger.Info.Println(filter)
	//logger.Info.Println(opts)
	events, err := c.client.CoreV1().Events(namespace).List(opts)
	if err != nil {
		return
	}
	for _, event := range events.Items {
		eventList = append(eventList, EventInfo{
			Kind: event.Kind,
			Name: event.Name,
			Reason: event.Reason,
			Message: event.Message,
			LastTimestamp: event.LastTimestamp.Format("2006-01-02 15:04:05"),
			Count: int(event.Count),
			Type: event.Type,
		})
	}
	return
}