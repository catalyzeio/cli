package associate

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/commands/git"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "associate",
	ShortHelp: "Associates an environment",
	LongHelp: "`associate` is the entry point of the cli. You need to associate an environment before you can run most other commands. " +
		"Check out [scope](#global-scope) and [aliases](#environment-aliases) for more info on the value of the alias and default options. Here is a sample command\n\n" +
		"```\ncatalyze associate My-Production-Environment app01 -a prod\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			envName := cmd.StringArg("ENV_NAME", "", "The name of your environment")
			alias := cmd.StringOpt("a alias", "", "A shorter name to reference your environment by for local commands")
			remote := cmd.StringOpt("r remote", "catalyze", "The name of the remote")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdAssociate(*envName, *alias, *remote, New(settings), git.New(), environments.New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "ENV_NAME [-a] [-r]"
		}
	},
}

// interfaces are the API calls
type IAssociate interface {
	Associate(name, remote string, env *models.Environment) error
}

// SAssociate is a concrete implementation of IAssociate
type SAssociate struct {
	Settings *models.Settings
}

func New(settings *models.Settings) IAssociate {
	return &SAssociate{
		Settings: settings,
	}
}
