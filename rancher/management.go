package rancher

import (
	"fmt"
	"github.com/rancher/log"
)

func (s *Server) figureOutManagementPlaneDetails() error {
	clustersInfo, err := s.V3Client.Cluster.List(nil)
	if err != nil {
		return err
	}
	log.Debugf("clusters=%+v", clustersInfo.Data)

	// TODO: handle k3s

	// If local cluster is available from the API, there is no need to ask
	// for RKE cluster kubeconfig

	//// Get the kubeconfig file location of the RKE cluster
	//rkeKubeconfigFileLocation, err := s.getRKEKubeconfigFileLocation()
	//if err != nil {
	//	return err
	//}
	//log.Debugf("rkeKubeconfigFileLocation: %v", rkeKubeconfigFileLocation)

	// Select which cluster(s) to collect logs for
	isLocalClusterPresent := false
	clusterOptions := []string{}
	for _, cluster := range clustersInfo.Data {
		if cluster.Name == "local" {
			isLocalClusterPresent = true
			continue
		}
		// TODO: Figure out what do with unavailable/provisioning clusters
		option := fmt.Sprintf("%v/%v", cluster.Name, cluster.ID)
		clusterOptions = append(clusterOptions, option)
	}
	log.Debugf("isLocalClusterPresent: %v", isLocalClusterPresent)

	// TODO: Check if using client-go will be a problem when working with older k8s clusters.
	// 		 If yes, use kubectl binary directly based on version.
	//		 Check support matrix and see what is the lowest version of k8s supported

	// TODO: Connect to local cluster, grab num of cores, memory, sysctl etc settings
	// 		 and see if they are sufficient

	return nil
}
