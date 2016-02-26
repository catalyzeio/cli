package certs

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func CmdList(ic ICerts, is services.IServices) error {
	service, err := is.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}
	certs, err := ic.List(service.ID)
	if err != nil {
		return err
	}
	if certs == nil || len(*certs) == 0 {
		logrus.Println("No certs found")
		return nil
	}
	logrus.Println("NAME")
	for _, cert := range *certs {
		logrus.Println(cert.Name)
	}
	return nil
}

func (c *SCerts) List(svcID string) (*[]models.Cert, error) {
	headers := httpclient.GetHeaders(c.Settings.SessionToken, c.Settings.Version, c.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/certs", c.Settings.PaasHost, c.Settings.PaasHostVersion, c.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return nil, err
	}
	var certs []models.Cert
	err = httpclient.ConvertResp(resp, statusCode, &certs)
	if err != nil {
		return nil, err
	}
	return &certs, nil
}
