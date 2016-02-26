package sites

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func CmdCreate(name, serviceName, hostname string, is ISites, iservices services.IServices) error {
	upstreamService, err := iservices.RetrieveByLabel(serviceName)
	if err != nil {
		return err
	}
	if upstreamService == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"", serviceName)
	}

	serviceProxy, err := iservices.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}

	site, err := is.Create(name, hostname, upstreamService.ID, serviceProxy.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Created '%s'", site.Name)
	return nil
}

func (s *SSites) Create(name, cert, upstreamServiceID, svcID string) (*models.Site, error) {
	site := models.Site{
		Name:            name,
		Cert:            cert,
		UpstreamService: upstreamServiceID,
	}
	b, err := json.Marshal(site)
	if err != nil {
		return nil, err
	}
	headers := httpclient.GetHeaders(s.Settings.SessionToken, s.Settings.Version, s.Settings.Pod)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/sites", s.Settings.PaasHost, s.Settings.PaasHostVersion, s.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return nil, err
	}
	var createdSite models.Site
	err = httpclient.ConvertResp(resp, statusCode, &createdSite)
	if err != nil {
		return nil, err
	}
	return &createdSite, nil
}
