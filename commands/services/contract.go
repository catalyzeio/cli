package services

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "services",
	ShortHelp: "List all services for your environment",
	LongHelp:  "List all services for your environment",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdServices(New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

// IServices
type IServices interface {
	List() (*[]models.Service, error)
	ListByEnvID(envID, podID string) (*[]models.Service, error)
	Retrieve(svcID string) (*models.Service, error)
	RetrieveByLabel(label string) (*models.Service, error)
}

// SServices is a concrete implementation of IServices
type SServices struct {
	Settings *models.Settings
}

// New generates a new instance of IServices
func New(settings *models.Settings) IServices {
	return &SServices{
		Settings: settings,
	}
}
