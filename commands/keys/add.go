package keys

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/mitchellh/go-homedir"
)

func CmdAdd(name, path string, ik IKeys) error {
	fullPath, err := homedir.Expand(path)
	if err != nil {
		return err
	}
	keyBytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	err = ik.Add(name, string(keyBytes))
	if err != nil {
		return err
	}
	logrus.Printf("Key '%s' added to your account.", name)
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
	headers := httpclient.GetHeaders(k.Settings.SessionToken, k.Settings.Version, k.Settings.Pod)
	resp, status, err := httpclient.Post(body, fmt.Sprintf("%s%s/keys", k.Settings.AuthHost, k.Settings.AuthHostVersion), headers)
	if err != nil {
		return err
	}
	if httpclient.IsError(status) {
		return httpclient.ConvertError(resp, status)
	}
	return nil
}
