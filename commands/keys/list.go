package keys

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/deploykeys"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/olekukonko/tablewriter"
)

func CmdList(ik IKeys, id deploykeys.IDeployKeys) error {
	keys, err := ik.List()
	if err != nil {
		return err
	}

	invalidKeys := map[string]string{}

	data := [][]string{{"NAME", "FINGERPRINT"}}
	for _, key := range *keys {
		s, err := id.ParsePublicKey([]byte(key.Key))
		if err != nil {
			invalidKeys[key.Name] = err.Error()
			continue
		}
		h := sha256.New()
		h.Write(s.Marshal())
		fingerprint := base64.StdEncoding.EncodeToString(h.Sum(nil))
		data = append(data, []string{key.Name, fmt.Sprintf("SHA256:%s", strings.TrimRight(fingerprint, "="))})
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

func (k *SKeys) List() (*[]models.UserKey, error) {
	headers := httpclient.GetHeaders(k.Settings.SessionToken, k.Settings.Version, k.Settings.Pod)
	resp, status, err := httpclient.Get(nil, fmt.Sprintf("%s%s/keys", k.Settings.AuthHost, k.Settings.AuthHostVersion), headers)
	if err != nil {
		return nil, err
	}

	keys := []models.UserKey{}
	err = httpclient.ConvertResp(resp, status, &keys)
	return &keys, err
}
