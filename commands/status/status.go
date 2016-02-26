package status

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/models"
	"github.com/pmylund/sortutil"
)

func CmdStatus(envID string, is IStatus, ie environments.IEnvironments, iservices services.IServices) error {
	env, err := ie.Retrieve(envID)
	if err != nil {
		return err
	}
	svcs, err := iservices.ListByEnvID(env.ID, env.Pod)
	if err != nil {
		return err
	}
	return is.Status(env, svcs)
}

// Status prints out all of the non-utility services and their running jobs
func (s *SStatus) Status(env *models.Environment, services *[]models.Service) error {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, '\t', 0)

	fmt.Fprintln(w, env.Name+" (environment ID = "+env.ID+"):")
	fmt.Fprintln(w, "Label\tStatus\tCreated At")

	sortutil.AscByField(*services, "Label")

	for _, service := range *services {
		if service.Type != "utility" && service.Type != "" {
			s.Settings.ServiceID = service.ID
			jobs, err := s.Jobs.RetrieveByStatus("running")
			if err != nil {
				return err
			}
			for _, job := range *jobs {
				displayType := service.Label
				if job.Type != "deploy" {
					displayType = fmt.Sprintf("%s (%s)", service.Label, job.Type)
					if job.Type == "worker" {
						// fetch the worker separately to get the procfile target run
						workerJob, err := s.Jobs.Retrieve(job.ID, service.ID, true)
						if err != nil {
							return err
						}
						if workerJob.Spec != nil && workerJob.Spec.Payload != nil && workerJob.Spec.Payload.Environment != nil {
							if target, contains := workerJob.Spec.Payload.Environment["PROCFILE_TARGET"]; contains {
								displayType = fmt.Sprintf("%s (%s: target=%s)", service.Label, job.Type, target)
							}
						}
					}
				}

				const dateForm = "2006-01-02T15:04:05"
				t, _ := time.Parse(dateForm, job.CreatedAt)
				fmt.Fprintln(w, displayType+"\t"+job.Status+"\t"+t.Local().Format(time.Stamp))
			}
			if service.Type == "code" {
				latestBuildJobs, err := s.Jobs.RetrieveByType("build", 1, 1)
				if err != nil {
					return err
				}
				for _, latestBuildJob := range *latestBuildJobs {
					if latestBuildJob.ID == "" {
						fmt.Fprintln(w, "--------"+"\t"+service.Label+"\t"+"-------"+"\t"+"---------------")
					} else if latestBuildJob.ID != "" {
						const dateForm = "2006-01-02T15:04:05"
						t, _ := time.Parse(dateForm, latestBuildJob.CreatedAt)
						displayType := service.Label
						displayType = fmt.Sprintf("%s (%s)", displayType, latestBuildJob.Type)
						fmt.Fprintln(w, displayType+"\t"+latestBuildJob.Status+"\t"+t.Local().Format(time.Stamp))
					}
				}
			}
		}
	}
	w.Flush()
	return nil
}
