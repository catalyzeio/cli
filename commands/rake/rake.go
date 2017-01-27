package rake

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
)

func CmdRake(svcName, taskName string, ir IRake, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	logrus.Printf("Executing Rake task: %s", taskName)
	err = ir.Run(taskName, service.ID)
	if err != nil {
		return err
	}
	logrus.Println("Rake task output viewable in your logging dashboard")
	return nil
}

// Run executes a rake task. This is only applicable for ruby-based
// applications.
func (r *SRake) Run(taskName, svcID string) error {
	rakeTask := map[string]string{
		"command": taskName,
	}
	b, err := json.Marshal(rakeTask)
	if err != nil {
		return err
	}
	headers := r.Settings.HTTPManager.GetHeaders(r.Settings.SessionToken, r.Settings.Version, r.Settings.Pod, r.Settings.UsersID)
	resp, statusCode, err := r.Settings.HTTPManager.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/rake", r.Settings.PaasHost, r.Settings.PaasHostVersion, r.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return err
	}
	return r.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
