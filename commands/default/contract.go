package defaultcmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "default",
	ShortHelp: "Set the default associated environment",
	LongHelp:  "Set the default associated environment",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			alias := cmd.StringArg("ENV_ALIAS", "", "The alias of an already associated environment to set as the default")
			cmd.Action = func() {
				if err := config.CheckRequiredAssociation(true, false, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdDefault(*alias, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "ENV_ALIAS"
		}
	},
}

// IDefault
type IDefault interface {
	Set(alias string) error
}

// SDefault is a concrete implementation of IDefault
type SDefault struct {
	Settings *models.Settings
}

// New returns an instance of IDefault
func New(settings *models.Settings) IDefault {
	return &SDefault{
		Settings: settings,
	}
}
