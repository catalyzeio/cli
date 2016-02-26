package dashboard

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "dashboard",
	ShortHelp: "Open the Catalyze Dashboard in your default browser",
	LongHelp:  "Open the Catalyze Dashboard in your default browser",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				err := CmdDashboard(New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

// IDashboard
type IDashboard interface {
	Open() error
}

// SDashboard is a concrete implementation of IDashboard
type SDashboard struct {
	Settings *models.Settings
}

// New returns an instance of IDashboard
func New(settings *models.Settings) IDashboard {
	return &SDashboard{
		Settings: settings,
	}
}
