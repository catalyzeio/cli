package maintenance

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
)

func CmdDisable(svcName string, im IMaintenance, is services.IServices) error {
	upstreamService, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if upstreamService == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	if upstreamService.Type != "code" {
		return fmt.Errorf("Maintenance mode can only be disabled for code services, not %s services", upstreamService.Type)
	}

	serviceProxy, err := is.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}

	svcMaintenance, err := im.List(serviceProxy.ID)
	if err != nil {
		return err
	}
	enabled := false
	for _, mm := range *svcMaintenance {
		if mm.UpstreamID == upstreamService.ID {
			enabled = true
		}
	}
	if !enabled {
		return fmt.Errorf("Maintenance mode is not currently enabled for the service %s", svcName)
	}

	err = im.Disable(serviceProxy.ID, upstreamService.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Maintenance mode disabled for service %s (ID = %s)", upstreamService.Label, upstreamService.ID)
	return nil
}

func (m *SMaintenance) Disable(svcProxyID, upstreamID string) error {
	headers := m.Settings.HTTPManager.GetHeaders(m.Settings.SessionToken, m.Settings.Version, m.Settings.Pod, m.Settings.UsersID)
	resp, statusCode, err := m.Settings.HTTPManager.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/maintenance?upstream=%s", m.Settings.PaasHost, m.Settings.PaasHostVersion, m.Settings.EnvironmentID, svcProxyID, upstreamID), headers)
	if err != nil {
		return err
	}
	return m.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
