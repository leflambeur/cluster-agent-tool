package cattle

import (
	"cluster-agent-tool/rancher"
	"fmt"
	"github.com/rancher/norman/types"
	"os"
)

func GetDeployment(cattleToken string, cattleCluster string, cattleServer string, apply bool) error {
	if len(cattleCluster) == 0 {
		//cattleNode := GetEnvVars("CATTLE_NODE_NAME")
		//cattleCluster := GetClusterID(cattleNode)
	}
	if len(cattleServer) == 0 {
		//cattleServer := GetEnvVars("CATTLE_URL")
	}

	//cattleCAChecksum := GetEnvVars("CATTLE_CA_CHECKSUM")
	return nil
}

func GetEnvVars(cattleType string) string {
	cattleValue := os.Getenv(cattleType)
	return cattleValue
}

func GetClusterID(cattleNode string, cattleServer string) (string, error) {
	s, err := rancher.NewServer(true, cattleServer)
	if err != nil {
		panic(err)
	}
	nodesInfo, err := s.V3Client.Node.List(&types.ListOpts{
		Filters: map[string]interface{}{"nodeName": cattleNode},
	})
	if err != nil {
		panic(err)
	}
	clusterID := fmt.Sprintf("%v", nodesInfo.Data)
	return clusterID, err
}
