package associated

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "associated",
	ShortHelp: "Lists all associated environments",
	LongHelp: "`associated` outputs information about all previously associated environments on your local machine. " +
		"The information that is printed out includes the alias, environment ID, actual environment name, service ID, and the git repo directory. Here is a sample command\n\n" +
		"```\ncatalyze associated\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				err := CmdAssociated(New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

// IAssociated
type IAssociated interface {
	Associated() (map[string]models.AssociatedEnvV2, error)
}

// SAssociated is a concrete implementation of IAssociated
type SAssociated struct {
	Settings *models.Settings
}

// New returns an instance of IAssociated
func New(settings *models.Settings) IAssociated {
	return &SAssociated{
		Settings: settings,
	}
}
