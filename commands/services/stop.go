package services

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/lib/prompts"
)

// CmdStop stops all instances of a given service. All workers and rake tasks will also be stopped
// if applicable.
func CmdStop(svcName string, is IServices, ij jobs.IJobs, ip prompts.IPrompts) error {
	err := ip.YesNo(fmt.Sprintf("Stopping %s will stop all instances of the service, all workers, all rake tasks, and all currently open consoles.", svcName), fmt.Sprintf("Are you sure you want to stop %s? (y/n) ", svcName))
	if err != nil {
		return err
	}
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	if !service.Redeployable {
		return fmt.Errorf("This service cannot be stopped. Please contact Catalyze Support at support@catalyze.io if you need the \"%s\" service stopped.", svcName)
	}

	page := 0
	pageSize := 100
	for {
		jobs, err := ij.List(service.ID, page, pageSize)
		if err != nil {
			return err
		}

		for _, job := range *jobs {
			if job.Status != "scheduled" && job.Status != "queued" && job.Status != "started" && job.Status != "running" && job.Status != "waiting" {
				logrus.Debugf("Skipping %s job (%s)", job.Status, job.ID)
				continue
			}
			logrus.Debugf("Deleting %s job (%s) on service %s", job.Type, job.ID, service.ID)
			err = ij.Delete(job.ID, service.ID)
			if err != nil {
				return err
			}
		}
		if len(*jobs) < pageSize {
			break
		}
		page++
	}

	logrus.Printf("Successfully stopped %s. Run \"catalyze redeploy %s\" to start this service again.", svcName, svcName)
	return nil
}
