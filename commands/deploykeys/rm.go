package deploykeys

import (
	"fmt"

	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
)

func CmdRm(name, svcName string, private bool, id IDeployKeys, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"", svcName)
	}
	if service.Type != "code" {
		return fmt.Errorf("You can only remove a deploy keys from code services, not %s services", service.Type)
	}
	keyType := "ssh"
	if private {
		keyType = "ssh_private"
	}
	return id.Rm(name, keyType, service.ID)
}

func (d *SDeployKeys) Rm(name, keyType, svcID string) error {
	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod)
	resp, statusCode, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/ssh_keys/%s/type/%s", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, svcID, name, keyType), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
