package status

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "status",
	ShortHelp: "Get quick readout of the current status of your associated environment and all of its services",
	LongHelp:  "Get quick readout of the current status of your associated environment and all of its services",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdStatus(settings.EnvironmentID, New(settings, jobs.New(settings)), environments.New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

// IStatus
type IStatus interface {
	Status(env *models.Environment, services *[]models.Service) error
}

// SStatus is a concrete implementation of IStatus
type SStatus struct {
	Settings *models.Settings
	Jobs     jobs.IJobs
}

// New returns an instance of IStatus
func New(settings *models.Settings, ij jobs.IJobs) IStatus {
	return &SStatus{
		Settings: settings,
		Jobs:     ij,
	}
}
