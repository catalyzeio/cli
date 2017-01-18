package services

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "services",
	ShortHelp: "Perform operations on an environment's services",
	LongHelp:  "The `services` command allows you to manage your services. The services command cannot be run directly but has sub commands.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.CommandLong(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.LongHelp, ListSubCmd.CmdFunc(settings))
			cmd.CommandLong(StopSubCmd.Name, StopSubCmd.ShortHelp, StopSubCmd.LongHelp, StopSubCmd.CmdFunc(settings))
			cmd.CommandLong(RenameSubCmd.Name, RenameSubCmd.ShortHelp, RenameSubCmd.LongHelp, RenameSubCmd.CmdFunc(settings))
			cmd.Action = func() {
				logrus.Warnln("This command has been moved! Please use \"catalyze services list\" instead. This alias will be removed in the next CLI update.")
				logrus.Warnln("You can list all available services subcommands by running \"catalyze services --help\".")
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdServices(New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all services for your environment",
	LongHelp: "`services list` prints out a list of all services in your environment and their sizes. " +
		"The services will be printed regardless of their currently running state. " +
		"To see which services are currently running and which are not, use the [status](#status) command. " +
		"Here is a sample command\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" services list\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdServices(New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

var RenameSubCmd = models.Command{
	Name:      "rename",
	ShortHelp: "Rename a service",
	LongHelp: "`services rename` allows you to rename any service in your environment. Here is a sample command\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" services rename code-1 api-svc\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The service to rename")
			label := subCmd.StringArg("NEW_NAME", "", "The new name for the service")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRename(*serviceName, *label, New(settings))
				if err != nil {
					logrus.Fatalln(err.Error())
				}
			}
			subCmd.Spec = "SERVICE_NAME NEW_NAME"
		}
	},
}

var StopSubCmd = models.Command{
	Name:      "stop",
	ShortHelp: "Stop all instances of a given service (including all workers, rake tasks, and open consoles)",
	LongHelp: "`services stop` shuts down all running instances of a given service. " +
		"This is useful when performing maintenance and a service must be shutdown to perform that maintenance. " +
		"Take caution when running this command as all instances of the service, all workers, all rake tasks, and all open console sessions will be stopped. " +
		"Here is a sample command\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" services stop code-1\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			svcName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to stop")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdStop(*svcName, New(settings), jobs.New(settings), prompts.New())
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "SERVICE_NAME"
		}
	},
}

// IServices
type IServices interface {
	List() (*[]models.Service, error)
	ListByEnvID(envID, podID string) (*[]models.Service, error)
	RetrieveByLabel(label string) (*models.Service, error)
	Update(svcID string, updates map[string]string) error
}

// SServices is a concrete implementation of IServices
type SServices struct {
	Settings *models.Settings
}

// New generates a new instance of IServices
func New(settings *models.Settings) IServices {
	return &SServices{
		Settings: settings,
	}
}
