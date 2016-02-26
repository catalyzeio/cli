package certs

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
)

func CmdRm(hostname string, ic ICerts, is services.IServices) error {
	service, err := is.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}
	hostname = strings.Replace(hostname, "*", "star", -1)
	err = ic.Rm(hostname, service.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Removed '%s'", hostname)
	return nil
}

func (c *SCerts) Rm(hostname, svcID string) error {
	headers := httpclient.GetHeaders(c.Settings.SessionToken, c.Settings.Version, c.Settings.Pod)
	resp, statusCode, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/certs/%s", c.Settings.PaasHost, c.Settings.PaasHostVersion, c.Settings.EnvironmentID, svcID, hostname), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
