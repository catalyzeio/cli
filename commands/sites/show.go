package sites

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/forana/simpletable"
)

func CmdShow(name string, is ISites, iservices services.IServices) error {
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
	site, err = is.Retrieve(site.ID, serviceProxy.ID)
	if err != nil {
		return err
	}
	table, err := simpletable.New(simpletable.HeadersForType(models.Site{}), []models.Site{*site})
	if err != nil {
		return err
	}
	table.Write(logrus.StandardLogger().Out)
	return nil
}

func (s *SSites) Retrieve(siteID int, svcID string) (*models.Site, error) {
	headers := httpclient.GetHeaders(s.Settings.SessionToken, s.Settings.Version, s.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/sites/%d", s.Settings.PaasHost, s.Settings.PaasHostVersion, s.Settings.EnvironmentID, svcID, siteID), headers)
	if err != nil {
		return nil, err
	}
	var site models.Site
	err = httpclient.ConvertResp(resp, statusCode, &site)
	if err != nil {
		return nil, err
	}
	return &site, nil
}
