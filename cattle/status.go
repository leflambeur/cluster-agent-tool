package cattle

import (
	"cluster-agent-tool/rancher"
	"context"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Env struct {
	Node             string
	Server           string
	AgentCAChecksum  string
	ClusterID        string
	Deployment       string
	ServerCAChecksum string
}

func PrintStatus(vb bool) (string, error) {
	var p string
	if vb != true {
		podStatus, err := PodStatus()
		p += fmt.Sprintf("\nAgent Pod Status:\n\n%s\n", podStatus)
		return p, err
	} else {
		podStatus, err := PodStatus()
		envStatus, s, err := EnvStatus()
		p += fmt.Sprintf("\nServer: %s\nToken: %s\n\nCluster: %s\nNode: %s\n\nLocal CA Checksum: %s\nRancher CA Checksum: %s\n\nAgent Deployment URL: %s\nAgent Pod Status:\n\n%s\n",
			envStatus.Server, s.Token, envStatus.ClusterID, envStatus.Node, envStatus.AgentCAChecksum, envStatus.ServerCAChecksum, envStatus.Deployment, podStatus)
		return p, err
	}
}

func NewRow(rowcount int) string {
	newline1 := "-"
	res1 := strings.Repeat(newline1, rowcount)
	return res1
}

func PodStatus() (string, error) {
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

func EnvStatus() (Env, *rancher.Server, error) {
	cattleNode := GetEnvVar("CATTLE_NODE_NAME")
	cattleServer := GetEnvVar("CATTLE_SERVER")
	cattleCAChecksum := GetEnvVar("CATTLE_CA_CHECKSUM")
	s, err := rancher.NewServer(true, cattleServer)
	if err != nil {
		panic(err)
	}
	cattleCluster, err := GetClusterID(s, cattleNode)
	cattleDeployment, err := GetDeploymentURL(s, cattleCluster)
	serverCAChecksum, err := GetServerCAChecksum(s)
	currCattleEnv := Env{cattleNode, cattleServer, cattleCAChecksum,
		cattleCluster, cattleDeployment, serverCAChecksum}
	return currCattleEnv, s, err
}
