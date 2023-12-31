package main

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"
)

// 集群内 创建Client
func main() {
	config, err := rest.InClusterConfig()
	// 创建 clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		//通过省略命名空间获取所有命名空间中的pod
		//或者指定namespace来获取特定命名空间中的pod
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		//错误处理的例子:
		// -使用helper函数，例如:errors.IsNotFound()
		// -和/或转换为StatusError，并使用它的属性，如ErrStatus。消息
		_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "test-node-local-dns", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod test-node-local-dns not found in default namespace\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found test-node-local-dns pod in default namespace\n")
		}

		time.Sleep(10 * time.Second)
	}
}
