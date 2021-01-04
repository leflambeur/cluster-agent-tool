package cattle

import (
	"context"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetStatus(vb bool) (string, error) {
	var p string
	if vb != true {
		podStatus, err := GetPodStatus()
		cattleStatus, err := GetCattleEnvValues()
		p += fmt.Sprintf("%s\n\n%s", podStatus, cattleStatus)
		return p, err
	} else {
		var podStatus, err = GetPodStatus()
		p += fmt.Sprintf("%s", podStatus)
		return p, err
	}
}

func NewRow(rowcount int) string {
	newline1 := "-"
	res1 := strings.Repeat(newline1, rowcount)
	return res1
}

func GetPodStatus() (string, error) {
	var ppod string
	var newline string

	newline += fmt.Sprintf("|%-41s|%-9s|%-16s|%-16s|", NewRow(41), NewRow(9), NewRow(16), NewRow(16))
	config, err := rest.InClusterConfig()
	if err != nil {
		 panic(err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	pods, err := clientset.CoreV1().Pods("cattle-system").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	ppod += fmt.Sprintf("%s\n| %-40s| %-8s| %-15s| %-15s|\n%s", newline, "Pod Name", "Status", "Age", "Node", newline)
	for _, pod := range pods.Items {
		podCreationTime := pod.GetCreationTimestamp()
		age := time.Since(podCreationTime.Time).Round(time.Second)

		ppod += fmt.Sprintf("\n| %-40s| %-8s| %-15s| %-15s|\n%s", pod.Name, pod.Status.Phase,
			age.String(), pod.Spec.NodeName, newline)
	}
	ppod += fmt.Sprintf("\n\nTotal: %d \n", len(pods.Items))
	return ppod, err
}

func GetCattleEnvValues() (string, error){
	var pcattle string
	cattleNode := GetEnvVars("CATTLE_NODE_NAME")
	cattleServer := GetEnvVars("CATTLE_SERVER")
	cattleCluster, err := GetClusterID(cattleNode, cattleServer)
	cattleCAChecksum := GetEnvVars("CATTLE_CA_CHECKSUM")
	pcattle += fmt.Sprintf("\nNode: %s\nCluster ID: %s\nRancher Server: %s", cattleCluster, cattleServer, cattleCAChecksum)
	return pcattle, err
}
