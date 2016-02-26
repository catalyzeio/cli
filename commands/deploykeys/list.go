package deploykeys

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/olekukonko/tablewriter"
)

func CmdList(svcName string, id IDeployKeys, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"", svcName)
	}
	if service.Type != "code" {
		return fmt.Errorf("You can only list deploy keys to code services, not %s services", service.Type)
	}

	keys, err := id.List(service.ID)
	if err != nil {
		return err
	}

	invalidKeys := map[string]string{}

	data := [][]string{{"NAME", "TYPE", "FINGERPRINT"}}
	for _, key := range *keys {
		var s ssh.PublicKey
		if key.Type == "ssh_private" {
			privKey, err := id.ParsePrivateKey([]byte(key.Key))
			if err != nil {
				invalidKeys[key.Name] = err.Error()
				continue
			}
			s, err = id.ExtractPublicKey(privKey)
			if err != nil {
				invalidKeys[key.Name] = err.Error()
				continue
			}
		} else {
			s, err = id.ParsePublicKey([]byte(key.Key))
			if err != nil {
				invalidKeys[key.Name] = err.Error()
				continue
			}
		}
		h := sha256.New()
		h.Write(s.Marshal())
		fingerprint := base64.StdEncoding.EncodeToString(h.Sum(nil))
		data = append(data, []string{key.Name, key.Type, fmt.Sprintf("SHA256:%s", strings.TrimRight(fingerprint, "="))})
	}

	table := tablewriter.NewWriter(logrus.StandardLogger().Out)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.AppendBulk(data)
	table.Render()

	if len(invalidKeys) > 0 {
		logrus.Println("\nInvalid Keys:")
		for keyName, reason := range invalidKeys {
			logrus.Printf("%s: %s", keyName, reason)
		}
	}
	return nil
}

func (d *SDeployKeys) List(svcID string) (*[]models.DeployKey, error) {
	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/ssh_keys", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return nil, err
	}
	var deployKeys []models.DeployKey
	err = httpclient.ConvertResp(resp, statusCode, &deployKeys)
	if err != nil {
		return nil, err
	}
	return &deployKeys, nil
}
