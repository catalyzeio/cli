package maintenance

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
)

func CmdEnable(svcName string, im IMaintenance, is services.IServices) error {
	upstreamService, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if upstreamService == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	if upstreamService.Type != "code" {
		return fmt.Errorf("Maintenance mode can only be enabled for code services, not %s services", upstreamService.Type)
	}

	serviceProxy, err := is.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}

	err = im.Enable(serviceProxy.ID, upstreamService.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Maintenance mode enabled for service %s (ID = %s)", upstreamService.Label, upstreamService.ID)
	return nil
}

func (m *SMaintenance) Enable(svcProxyID, upstreamID string) error {
	body := map[string]string{
		"upstream": upstreamID,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	headers := m.Settings.HTTPManager.GetHeaders(m.Settings.SessionToken, m.Settings.Version, m.Settings.Pod, m.Settings.UsersID)
	resp, statusCode, err := m.Settings.HTTPManager.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/maintenance", m.Settings.PaasHost, m.Settings.PaasHostVersion, m.Settings.EnvironmentID, svcProxyID), headers)
	if err != nil {
		return err
	}
	return m.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
