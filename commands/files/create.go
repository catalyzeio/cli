package files

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func (f *SFiles) Create(svcID, filePath, name, mode string) (*models.ServiceFile, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	sf := models.ServiceFile{
		Contents:       string(b),
		GID:            0,
		Mode:           mode,
		Name:           name,
		UID:            0,
		EnableDownload: true,
	}
	body, err := json.Marshal(sf)
	if err != nil {
		return nil, err
	}
	headers := httpclient.GetHeaders(f.Settings.SessionToken, f.Settings.Version, f.Settings.Pod)
	resp, statusCode, err := httpclient.Post(body, fmt.Sprintf("%s%s/environments/%s/services/%s/files", f.Settings.PaasHost, f.Settings.PaasHostVersion, f.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return nil, err
	}
	var svcFile models.ServiceFile
	err = httpclient.ConvertResp(resp, statusCode, &svcFile)
	if err != nil {
		return nil, err
	}
	return &svcFile, nil
}
