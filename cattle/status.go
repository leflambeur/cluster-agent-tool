package cattle


import (
	"context"
	_ "context"
	"fmt"
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
		p += fmt.Sprintf("%d", len(pods.Items))
		return p, err
	}
	return "Error", err
}

