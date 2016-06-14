package update

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "update",
	ShortHelp: "Checks for available updates and updates the CLI if a new update is available",
	LongHelp:  "Checks for available updates and updates the CLI if a new update is available",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				err := CmdUpdate(New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
				updatePods(settings)
			}
		}
	},
}

// IUpdate
type IUpdate interface {
	Check() (bool, error)
	Update() error
}

// SUpdate is a concrete implementation of IUpdate
type SUpdate struct {
	Settings *models.Settings
}

// New returns an instance of IUpdate
func New(settings *models.Settings) IUpdate {
	return &SUpdate{
		Settings: settings,
	}
}
