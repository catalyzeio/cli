package services

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/olekukonko/tablewriter"
)

// CmdServices lists the names of all services for an environment.
func CmdServices(is IServices) error {
	svcs, err := is.List()
	if err != nil {
		return err
	}
	data := [][]string{{"NAME", "RAM (GB)", "CPU", "STORAGE (GB)"}}
	for _, s := range *svcs {
		data = append(data, []string{s.Label, fmt.Sprintf("%d", s.Size.RAM), fmt.Sprintf("%d", s.Size.CPU), fmt.Sprintf("%d", s.Size.Storage)})
	}

	table := tablewriter.NewWriter(logrus.StandardLogger().Out)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.AppendBulk(data)
	table.Render()
	return nil
}

func (s *SServices) List() (*[]models.Service, error) {
	return s.ListByEnvID(s.Settings.EnvironmentID, s.Settings.Pod)
}

func (s *SServices) ListByEnvID(envID, podID string) (*[]models.Service, error) {
	headers := httpclient.GetHeaders(s.Settings.SessionToken, s.Settings.Version, podID)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services", s.Settings.PaasHost, s.Settings.PaasHostVersion, envID), headers)
	if err != nil {
		return nil, err
	}
	var services []models.Service
	err = httpclient.ConvertResp(resp, statusCode, &services)
	if err != nil {
		return nil, err
	}
	return &services, nil
}

func (s *SServices) Retrieve(svcID string) (*models.Service, error) {
	headers := httpclient.GetHeaders(s.Settings.SessionToken, s.Settings.Version, s.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s", s.Settings.PaasHost, s.Settings.PaasHostVersion, s.Settings.EnvironmentID, s.Settings.ServiceID), headers)
	if err != nil {
		return nil, err
	}
	var service models.Service
	err = httpclient.ConvertResp(resp, statusCode, &service)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (s *SServices) RetrieveByLabel(label string) (*models.Service, error) {
	services, err := s.List()
	if err != nil {
		return nil, err
	}
	var service *models.Service
	for _, s := range *services {
		if s.Label == label {
			service = &s
			break
		}
	}
	return service, nil
}
