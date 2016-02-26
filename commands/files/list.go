package files

import (
	"fmt"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

// CmdList lists all service files that are able to be downloaded
// by a member of the environment. Typically service files of interest
// will be on the service_proxy.
func CmdList(svcName string, ifiles IFiles, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the name \"%s\"", svcName)
	}
	files, err := ifiles.List(service.ID)
	if err != nil {
		return err
	}
	if len(*files) == 0 {
		logrus.Println("No service files found")
		return nil
	}
	logrus.Println("NAME")
	for _, sf := range *files {
		logrus.Println(sf.Name)
	}
	return nil
}

func fileModeToRWXString(perms uint64) string {
	permissionString := ""
	binaryString := strconv.FormatUint(perms, 2)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if string(binaryString[len(binaryString)-1-i*3-j]) == "1" {
				switch j {
				case 0:
					permissionString = "x" + permissionString
				case 1:
					permissionString = "w" + permissionString
				case 2:
					permissionString = "r" + permissionString
				}
			} else {
				permissionString = "-" + permissionString
			}
		}
	}
	permissionString = "-" + permissionString // we don't store folders
	return permissionString
}

func (f *SFiles) List(svcID string) (*[]models.ServiceFile, error) {
	headers := httpclient.GetHeaders(f.Settings.SessionToken, f.Settings.Version, f.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/files", f.Settings.PaasHost, f.Settings.PaasHostVersion, f.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return nil, err
	}
	var svcFiles []models.ServiceFile
	err = httpclient.ConvertResp(resp, statusCode, &svcFiles)
	if err != nil {
		return nil, err
	}
	return &svcFiles, nil
}
