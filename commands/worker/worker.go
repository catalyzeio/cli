package worker

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/httpclient"
)

func CmdWorker(target, svcID string, iw IWorker) error {
	logrus.Printf("Initiating a background worker for Service: %s (procfile target = \"%s\")", svcID, target)
	err := iw.Start(target)
	if err != nil {
		return err
	}
	logrus.Println("Worker started.")
	return nil
}

// Start starts a Procfile target as a worker. Worker containers are intended
// to be short-lived, one-off tasks.
func (w *SWorker) Start(target string) error {
	worker := map[string]string{
		"target": target,
	}
	b, err := json.Marshal(worker)
	if err != nil {
		return err
	}
	headers := httpclient.GetHeaders(w.Settings.SessionToken, w.Settings.Version, w.Settings.Pod)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/worker", w.Settings.PaasHost, w.Settings.PaasHostVersion, w.Settings.EnvironmentID, w.Settings.ServiceID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
