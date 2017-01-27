package environments

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "environments",
	ShortHelp: "Manage environments for which you have access",
	LongHelp:  "The `environments` command allows you to manage your environments. The environments command can not be run directly but has sub commands.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.CommandLong(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.LongHelp, ListSubCmd.CmdFunc(settings))
			cmd.CommandLong(RenameSubCmd.Name, RenameSubCmd.ShortHelp, RenameSubCmd.LongHelp, RenameSubCmd.CmdFunc(settings))
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all environments you have access to",
	LongHelp: "`environments list` lists all environments that you are granted access to. " +
		"These environments include those you created and those that other Catalyze customers have added you to. " +
		"Here is a sample command\n\n" +
		"```\ncatalyze environments list\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatalln(err.Error())
				}
				err := CmdList(New(settings))
				if err != nil {
					logrus.Fatalln(err.Error())
				}
			}
		}
	},
}

var RenameSubCmd = models.Command{
	Name:      "rename",
	ShortHelp: "Rename an environment",
	LongHelp: "`environments rename` allows you to rename your environment. Here is a sample command\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" environments rename MyNewEnvName\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			name := subCmd.StringArg("NAME", "", "The new name of the environment")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRename(settings.EnvironmentID, *name, New(settings))
				if err != nil {
					logrus.Fatalln(err.Error())
				}
			}
			subCmd.Spec = "NAME"
		}
	},
}

// IEnvironments is an interface for interacting with environments
type IEnvironments interface {
	List() (*[]models.Environment, map[string]error)
	Retrieve(envID string) (*models.Environment, error)
	Update(envID string, updates map[string]string) error
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
