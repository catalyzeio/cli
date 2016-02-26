package jobs

import (
	"fmt"

	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func (j *SJobs) List(svcID string, page, pageSize int) (*[]models.Job, error) {
	headers := httpclient.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs?pageNumber=%d&pageSize=%d", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, svcID, page, pageSize), headers)
	if err != nil {
		return nil, err
	}
	var jobs []models.Job
	err = httpclient.ConvertResp(resp, statusCode, &jobs)
	if err != nil {
		return nil, err
	}
	return &jobs, nil
}
