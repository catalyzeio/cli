package vars

import (
	"fmt"
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/httpclient"
)

func CmdList(iv IVars) error {
	envVars, err := iv.List()
	if err != nil {
		return err
	}
	var keys []string
	for k := range envVars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		logrus.Printf("%s=%s", key, envVars[key])
	}
	return nil
}

// List lists all environment variables.
func (v *SVars) List() (map[string]string, error) {
	headers := httpclient.GetHeaders(v.Settings.SessionToken, v.Settings.Version, v.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/env", v.Settings.PaasHost, v.Settings.PaasHostVersion, v.Settings.EnvironmentID, v.Settings.ServiceID), headers)
	if err != nil {
		return nil, err
	}
	var envVars map[string]string
	err = httpclient.ConvertResp(resp, statusCode, &envVars)
	if err != nil {
		return nil, err
	}
	return envVars, nil
}
