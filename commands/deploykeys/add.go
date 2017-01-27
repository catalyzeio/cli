package deploykeys

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/models"
)

func CmdAdd(name, keyPath, svcName string, id IDeployKeys, is services.IServices) error {
	if strings.ContainsAny(name, config.InvalidChars) {
		return fmt.Errorf("Invalid SSH key name. Names must not contain the following characters: %s", config.InvalidChars)
	}
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return fmt.Errorf("A file does not exist at path '%s'", keyPath)
	}
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	if service.Type != "code" {
		return fmt.Errorf("You can only add deploy keys to code services, not %s services", service.Type)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return err
	}
	k, err := id.ParsePublicKey(key)
	if err != nil {
		return err
	}
	key = ssh.MarshalAuthorizedKey(k)
	return id.Add(name, "ssh", string(key), service.ID)
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
	headers := d.Settings.HTTPManager.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod, d.Settings.UsersID)
	resp, statusCode, err := d.Settings.HTTPManager.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/ssh_keys", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return err
	}
	return d.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
