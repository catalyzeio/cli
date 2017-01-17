package supportids

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "support-ids",
	ShortHelp: "Print out various IDs related to your associated environment to be used when contacting Catalyze support",
	LongHelp: "`support-ids` is helpful when contacting Catalyze support by sending an email to support@catalyze.io. " +
		"If you are having an issue with a CLI command or anything with your environment, it is helpful to run this command and copy the output into the initial correspondence with a Catalyze engineer. " +
		"This will help Catalyze identify the environment faster and help come to resolution faster. Here is a sample command\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" support-ids\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				err := CmdSupportIDs(New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

// ISupportIDs
type ISupportIDs interface {
	SupportIDs() (string, string, string, string, error)
}

// SSupportIDs is a concrete implementation of ISupportIDs
type SSupportIDs struct {
	Settings *models.Settings
}

// New returns an instance of ISupportIDs
func New(settings *models.Settings) ISupportIDs {
	return &SSupportIDs{
		Settings: settings,
	}
}
