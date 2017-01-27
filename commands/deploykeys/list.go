package deploykeys

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/models"
	"github.com/olekukonko/tablewriter"
)

func CmdList(svcName string, id IDeployKeys, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	if service.Type != "code" {
		return fmt.Errorf("You can only list deploy keys for code services, not %s services", service.Type)
	}

	keys, err := id.List(service.ID)
	if err != nil {
		return err
	}
	if keys == nil || len(*keys) == 0 {
		logrus.Println("No deploy-keys found")
		return nil
	}

	invalidKeys := map[string]string{}

	data := [][]string{{"NAME", "TYPE", "FINGERPRINT"}}
	for _, key := range *keys {
		if key.Type != "ssh" {
			continue
		}
		s, err := id.ParsePublicKey([]byte(key.Key))
		if err != nil {
			invalidKeys[key.Name] = err.Error()
			continue
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
	headers := d.Settings.HTTPManager.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod, d.Settings.UsersID)
	resp, statusCode, err := d.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/ssh_keys", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return nil, err
	}
	var deployKeys []models.DeployKey
	err = d.Settings.HTTPManager.ConvertResp(resp, statusCode, &deployKeys)
	if err != nil {
		return nil, err
	}
	return &deployKeys, nil
}
