package releases

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
)

func CmdRm(svcName, releaseName string, ir IReleases, is services.IServices) error {
	if strings.ContainsAny(releaseName, config.InvalidChars) {
		return fmt.Errorf("Invalid release name. Names must not contain the following characters: %s", config.InvalidChars)
	}
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	err = ir.Rm(releaseName, service.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Release '%s' has been successfully removed.", releaseName)
	return nil
}

func (r *SReleases) Rm(releaseName, svcID string) error {
	headers := r.Settings.HTTPManager.GetHeaders(r.Settings.SessionToken, r.Settings.Version, r.Settings.Pod, r.Settings.UsersID)
	resp, statusCode, err := r.Settings.HTTPManager.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/releases/%s", r.Settings.PaasHost, r.Settings.PaasHostVersion, r.Settings.EnvironmentID, svcID, releaseName), headers)
	if err != nil {
		return err
	}
	return r.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
