package jobs

import (
	"fmt"

	"github.com/catalyzeio/cli/lib/httpclient"
)

func (j *SJobs) Delete(jobID, svcID string) error {
	headers := httpclient.GetHeaders(j.Settings.SessionToken, j.Settings.Version, j.Settings.Pod)
	resp, statusCode, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs/%s", j.Settings.PaasHost, j.Settings.PaasHostVersion, j.Settings.EnvironmentID, svcID, jobID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
