package vars

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/httpclient"
)

func CmdUnset(key string, iv IVars) error {
	err := iv.Unset(key)
	if err != nil {
		return err
	}
	logrus.Println("Unset.")
	return nil
}

// Unset deletes an environment variable. Any changes to environment variables
// will not take effect until the service is redeployed by pushing new code
// or via `catalyze redeploy`.
func (v *SVars) Unset(variable string) error {
	headers := httpclient.GetHeaders(v.Settings.SessionToken, v.Settings.Version, v.Settings.Pod)
	resp, statusCode, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/env/%s", v.Settings.PaasHost, v.Settings.PaasHostVersion, v.Settings.EnvironmentID, v.Settings.ServiceID, variable), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
