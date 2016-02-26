package environments

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "environments",
	ShortHelp: "List all environments you have access to",
	LongHelp:  "List all environments you have access to",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdEnvironments(New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

// IEnvironments is an interface for interacting with environments
type IEnvironments interface {
	List() (*[]models.Environment, error)
	Retrieve(envID string) (*models.Environment, error)
}

// SEnvironments is a concrete implementation of IEnvironments
type SEnvironments struct {
	Settings *models.Settings
}

// New generates a new instance of IEnvironments
func New(settings *models.Settings) IEnvironments {
	return &SEnvironments{
		Settings: settings,
	}
}
