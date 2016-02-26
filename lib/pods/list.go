package pods

import (
	"fmt"

	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func (p *SPods) List() (*[]models.Pod, error) {
	headers := httpclient.GetHeaders(p.Settings.SessionToken, p.Settings.Version, p.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/pods", p.Settings.PaasHost, p.Settings.PaasHostVersion), headers)
	if err != nil {
		return nil, err
	}
	var podWrapper models.PodWrapper
	err = httpclient.ConvertResp(resp, statusCode, &podWrapper)
	if err != nil {
		return nil, err
	}
	return podWrapper.Pods, nil
}
