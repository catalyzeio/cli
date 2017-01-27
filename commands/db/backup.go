package db

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/models"
)

func CmdBackup(databaseName string, skipPoll bool, id IDb, is services.IServices, ij jobs.IJobs) error {
	service, err := is.RetrieveByLabel(databaseName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", databaseName)
	}
	job, err := id.Backup(service)
	if err != nil {
		return err
	}
	logrus.Printf("Backup started (job ID = %s)", job.ID)
	isSnapshotBackup := job.IsSnapshotBackup != nil && *job.IsSnapshotBackup
	if !skipPoll {
		// all because logrus treats print, println, and printf the same
		logrus.Println("Polling until backup finishes.")
		if isSnapshotBackup {
			logrus.Printf("This is a snapshot backup, it may be a while before this backup shows up in the \"catalyze db list %s\" command.", databaseName)
			err = ij.WaitToAppear(job.ID, service.ID)
			if err != nil {
				return err
			}
		}
		status, err := ij.PollTillFinished(job.ID, service.ID)
		if err != nil {
			return err
		}
		job.Status = status
		logrus.Printf("Ended in status '%s'", job.Status)
		err = id.DumpLogs("backup", job, service)
		if err != nil {
			return err
		}
		if job.Status != "finished" {
			return fmt.Errorf("Job finished with invalid status %s", job.Status)
		}
	} else if isSnapshotBackup {
		logrus.Printf("This is a snapshot backup, it may be a while before this backup shows up in the \"catalyze db list %s\" command.", databaseName)
	}
	logrus.Printf("You can download your backup with the \"catalyze db download %s %s ./output_file_path\" command", databaseName, job.ID)
	return nil
}

// Backup creates a new backup
func (d *SDb) Backup(service *models.Service) (*models.Job, error) {
	headers := d.Settings.HTTPManager.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod, d.Settings.UsersID)
	resp, statusCode, err := d.Settings.HTTPManager.Post(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/backup", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, service.ID), headers)
	if err != nil {
		return nil, err
	}
	var job models.Job
	err = d.Settings.HTTPManager.ConvertResp(resp, statusCode, &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}
