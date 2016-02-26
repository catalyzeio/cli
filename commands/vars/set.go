package vars

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/httpclient"
)

func CmdSet(variables []string, iv IVars) error {
	envVarsMap := make(map[string]string, len(variables))
	r := regexp.MustCompile("^[a-zA-Z_]+[a-zA-Z0-9_]*$")
	for _, envVar := range variables {
		pieces := strings.SplitN(envVar, "=", 2)
		if len(pieces) != 2 {
			return fmt.Errorf("Invalid variable format. Expected <key>=<value> but got %s", envVar)
		}
		name, value := pieces[0], pieces[1]
		if !r.MatchString(name) {
			return fmt.Errorf("Invalid environment variable name '%s'. Environment variable names must only contain letters, numbers, and underscores and must not start with a number.", name)
		}
		envVarsMap[name] = value
	}

	err := iv.Set(envVarsMap)
	if err != nil {
		return err
	}
	logrus.Println("Set. For these environment variables to take effect, you will need to redeploy your service with \"catalyze redeploy\"")
	return nil
}

// Set adds a new environment variables or updates the value of an existing
// environment variables. Any changes to environment variables will not take
// effect until the service is redeployed by pushing new code or via
// `catalyze redeploy`.
func (v *SVars) Set(envVarsMap map[string]string) error {
	b, err := json.Marshal(envVarsMap)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(v.Settings.SessionToken, v.Settings.Version, v.Settings.Pod)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/env", v.Settings.PaasHost, v.Settings.PaasHostVersion, v.Settings.EnvironmentID, v.Settings.ServiceID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
