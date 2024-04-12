package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func GetCRD(client *dynamic.DynamicClient) (*unstructured.UnstructuredList, error) {
	resource, _ := client.Resource(schema.GroupVersionResource{"node.nodedeploy", "v1", "nodedeploys"}).List(context.TODO(), metav1.ListOptions{})
	if resource == nil {
		return nil, errors.New("Not find")
	}
	return resource, nil
}
func CreateCRD(dynamicClient *dynamic.DynamicClient) {
	//使用scheme的包带入gvr
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	deploymentName := "ss"
	replicas := 1
	image := "nginx"
	//定义结构化数据
	deploymnet := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": deploymentName,
			},
			"spec": map[string]interface{}{
				"replicas": replicas,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": "demo",
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": "demo",
						},
					},
					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name":  "web",
								"image": image,
								"ports": []map[string]interface{}{
									{
										"name":          "http",
										"protocol":      "TCP",
										"containerPort": 80,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	create, err := dynamicClient.Resource(deploymentRes).Namespace("default").Create(context.TODO(), deploymnet, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(create)

}
func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	//kubeconfig = "C:\\Users\\huan\\.kube\\conf"

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	// 实例化对象
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	crd, _ := GetCRD(dynamicClient)
	fmt.Println(crd)

}
