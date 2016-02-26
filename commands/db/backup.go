package db

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/models"
)

func CmdBackup(databaseName string, skipPoll bool, id IDb, is services.IServices, ij jobs.IJobs) error {
	service, err := is.RetrieveByLabel(databaseName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"", databaseName)
	}
	job, err := id.Backup(service)
	if err != nil {
		return err
	}
	logrus.Printf("Backup started (job ID = %s)", job.ID)
	if !skipPoll {
		// all because logrus treats print, println, and printf the same
		logrus.StandardLogger().Out.Write([]byte("Polling until backup finishes."))
		status, err := ij.PollTillFinished(job.ID, service.ID)
		if err != nil {
			return err
		}
		job.Status = status
		logrus.Printf("\nEnded in status '%s'", job.Status)
		err = id.DumpLogs("backup", job, service)
		if err != nil {
			return err
		}
		if job.Status != "finished" {
			return fmt.Errorf("Job finished with invalid status %s", job.Status)
		}
	}
	return nil
}

// Backup creates a new backup
func (d *SDb) Backup(service *models.Service) (*models.Job, error) {
	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod)
	resp, statusCode, err := httpclient.Post(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/backup", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, service.ID), headers)
	if err != nil {
		return nil, err
	}
	var job models.Job
	err = httpclient.ConvertResp(resp, statusCode, &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}
