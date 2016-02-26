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

func CmdCreate(givenHostname, pubKeyPath, privKeyPath string, selfSigned, resolve bool, ic ICerts, is services.IServices, issl ssl.ISSL) error {
	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("A cert does not exist at path '%s'", pubKeyPath)
	}
	if _, err := os.Stat(privKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("A private key does not exist at path '%s'", privKeyPath)
	}
	err := issl.Verify(pubKeyPath, privKeyPath, givenHostname, selfSigned)
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
	modifiedHostname := strings.Replace(givenHostname, "*", "star", -1)
	err = ic.Create(modifiedHostname, string(pubKeyBytes), string(privKeyBytes), service.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Created '%s'", modifiedHostname)
	return nil
}

func (c *SCerts) Create(hostname, pubKey, privKey, svcID string) error {
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
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/certs", c.Settings.PaasHost, c.Settings.PaasHostVersion, c.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
