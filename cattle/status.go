package cattle


import (
	"context"
	_ "context"
	"fmt"
	_ "k8s.io/api/apps/v1"
	_ "time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)


func GetStatus(vb bool) (string, error) {
	var p string
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	if vb != true {
		pods, err := clientset.CoreV1().Pods("cattle-system").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		for _, pod := range pods.Items {
		p += fmt.Sprintf("| %-40s| %-8s| %-15s|\n", pod.Name, pod.Status.Phase, pod.Spec.NodeName)
		}
		p += fmt.Sprintf("----\nTotal: %d", len(pods.Items))

		return p, err
	} else {
		pods, err := clientset.CoreV1().Pods("cattle-system").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		for _, pod := range pods.Items {
			p += fmt.Sprintf("| %-40s| %-8s| %-15s|\n", pod.Name, pod.Status.Phase, pod.Spec.NodeName)
		}
		p += fmt.Sprintf("----\nTotal: %d \n", len(pods.Items))

		deployments, err := clientset.AppsV1().Deployments("cattle-system").List(context.TODO(), metav1.ListOptions{})
        if err != nil {
        	panic(err)
		}
		for _, deployments := range deployments.Items {
			p += fmt.Sprintf("| %-40s| %-2d| %-2d|", deployments.Name, deployments.Status.ReadyReplicas, deployments.Status.Replicas)
		}
		return p, err
	}
	return "Error", err
}

