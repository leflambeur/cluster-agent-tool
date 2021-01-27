package rancher

import (
	"fmt"

	"github.com/rancher/log"

	clusterv3 "github.com/rancher/types/client/cluster/v3"
	mgmtv3 "github.com/rancher/types/client/management/v3"
	mgmtv3public "github.com/rancher/types/client/management/v3public"
)

type Server struct {
	URL      string
	Token    string
	Username string
	Password string

	V3PublicClient *mgmtv3public.Client
	V3Client       *mgmtv3.Client
	ClusterClient  *clusterv3.Client

	AuthProvider string
}

func NewServer(insecure bool, cattleServer string) (*Server, error) {
	// TODO: Create a temporary directory to work during the session
	s := &Server{}
	s.URL = cattleServer
	if len(s.URL) == 0 {
		fmt.Println("No Rancher URL Detected!")
		url, err := askForRancherServerDetails()
		if err != nil {
			return nil, fmt.Errorf("error getting server details: %v", err)
		}
		s.URL = url
	}

	log.Debugf("Rancher URL: %v", s.URL)

	// TODO: Remove trailing slash. Check if /v3 is also provided
	// TODO: Do a /ping check and also see if self signed certs are being used

	if err := s.initV3PublicClient(insecure); err != nil {
		return nil, fmt.Errorf("error initializing v3public client: %v", err)
	}
	if err := s.doLogin(); err != nil {
		return nil, fmt.Errorf("error performing login: %v", err)
	}
	// initialize clients
	if err := s.initClients(insecure); err != nil {
		return nil, fmt.Errorf("error initializing clients: %v", err)
	}
	if err := s.figureOutManagementPlaneDetails(); err != nil {
		return nil, fmt.Errorf("error figuring management plane details: %v", err)
	}

	return s, nil
}
