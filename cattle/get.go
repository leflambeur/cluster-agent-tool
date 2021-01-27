package cattle

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/leflambeur/cluster-agent-tool/rancher"
	"github.com/rancher/log"
	"github.com/rancher/norman/types"
)

type WriteCounter struct {
	Total uint64
}

func GetDeploymentSetup(apply bool) (string, error) {
	var p string
	getLogin, s, err := EnvStatus()
	p += fmt.Sprintf(
		"\nServer: %s\nToken: %s\n\nCluster: %s\nNode: %s\n\nLocal CA Checksum: %s\nRancher CA Checksum: %s\n\nAgent Deployment URL: %s\n",
		getLogin.Server, s.Token, getLogin.ClusterID, getLogin.Node, getLogin.AgentCAChecksum, getLogin.ServerCAChecksum, getLogin.Deployment)
	workingDirectory, err := createWorkDir("/tmp")
	log.Debugf("created working directory: %v", workingDirectory)
	getDeployment(getLogin.Deployment, workingDirectory, apply)
	if err != nil {
		panic(err)
	}
	p += fmt.Sprintf("\n\nFile Created at %s", workingDirectory)
	return p, err
}

func getDeployment(cattleDeployment string, workingDirectory string, apply bool) error {
	fullPath := fmt.Sprintf("%s/deployment.yaml", workingDirectory)
	fmt.Println()
	resp, err := http.Get(cattleDeployment)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(fullPath + ".tmp")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	counter := &WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		panic(err)
	}

	fmt.Println()
	err = os.Rename(fullPath+".tmp", fullPath)
	if err != nil {
		return err
	}

	if apply == true {
		kubectlApply, err := ApplyDeployment(fullPath)
		fmt.Println(kubectlApply)
		fmt.Println("Apply was True")
		if err != nil {
			panic(err)
		}
	}
	return err
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.printProgress()
	return n, nil
}

func (wc WriteCounter) printProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 50))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

func getEnvVar(cattleType string) string {
	cattleValue := os.Getenv(cattleType)
	return cattleValue
}

func getClusterID(s *rancher.Server, cattleNode string) (string, error) {

	nodesInfo, err := s.V3Client.Node.List(&types.ListOpts{
		Filters: map[string]interface{}{"nodeName": cattleNode},
	})
	if err != nil {
		panic(err)
	}
	clusterID := fmt.Sprintf("%v", nodesInfo.Data[0].ClusterID)
	return clusterID, err
}

func getDeploymentURL(s *rancher.Server, cattleCluster string) (string, error) {
	cattleDeployment, err := s.V3Client.ClusterRegistrationToken.List(&types.ListOpts{
		Filters: map[string]interface{}{"clusterID": cattleCluster},
	})
	if err != nil {
		panic(err)
	}

	pDeployment := fmt.Sprintf("%v", cattleDeployment.Data[0].ManifestURL)

	return pDeployment, err
}

func getServerCAChecksum(s *rancher.Server) (string, error) {
	serverCASetting, err := s.V3Client.Setting.List(&types.ListOpts{
		Filters: map[string]interface{}{"name": "cacerts"},
	})
	serverCA := serverCASetting.Data[0].Value
	if serverCA != "" {
		if !strings.HasSuffix(serverCA, "\n") {
			serverCA += "\n"
		}
		digest := sha256.Sum256([]byte(serverCA))
		return hex.EncodeToString(digest[:]), err
	}
	return "No Checksum!", err
}
