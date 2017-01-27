package deploykeys

import (
	"fmt"
	"strings"

	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
)

func CmdRm(name, svcName string, id IDeployKeys, is services.IServices) error {
	if strings.ContainsAny(name, config.InvalidChars) {
		return fmt.Errorf("Invalid SSH key name. Names must not contain the following characters: %s", config.InvalidChars)
	}
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	if service.Type != "code" {
		return fmt.Errorf("You can only remove deploy keys from code services, not %s services", service.Type)
	}
	return id.Rm(name, "ssh", service.ID)
}

func (d *SDeployKeys) Rm(name, keyType, svcID string) error {
	headers := d.Settings.HTTPManager.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod, d.Settings.UsersID)
	resp, statusCode, err := d.Settings.HTTPManager.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/ssh_keys/%s/type/%s", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, svcID, name, keyType), headers)
	if err != nil {
		return err
	}
	return d.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
