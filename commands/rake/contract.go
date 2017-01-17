package rake

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "rake",
	ShortHelp: "Execute a rake task",
	LongHelp: "`rake` executes a rake task by its name asynchronously. " +
		"Once executed, the output of the task can be seen through your logging Dashboard. Here is a sample command\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" rake code-1 db:migrate\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The service that will run the rake task. Defaults to the associated service.")
			taskName := cmd.StringArg("TASK_NAME", "", "The name of the rake task to run")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRake(*serviceName, *taskName, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "SERVICE_NAME TASK_NAME"
		}
	},
}

// IRake
type IRake interface {
	Run(taskName, svcID string) error
}

// SRake is a concrete implementation of IRake
type SRake struct {
	Settings *models.Settings
}

// New returns an instance of IRake
func New(settings *models.Settings) IRake {
	return &SRake{
		Settings: settings,
	}
}
