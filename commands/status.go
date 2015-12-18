package commands

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/catalyzeio/catalyze/helpers"
	"github.com/catalyzeio/catalyze/models"
	"github.com/pmylund/sortutil"
)

// Status prints out all of the non-utility services and their running jobs
func Status(settings *models.Settings) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, '\t', 0)

	helpers.SignIn(settings)
	env := helpers.RetrieveEnvironment("pod", settings)

	fmt.Fprintln(w, env.Name+" (environment ID = "+env.ID+"):")
	fmt.Fprintln(w, "Label\tStatus\tCreated At")

	services := *env.Services
	sortutil.AscByField(services, "Label")

	for _, service := range services {
		if service.Type != "utility" && service.Type != "" {
			jobs := helpers.RetrieveRunningJobs(service.ID, settings)
			for jobID, job := range *jobs {
				displayType := service.Label
				if job.Type != "deploy" {
					displayType = fmt.Sprintf("%s (%s)", service.Label, job.Type)
					if job.Type == "worker" {
						// fetch the worker separately to get the procfile target run
						workerJob := helpers.RetrieveJob(jobID, service.ID, settings)
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
				latestBuildMap := helpers.RetrieveLatestBuildJob(service.ID, settings)
				for latestBuildID, latestBuild := range *latestBuildMap {
					if latestBuildID == "" {
						fmt.Fprintln(w, "--------"+"\t"+service.Label+"\t"+"-------"+"\t"+"---------------")
					} else if latestBuildID != "" {
						const dateForm = "2006-01-02T15:04:05"
						t, _ := time.Parse(dateForm, latestBuild.CreatedAt)
						displayType := service.Label
						displayType = fmt.Sprintf("%s (%s)", displayType, latestBuild.Type)
						fmt.Fprintln(w, displayType+"\t"+latestBuild.Status+"\t"+t.Local().Format(time.Stamp))
					}
				}
			}
		}
	}
	w.Flush()
}
