package rancher

import (
	"fmt"

	normanclientbase "github.com/rancher/norman/clientbase"
	v3 "github.com/rancher/types/client/management/v3"
	v3public "github.com/rancher/types/client/management/v3public"
)

func (s *Server) initV3PublicClient(insecure bool) error {
	var err error
	opts := &normanclientbase.ClientOpts{
		URL:      fmt.Sprintf("%v/v3-public", s.URL),
		Insecure: insecure,
	}

	// TODO: Deal with self signed certs/ Insecure Client
	s.V3PublicClient, err = v3public.NewClient(opts)
	if err != nil {
		return fmt.Errorf("error creating v3public client: %v", err)
	}

	return nil
}

func (s *Server) initClients(insecure bool) error {
	var err error
	opts := &normanclientbase.ClientOpts{
		URL:      fmt.Sprintf("%v/v3", s.URL),
		TokenKey: s.Token,
		Insecure: insecure,
	}

	// TODO: Deal with self signed certs/ Insecure Client
	s.V3Client, err = v3.NewClient(opts)
	if err != nil {
		return fmt.Errorf("error creating v3 client: %v", err)
	}

	return nil
}
