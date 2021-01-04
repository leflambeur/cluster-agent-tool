package cattle

import (
	"cluster-agent-tool/rancher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/rancher/norman/types"
	"os"
	"strings"
)

//func GetDeployment(){
//	getLogin, s, err := EnvStatus()
//}

func GetEnvVar(cattleType string) string {
	cattleValue := os.Getenv(cattleType)
	return cattleValue
}

func GetClusterID(s *rancher.Server, cattleNode string) (string, error) {

	nodesInfo, err := s.V3Client.Node.List(&types.ListOpts{
		Filters: map[string]interface{}{"nodeName": cattleNode},
	})
	if err != nil {
		panic(err)
	}
	clusterID := fmt.Sprintf("%v", nodesInfo.Data[0].ClusterID)
	return clusterID, err
}

func GetDeploymentURL(s *rancher.Server, cattleCluster string) (string, error) {
	cattleDeployment, err := s.V3Client.ClusterRegistrationToken.List(&types.ListOpts{
		Filters: map[string]interface{}{"clusterID": cattleCluster},
	})
	if err != nil {
		panic(err)
	}

	pDeployment := fmt.Sprintf("%v", cattleDeployment.Data[0].ManifestURL)

	return pDeployment, err
}

func GetServerCAChecksum(s *rancher.Server)(string, error){
	serverCASetting, err := s.V3Client.Setting.List(&types.ListOpts{
		Filters: map[string]interface{}{"name": "cacerts"},
	})
	serverCA := serverCASetting.Data[0].Value
	if serverCA != "" {
		if !strings.HasSuffix(serverCA, "\n") {
			serverCA += "\n"
		}
		digest := sha256.Sum256([]byte(serverCA))
		return hex.EncodeToString(digest[:]),err
	}
	return "No Checksum!", err
}

