package keys

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/httpclient"
)

func CmdRemove(name string, ik IKeys) error {
	err := ik.Remove(name)
	if err != nil {
		return err
	}
	logrus.Printf("Key '%s' has been removed from your account.", name)
	return nil
}

func (k *SKeys) Remove(name string) error {
	headers := httpclient.GetHeaders(k.Settings.SessionToken, k.Settings.Version, k.Settings.Pod)
	resp, status, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/keys/%s", k.Settings.AuthHost, k.Settings.AuthHostVersion, name), headers)
	if err != nil {
		return err
	}
	if httpclient.IsError(status) {
		return httpclient.ConvertError(resp, status)
	}
	return nil
}
