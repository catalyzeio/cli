package worker

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/lib/prompts"
)

func CmdRm(svcName, target string, iw IWorker, is services.IServices, ip prompts.IPrompts, ij jobs.IJobs) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	err = ip.YesNo(fmt.Sprintf("Removing the worker target %s for service %s will automatically stop all existing worker jobs with that target.", target, svcName), "Would you like to proceed? (y/n) ")
	if err != nil {
		return err
	}
	jobs, err := ij.RetrieveByTarget(service.ID, target, 1, 1000)
	if err != nil {
		return err
	}
	for _, j := range *jobs {
		err = ij.Delete(j.ID, service.ID)
		if err != nil {
			return err
		}
	}
	workers, err := iw.Retrieve(service.ID)
	if err != nil {
		return err
	}
	delete(workers.Workers, target)
	err = iw.Update(service.ID, workers)
	if err != nil {
		return err
	}
	logrus.Printf("Successfully removed all workers with target %s for service %s", target, svcName)
	return nil
}
