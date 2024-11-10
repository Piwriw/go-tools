package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func exitWhenTerminate(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		fmt.Println("Received terminate,start shutting down...")
		cancel()
	}()
}

func main() {
	var id string
	flag.StringVar(&id, "id", uuid.New().String(), "id of the leader (required)")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exitWhenTerminate(cancel)

	config, err := clientcmd.BuildConfigFromFlags("", "/Users/joohwan/.kube/config")
	if err != nil {
		panic(err)
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      "example",
			Namespace: "default",
		},
		Client: kubeClient.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				run(ctx)
			},
			OnStoppedLeading: func() {
				fmt.Printf("leader lost %s", id)
				os.Exit(1)
			},
			OnNewLeader: func(identity string) {
				if identity == id {
					return
				}
				fmt.Printf("new leader elected: %s", id)
			},
		},
	})

}
func run(ctx context.Context) {
	fmt.Println("Controller loop...")
	select {}

}
