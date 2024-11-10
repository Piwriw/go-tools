package main

import (
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"strings"
)

func UserIndexFunc(obj interface{}) ([]string, error) {
	pod := obj.(*v1.Pod)
	userString := pod.Annotations["users"]
	return strings.Split(userString, ","), nil
}
func main() {
	indexer := cache.NewIndexer(cache.MetaNamespaceKeyFunc, cache.Indexers{"byUser": UserIndexFunc})
	pod1 := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "one", Annotations: map[string]string{"users": "user1"}}}
	pod2 := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "two", Annotations: map[string]string{"users": "user2"}}}
	pod3 := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "three", Annotations: map[string]string{"users": "user2"}}}
	indexer.Add(pod1)
	indexer.Add(pod2)
	indexer.Add(pod3)

	pods, err := indexer.ByIndex("byUser", "user1")
	if err != nil {
		panic(err)
	}
	for _, pod := range pods {
		fmt.Println(pod.(*v1.Pod).Name)
	}
}
