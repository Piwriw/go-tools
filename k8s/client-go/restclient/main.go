package main

import (
	"context"
	"flag"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	// path: /api/v1/namespaces/{namespace}/pods
	config.APIPath = "api"
	// pod的Group是空字符串
	config.GroupVersion = &corev1.SchemeGroupVersion
	// 指定序列化工具
	config.NegotiatedSerializer = scheme.Codecs

	// 创建RESTClient实例
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}
	// 存放podList
	res := &corev1.PodList{}
	namespace := "kube-system"
	err = restClient.Get().Namespace(namespace).
		// 请求资源
		Resource("pods").
		// 指定大小限制和序列化工具
		VersionedParams(&metav1.ListOptions{Limit: 100}, scheme.ParameterCodec).
		Do(context.TODO()).
		// 结果存放
		Into(res)
	if err != nil {
		panic(err)
	}
	for _, item := range res.Items {
		fmt.Printf("%v\t %v\t %v\n",
			item.Namespace,
			item.Status.Phase,
			item.Name)
	}

}
