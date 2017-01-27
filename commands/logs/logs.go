package logs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/commands/sites"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
)

const size = 50

// CmdLogs is a way to stream logs from Kibana to your local terminal. This is
// useful because Kibana is hard to look at because it splits every single
// log statement into a separate block that spans multiple lines so it's
// not very cohesive. This is intended to be similar to the `heroku logs`
// command.
func CmdLogs(queryString string, follow bool, hours, minutes, seconds int, envID string, settings *models.Settings, il ILogs, ip prompts.IPrompts, ie environments.IEnvironments, is services.IServices, isites sites.ISites) error {
	if follow && (hours > 0 || minutes > 0 || seconds > 0) {
		return fmt.Errorf("Specifying \"-f\" in combination with \"--hours\", \"--minutes\", or \"--seconds\" is unsupported.")
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
			break
		}
	}
	if domain == "" {
		return errors.New("Could not determine the fully qualified domain name of your environment. Please contact Catalyze Support at support@catalyze.io with this error message to resolve this issue.")
	}
	if follow {
		if err = il.Watch(queryString, domain, settings.SessionToken); err != nil {
			logrus.Debugf("Error attempting to stream logs from logwatch: %s", err)
		} else {
			return nil
		}
	}
	from := 0
	offset := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	timestamp := time.Now().In(time.UTC).Add(-1 * offset)
	from, err = il.Output(queryString, settings.SessionToken, domain, from, timestamp, time.Now())
	if err != nil {
		return err
	}
	if follow {
		return il.Stream(queryString, settings.SessionToken, domain, from, timestamp)
	}
	return nil
}

func (l *SLogs) Output(queryString, sessionToken, domain string, from int, startTimestamp, endTimestamp time.Time) (int, error) {
	appLogsIdentifier := "source"
	appLogsValue := "app"
	if strings.HasPrefix(domain, "pod01") || strings.HasPrefix(domain, "csb01") {
		appLogsIdentifier = "syslog_program"
		appLogsValue = "supervisord"
	}

	urlString := fmt.Sprintf("https://%s/__es", domain)

	headers := map[string][]string{"Cookie": {"sessionToken=" + url.QueryEscape(sessionToken)}}

	logrus.Println("        @timestamp       -        message")
	for {
		queryBytes := generateQuery(queryString, appLogsIdentifier, appLogsValue, startTimestamp, from)

		resp, statusCode, err := l.Settings.HTTPManager.Get(queryBytes, fmt.Sprintf("%s/_search", urlString), headers)
		if err != nil {
			return from, err
		}
		var logs models.Logs
		err = l.Settings.HTTPManager.ConvertResp(resp, statusCode, &logs)
		if err != nil {
			return from, err
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
	return from, nil
}

func (l *SLogs) Stream(queryString, sessionToken, domain string, from int, timestamp time.Time) error {
	for {
		f, err := l.Output(queryString, sessionToken, domain, from, timestamp, time.Now())
		if err != nil {
			return err
		}
		from = f
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
