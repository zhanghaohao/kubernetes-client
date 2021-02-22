# kubernetes-client
kubernetes客户端，通过加载kubeconfig来和kubernetes集群通讯。支持多集群，公共资源对象抽象，可自定义资源对象方法。
## 安装使用
### go vendor方式
这种方式就是使用vendor目录来管理依赖包。
#### 创建一个golang项目 
`mkdir -p test/src/test`
#### 创建go.mod
`cd test/src/test`     
`go mod init test`
#### 创建main函数
`touch test.go`

test.go文件如下：
```golang
package main

import (
	"fmt"
	k8sCli "github.com/zhanghaohao/kubernetes-client"
)

func main()  {
	clients, err := k8sCli.NewFromDefaultKubeconfigPaths()
	if err != nil {
		panic(err)
	}
	namespace := "default"
	deploymentName := "you-want-delete"
	err = clients.GetClient("default").CommonResourceObject(k8sCli.ResourceObjectType(k8sCli.KubernetesDeployment)).Delete(namespace, deploymentName)
	if err != nil {
		fmt.Println(err)
		return
	}
}
```
#### 配置GOPROXY
由于国内下载golang依赖包很慢，而且有些包还无法下载，所以通常我们都需要先配置GOPROXY。    
`export GOPROXY=https://goproxy.io`    
为了一劳永逸，你可以把上面这一行放在~/.bash_profile里面，这样每次打开终端就会自动生效。
#### 下载kubernetes-client及其他依赖包
`go mod vendor`
#### 修改依赖包版本
默认下载的依赖包都是latest版本，有些依赖包版本我们需要特定的，比如`k8s.io/client-go v0.15.7`
```
module test

go 1.13

require (
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/zhanghaohao/kubernetes-client v0.0.0-20210221043837-4c14a58e7c54
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/oauth2 v0.0.0-20210220000619-9bb904979d93 // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	k8s.io/client-go v0.15.7 // indirect
)
```
#### 再次重新下载依赖包
`go mod vendor`
生成的依赖包会自动放在vendor目录下面
#### 编译运行程序
`go build -o bin/test test.go`


### go modules方式

## 构建多集群k8s客户端
支持两种方式构建多集群k8s客户端
### 通过指定kubeconfig文件路径
```golang
kubeConfigPaths := []k8sCli.KubeConfigPath{
		{
			ClusterName: "cluster1",
			Path: "/kubeconfig/path1",
		},
		{
			ClusterName: "cluster2",
			Path: "/kubeconfig/path2",
		},
	}
clients, err := k8sCli.NewFromKubeconfigPaths(kubeConfigPaths)
if err != nil {
	fmt.Println(err)
	return
}
```
### 通过kubeconfig结构体
## k8s公共资源对象
有k8s公共资源对象`CommonResourceObject`，里面包含`Create`,`Update`,`Delete`,`Get`等方法，这样你在不知道要操作的是哪种资源对象的时候可以不用写`if ... else ...`来判断资源对象类型了，代码更加简洁和高效。
## 自定义资源对象的方法
本代码里面含有自定义资源对象的方法，比如deployment对象有`GetStatus`方法：
```golang
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
```
你也可以自行扩展你需要的方法:laughing:

