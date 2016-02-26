package redeploy

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "redeploy",
	ShortHelp: "Redeploy a service without having to do a git push",
	LongHelp:  "Redeploy a service without having to do a git push",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to redeploy (i.e. 'app01')")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRedeploy(*serviceName, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "SERVICE_NAME"
		}
	},
}

// IRedeploy
type IRedeploy interface {
	Redeploy(service *models.Service) error
}

// SRedeploy is a concrete implementation of IRedeploy
type SRedeploy struct {
	Settings *models.Settings
}

// New returns an instance of IRedeploy
func New(settings *models.Settings) IRedeploy {
	return &SRedeploy{
		Settings: settings,
	}
}
