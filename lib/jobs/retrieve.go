package jobs

import (
	"fmt"

	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func (j *SJobs) Retrieve(jobID, svcID string, includeSpec bool) (*models.Job, error) {
	headers := httpclient.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs/%s?spec=true", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, svcID, jobID), headers)
	if err != nil {
		return nil, err
	}
	var job models.Job
	err = httpclient.ConvertResp(resp, statusCode, &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (j *SJobs) RetrieveByStatus(status string) (*[]models.Job, error) {
	headers := httpclient.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs?status=%s", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, j.Settings.ServiceID, status), headers)
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

func (j *SJobs) RetrieveByType(jobType string, page, pageSize int) (*[]models.Job, error) {
	headers := httpclient.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs?type=%s&pageNumber=%d&pageSize=%d", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, j.Settings.ServiceID, jobType, page, pageSize), headers)
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
