package certs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/commands/ssl"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
)

func CmdUpdate(hostname, pubKeyPath, privKeyPath string, selfSigned, resolve bool, ic ICerts, is services.IServices, issl ssl.ISSL) error {
	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("A cert does not exist at path '%s'", pubKeyPath)
	}
	if _, err := os.Stat(privKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("A private key does not exist at path '%s'", privKeyPath)
	}
	err := issl.Verify(pubKeyPath, privKeyPath, hostname, selfSigned)
	var pubKeyBytes []byte
	var privKeyBytes []byte
	if err != nil {
		if ssl.IsIncompleteChainErr(err) && resolve {
			pubKeyBytes, err = issl.Resolve(pubKeyPath)
			if err != nil {
				return fmt.Errorf("Could not resolve the incomplete certificate chain. If this is a self signed certificate, please re-run this command with the '-s' option: %s", err.Error())
			}
		} else {
			return err
		}
	}
	service, err := is.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}
	if pubKeyBytes == nil {
		pubKeyBytes, err = ioutil.ReadFile(pubKeyPath)
		if err != nil {
			return err
		}
	}
	if privKeyBytes == nil {
		privKeyBytes, err = ioutil.ReadFile(privKeyPath)
		if err != nil {
			return err
		}
	}
	hostname = strings.Replace(hostname, "*", "star", -1)
	err = ic.Update(hostname, string(pubKeyBytes), string(privKeyBytes), service.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Updated '%s'", hostname)
	return nil
}

func (c *SCerts) Update(hostname, pubKey, privKey, svcID string) error {
	cert := models.Cert{
		Name:    hostname,
		PubKey:  pubKey,
		PrivKey: privKey,
	}
	b, err := json.Marshal(cert)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(c.Settings.SessionToken, c.Settings.Version, c.Settings.Pod)
	resp, statusCode, err := httpclient.Put(b, fmt.Sprintf("%s%s/environments/%s/services/%s/certs/%s", c.Settings.PaasHost, c.Settings.PaasHostVersion, c.Settings.EnvironmentID, svcID, hostname), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
