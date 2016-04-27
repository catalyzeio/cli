package files

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func CmdUpload(svcName, fileName, localFilePath string, ifiles IFiles, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services\" command.", svcName)
	}
	file, err := ifiles.Retrieve(fileName, service.ID)
	if err != nil {
		return err
	}
	if file == nil {
		return fmt.Errorf("File with name %s does not exist. Try listing files again by running \"catalyze files list %s\"", fileName, svcName)
	}
	err = ifiles.Update(service.ID, localFilePath, file.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Successfully updated %s", fileName)
	logrus.Printf("To make your changes go live, you must redeploy your service with the \"catalyze redeploy %s\" command", service.Label)
	return nil
}

func (f *SFiles) Update(svcID, filePath string, fileID int) error {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	sf := models.ServiceFile{
		Contents: string(b),
	}
	body, err := json.Marshal(sf)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(f.Settings.SessionToken, f.Settings.Version, f.Settings.Pod, f.Settings.UsersID)
	resp, statusCode, err := httpclient.Put(body, fmt.Sprintf("%s%s/environments/%s/services/%s/files/%d", f.Settings.PaasHost, f.Settings.PaasHostVersion, f.Settings.EnvironmentID, svcID, fileID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
