package users

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/invites"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "users",
	ShortHelp: "Manage users who have access to the given organization",
	LongHelp:  "Manage users who have access to the given organization",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(RmSubCmd.Name, RmSubCmd.ShortHelp, RmSubCmd.CmdFunc(settings))
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all users who have access to the given organization",
	LongHelp:  "List all users who have access to the given organization",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdList(settings.UsersID, New(settings), invites.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

var RmSubCmd = models.Command{
	Name:      "rm",
	ShortHelp: "Revoke access to the given organization for the given user",
	LongHelp:  "Revoke access to the given organization for the given user",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			email := subCmd.StringArg("EMAIL", "", "The email address of the user to revoke access from for the given organization")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRm(*email, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "EMAIL"
		}
	},
}

// IUsers
type IUsers interface {
	List() (*[]models.OrgUser, error)
	Rm(usersID string) error
}

// SUsers is a concrete implementation of IUsers
type SUsers struct {
	Settings *models.Settings
}

// New generates a new instance of IUsers
func New(settings *models.Settings) IUsers {
	return &SUsers{
		Settings: settings,
	}
}
