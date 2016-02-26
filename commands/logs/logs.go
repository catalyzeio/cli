package logs

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/commands/sites"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
)

const size = 50

// CmdLogs is a way to stream logs from Kibana to your local terminal. This is
// useful because Kibana is hard to look at because it splits every single
// log statement into a separate block that spans multiple lines so it's
// not very cohesive. This is intended to be similar to the `heroku logs`
// command.
func CmdLogs(queryString string, follow bool, hours, minutes, seconds int, envID string, il ILogs, ip prompts.IPrompts, ie environments.IEnvironments, is services.IServices, isites sites.ISites) error {
	username := os.Getenv(config.CatalyzeUsernameEnvVar)
	password := os.Getenv(config.CatalyzePasswordEnvVar)
	if username == "" || password == "" {
		logrus.Println("Your dashboard credentials are required to fetch logs")
		u, p, err := ip.UsernamePassword()
		if err != nil {
			return err
		}
		username = u
		password = p
	}
	env, err := ie.Retrieve(envID)
	if err != nil {
		return err
	}
	serviceProxy, err := is.RetrieveByLabel("service_proxy")
	if err != nil {
		return err
	}
	sites, err := isites.List(serviceProxy.ID)
	if err != nil {
		return err
	}
	domain := ""
	for _, site := range *sites {
		if strings.HasPrefix(site.Name, env.Namespace) {
			domain = site.Name
		}
	}
	if domain == "" {
		return errors.New("Could not determine the fully qualified domain name of your environment")
	}
	from := 0
	offset := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	timestamp := time.Now().In(time.UTC).Add(-1 * offset)
	from, timestamp, err = il.Output(queryString, username, password, domain, follow, hours, minutes, seconds, from, timestamp, time.Now(), env)
	if err != nil {
		return err
	}
	if follow {
		return il.Stream(queryString, username, password, domain, follow, hours, minutes, seconds, from, timestamp, env)
	}
	return nil
}

func (l *SLogs) Output(queryString, username, password, domain string, follow bool, hours, minutes, seconds, from int, startTimestamp, endTimestamp time.Time, env *models.Environment) (int, time.Time, error) {
	appLogsIdentifier := "source"
	appLogsValue := "app"
	if strings.HasPrefix(domain, "pod01") || strings.HasPrefix(domain, "csb01") {
		appLogsIdentifier = "syslog_program"
		appLogsValue = "supervisord"
	}

	urlString := fmt.Sprintf("https://%s/__es", domain)

	basicAuth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	headers := map[string][]string{"Authorization": {"Basic " + basicAuth}}

	logrus.Println("        @timestamp       -        message")
	for {
		queryBytes := generateQuery(queryString, appLogsIdentifier, appLogsValue, startTimestamp, from)

		resp, statusCode, err := httpclient.Get(queryBytes, fmt.Sprintf("%s/_search", urlString), headers)
		if err != nil {
			return from, startTimestamp, err
		}
		var logs models.Logs
		err = httpclient.ConvertResp(resp, statusCode, &logs)
		if err != nil {
			return from, startTimestamp, err
		}

		end := time.Time{}
		for _, lh := range *logs.Hits.Hits {
			logrus.Printf("%s - %s", lh.Fields["@timestamp"][0], lh.Fields["message"][0])
			end, _ = time.Parse(time.RFC3339Nano, lh.Fields["@timestamp"][0])
		}
		amount := len(*logs.Hits.Hits)

		from += len(*logs.Hits.Hits)
		// TODO this will infinite loop if it always retrieves `size` hits
		// and it fails to parse the end timestamp. very small window of opportunity.
		if amount < size || end.After(endTimestamp) {
			break
		}
		time.Sleep(config.JobPollTime * time.Second)
	}
	return from, startTimestamp, nil
}

func (l *SLogs) Stream(queryString, username, password, domain string, follow bool, hours, minutes, seconds, from int, timestamp time.Time, env *models.Environment) error {
	for {
		f, t, err := l.Output(queryString, username, password, domain, follow, hours, minutes, seconds, from, timestamp, time.Now(), env)
		if err != nil {
			return err
		}
		from = f
		timestamp = t
		time.Sleep(config.LogPollTime * time.Second)
	}
}

func generateQuery(queryString, appLogsIdentifier, appLogsValue string, timestamp time.Time, from int) []byte {
	query := `{
	"fields": ["@timestamp", "message", "` + appLogsIdentifier + `"],
	"query": {
		"wildcard": {
			"message": "` + queryString + `"
		}
	},
	"filter": {
		"bool": {
			"must": [
				{"term": {"` + appLogsIdentifier + `": "` + appLogsValue + `"}},
				{"range": {"@timestamp": {"gt": "` + fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", timestamp.Year(), timestamp.Month(), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second()) + `"}}}
			]
		}
	},
	"sort": {
		"@timestamp": {
			"order": "asc"
		},
		"message": {
			"order": "asc"
		}
	},
	"from": ` + fmt.Sprintf("%d", from) + `,
	"size": ` + fmt.Sprintf("%d", size) + `
	}`
	var buf bytes.Buffer
	json.Compact(&buf, []byte(query))
	return buf.Bytes()
}
