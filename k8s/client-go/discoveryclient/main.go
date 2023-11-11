package main

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
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

	// 新建discoveryClient
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}
	// 获取所有的分组和资源数据
	apiGroups, APIResourceListSlice, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		panic(err)
	}
	fmt.Printf("APIGroup:\n\n %v\n\n", apiGroups)
	for _, resourceList := range APIResourceListSlice {
		version := resourceList.GroupVersion
		fmt.Printf("%s", version)

		//把字符串转换为数据结构
		groupVersion, err := schema.ParseGroupVersion(version)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%v", groupVersion)

	}

}
