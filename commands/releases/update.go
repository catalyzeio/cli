package releases

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/models"
)

func CmdUpdate(svcName, releaseName, notes, newReleaseName string, ir IReleases, is services.IServices) error {
	if strings.ContainsAny(releaseName, config.InvalidChars) {
		return fmt.Errorf("Invalid existing release name. Names must not contain the following characters: %s", config.InvalidChars)
	}
	if strings.ContainsAny(newReleaseName, config.InvalidChars) {
		return fmt.Errorf("Invalid updated release name. Names must not contain the following characters: %s", config.InvalidChars)
	}
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	err = ir.Update(releaseName, service.ID, notes, newReleaseName)
	if err != nil {
		return err
	}
	logrus.Printf("Release '%s' successfully updated", releaseName)
	return nil
}

func (r *SReleases) Update(releaseName, svcID, notes, newReleaseName string) error {
	rls := models.Release{
		Name:  newReleaseName,
		Notes: notes,
	}
	b, err := json.Marshal(rls)
	if err != nil {
		return err
	}
	headers := r.Settings.HTTPManager.GetHeaders(r.Settings.SessionToken, r.Settings.Version, r.Settings.Pod, r.Settings.UsersID)
	resp, statusCode, err := r.Settings.HTTPManager.Put(b, fmt.Sprintf("%s%s/environments/%s/services/%s/releases/%s", r.Settings.PaasHost, r.Settings.PaasHostVersion, r.Settings.EnvironmentID, svcID, releaseName), headers)
	if err != nil {
		return err
	}
	return r.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
