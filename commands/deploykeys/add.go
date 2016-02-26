package deploykeys

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func CmdAdd(name, keyPath, svcName string, private bool, id IDeployKeys, is services.IServices) error {
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return fmt.Errorf("A file does not exist at path '%s'", keyPath)
	}
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"", svcName)
	}
	if service.Type != "code" {
		return fmt.Errorf("You can only add deploy keys to code services, not %s services", service.Type)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return err
	}
	keyType := "ssh"
	if private {
		keyType = "ssh_private"
		_, err := id.ParsePrivateKey(key)
		if err != nil {
			return err
		}
	} else {
		_, err := id.ParsePublicKey(key)
		if err != nil {
			return err
		}
	}
	return id.Add(name, keyType, string(key), service.ID)
}

// Add adds a new public key to the authenticated user's account
func (d *SDeployKeys) Add(name, keyType, key, svcID string) error {
	deployKey := models.DeployKey{
		Key:  key,
		Name: name,
		Type: keyType,
	}
	b, err := json.Marshal(deployKey)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/ssh_keys", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
