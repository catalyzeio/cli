package keys

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/deploykeys"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/models"
)

func CmdAdd(name, path string, ik IKeys, id deploykeys.IDeployKeys) error {
	if strings.ContainsAny(name, config.InvalidChars) {
		return fmt.Errorf("Invalid key name. Names must not contain the following characters: %s", config.InvalidChars)
	}
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	keyBytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	k, err := id.ParsePublicKey(keyBytes)
	if err != nil {
		return err
	}
	key := ssh.MarshalAuthorizedKey(k)
	if err != nil {
		return err
	}
	err = ik.Add(name, string(key))
	if err != nil {
		return err
	}
	logrus.Printf("Key '%s' added to your account.", name)
	logrus.Println("If you use an ssh-agent, make sure you add this key to your ssh-agent in order to push code")
	return nil
}

// Add adds a new public key to the authenticated user's account
func (k *SKeys) Add(name, publicKey string) error {
	body, err := json.Marshal(models.UserKey{
		Key:  publicKey,
		Name: name,
	})
	if err != nil {
		return err
	}
	headers := k.Settings.HTTPManager.GetHeaders(k.Settings.SessionToken, k.Settings.Version, k.Settings.Pod, k.Settings.UsersID)
	resp, status, err := k.Settings.HTTPManager.Post(body, fmt.Sprintf("%s%s/keys", k.Settings.AuthHost, k.Settings.AuthHostVersion), headers)
	if err != nil {
		return err
	}
	return k.Settings.HTTPManager.ConvertResp(resp, status, nil)
}
