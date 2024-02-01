package main

import (
	"context"
	"fmt"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/joohwan/k8sipconfig/10.10.102.96")
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	secret, err := clientset.CoreV1().Secrets("kubeedge").Get(context.Background(), "tokensecret", metaV1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Kubeedge JOIN TOKEN:%s", secret.Data["tokendata"])
}
