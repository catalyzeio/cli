package console

import (
	"github.com/Sirupsen/logrus"
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
	Name:      "console",
	ShortHelp: "Open a secure console to a service",
	LongHelp:  "Open a secure console to a service",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to open up a console for")
			command := cmd.StringArg("COMMAND", "", "An optional command to run when the console becomes available")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdConsole(*serviceName, *command, New(settings, jobs.New(settings)), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "SERVICE_NAME [COMMAND]"
		}
	},
}

// IConsole
type IConsole interface {
	Open(command string, service *models.Service) error
	Request(command string, service *models.Service) (*models.Job, error)
	RetrieveTokens(jobID string, service *models.Service) (*models.ConsoleCredentials, error)
	Destroy(jobID string, service *models.Service) error
}

// SConsole is a concrete implementation of IConsole
type SConsole struct {
	Settings *models.Settings
	Jobs     jobs.IJobs
}

// New returns an instance of IConsole
func New(settings *models.Settings, jobs jobs.IJobs) IConsole {
	return &SConsole{
		Settings: settings,
		Jobs:     jobs,
	}
}
