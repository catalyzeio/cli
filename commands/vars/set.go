package vars

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
)

func CmdSet(svcName string, variables []string, fileName string, iv IVars, is services.IServices) error {
	var envVarsMap map[string]string
	var err error
	if fileName != "" {
		if _, err = os.Stat(fileName); err != nil {
			return fmt.Errorf("A file does not exist at path '%s'", fileName)
		}
		fileData, err := ioutil.ReadFile(fileName)
		if err != nil {
			return err
		}
		envVarsMap, err = parseFileData(fileData)
		if err != nil {
			return err
		}
	} else {
		envVarsMap, err = parseKeyValue(variables)
		if err != nil {
			return err
		}
	}
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services\" command.", svcName)
	}

	err = iv.Set(service.ID, envVarsMap)
	if err != nil {
		return err
	}
	logrus.Printf("Set. For these environment variables to take effect, you will need to redeploy your service with \"catalyze redeploy %s\"", svcName)
	return nil
}

func parseFileData(fileData []byte) (map[string]string, error) {
	envVarsMap, err := parseYAML(fileData)
	if err != nil {
		envVarsMap, err = parseJSON(fileData)
		if err != nil {
			variables := strings.Split(string(fileData), "\n")
			envVarsMap, err = parseKeyValue(variables)
			if err != nil {
				return nil, errors.New("Invalid file format. Specified file must be contain key-value pairs in YAML, JSON, or KEY=VALUE format.")
			}
		}
	}
	return envVarsMap, nil
}

func parseYAML(fileData []byte) (map[string]string, error) {
	data := map[string]string{}
	err := yaml.Unmarshal(fileData, &data)
	return data, err
}

func parseJSON(fileData []byte) (map[string]string, error) {
	data := map[string]string{}
	err := json.Unmarshal(fileData, &data)
	return data, err
}

func parseKeyValue(variables []string) (map[string]string, error) {
	data := map[string]string{}
	r := regexp.MustCompile("^[a-zA-Z_]+[a-zA-Z0-9_]*$")
	for _, envVar := range variables {
		if len(envVar) == 0 {
			continue
		}
		pieces := strings.SplitN(envVar, "=", 2)
		if len(pieces) != 2 {
			return nil, fmt.Errorf("Invalid variable format. Expected <key>=<value> but got %s", envVar)
		}
		name, value := pieces[0], pieces[1]
		if !r.MatchString(name) {
			return nil, fmt.Errorf("Invalid environment variable name '%s'. Environment variable names must only contain letters, numbers, and underscores and must not start with a number.", name)
		}
		data[name] = value
	}
	return data, nil
}

// Set adds a new environment variables or updates the value of an existing
// environment variables. Any changes to environment variables will not take
// effect until the service is redeployed by pushing new code or via
// `catalyze redeploy`.
func (v *SVars) Set(svcID string, envVarsMap map[string]string) error {
	b, err := json.Marshal(envVarsMap)
	if err != nil {
		return err
	}
	headers := v.Settings.HTTPManager.GetHeaders(v.Settings.SessionToken, v.Settings.Version, v.Settings.Pod, v.Settings.UsersID)
	resp, statusCode, err := v.Settings.HTTPManager.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/env", v.Settings.PaasHost, v.Settings.PaasHostVersion, v.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return err
	}
	return v.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
