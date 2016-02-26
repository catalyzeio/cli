package associate

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/git"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "associate",
	ShortHelp: "Associates an environment",
	LongHelp:  "Associates an environment and service using the given alias to your local machine. For all further commands, the alias will be used instead. If no alias is given, the environment name is used.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			envName := cmd.StringArg("ENV_NAME", "", "The name of your environment")
			serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the primary code service to associate with this environment (i.e. 'app01')")
			alias := cmd.StringOpt("a alias", "", "A shorter name to reference your environment by for local commands")
			remote := cmd.StringOpt("r remote", "catalyze", "The name of the remote")
			defaultEnv := cmd.BoolOpt("d default", false, "Specifies whether or not the associated environment will be the default")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdAssociate(*envName, *serviceName, *alias, *remote, *defaultEnv, New(settings), git.New(), environments.New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "ENV_NAME SERVICE_NAME [-a] [-r] [-d]"
		}
	},
}

// interfaces are the API calls
type IAssociate interface {
	Associate(name, remote string, defaultEnv bool, env *models.Environment, chosenService *models.Service) error
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
