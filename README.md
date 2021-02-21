# kubernetes-client
kubernetes客户端，通过加载kubeconfig来和kubernetes集群通讯。支持多集群，公共资源对象抽象，可自定义资源对象方法。
## 安装使用
### go vendor方式
#### 创建一个golang项目 
`mkdir -p test/src/test`
#### 创建go.mod
`cd test/src/test`
`go mod init test`
#### 创建main函数
`touch test.go`

test.go文件如下：
```
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
#### 下载kubernetes-client及其他依赖包
`go mod vendor`
#### 修改依赖包版本
默认下载的依赖包都是latest版本，有些依赖包版本我们需要特定的，比如和k8s.io下面的
```
module test

go 1.13

require (
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/zhanghaohao/kubernetes-client v0.0.0-20210221043837-4c14a58e7c54
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/oauth2 v0.0.0-20210220000619-9bb904979d93 // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba // indirect
	k8s.io/client-go {red}{v0.15.7} // indirect
)
```
#### 再次重新下载依赖包
`go mod vendor`
生成的依赖包会自动放在vendor目录下面
#### 编译运行程序
`go build -o bin/test test.go`


### go modules方式
