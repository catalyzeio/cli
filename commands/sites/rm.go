package sites

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func CmdRm(name string, is ISites, iservices services.IServices) error {
	serviceProxy, err := iservices.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}
	sites, err := is.List(serviceProxy.ID)
	if err != nil {
		return err
	}
	var site *models.Site
	for _, s := range *sites {
		if s.Name == name {
			site = &s
		}
	}
	if site == nil {
		return fmt.Errorf("Could not find a site with the name \"%s\"", name)
	}
	err = is.Rm(site.ID, serviceProxy.ID)
	if err != nil {
		return err
	}
	logrus.Println("Site removed")
	return nil
}

func (s *SSites) Rm(siteID int, svcID string) error {
	headers := httpclient.GetHeaders(s.Settings.SessionToken, s.Settings.Version, s.Settings.Pod)
	resp, statusCode, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/sites/%d", s.Settings.PaasHost, s.Settings.PaasHostVersion, s.Settings.EnvironmentID, svcID, siteID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
