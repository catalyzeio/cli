package invites

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "invites",
	ShortHelp: "Manage invitations for your organizations",
	LongHelp:  "Manage invitations for your organizations",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(AcceptSubCmd.Name, AcceptSubCmd.ShortHelp, AcceptSubCmd.CmdFunc(settings))
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(RmSubCmd.Name, RmSubCmd.ShortHelp, RmSubCmd.CmdFunc(settings))
			cmd.Command(SendSubCmd.Name, SendSubCmd.ShortHelp, SendSubCmd.CmdFunc(settings))
		}
	},
}

var AcceptSubCmd = models.Command{
	Name:      "accept",
	ShortHelp: "Accept an organization invite",
	LongHelp:  "Accept an organization invite",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			inviteCode := subCmd.StringArg("INVITE_CODE", "", "The invite code that was sent in the invite email")
			subCmd.Action = func() {
				p := prompts.New()
				a := auth.New(settings, p)
				if _, err := a.Signin(); err != nil {
					logrus.Fatal(err.Error())
				}

				err := CmdAccept(*inviteCode, New(settings), a, p)
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "INVITE_CODE"
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all pending organization invitations",
	LongHelp:  "List all pending organization invitations",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdList(settings.EnvironmentName, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

var RmSubCmd = models.Command{
	Name:      "rm",
	ShortHelp: "Remove a pending organization invitation",
	LongHelp:  "Remove a pending organization invitation",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			inviteID := subCmd.StringArg("INVITE_ID", "", "The ID of an invitation to remove")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRm(*inviteID, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "INVITE_ID"
		}
	},
}

var SendSubCmd = models.Command{
	Name:      "send",
	ShortHelp: "Send an invite to a user by email for a given organization",
	LongHelp:  "Send an invite to a user by email for a given organization",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			email := subCmd.StringArg("EMAIL", "", "The email of a user to invite to the associated environment. This user does not need to have a Catalyze account prior to sending the invitation")
			subCmd.BoolOpt("m member", true, "Whether or not the user will be invited as a basic member")
			adminRole := subCmd.BoolOpt("a admin", false, "Whether or not the user will be invited as an admin")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				role := "member"
				if *adminRole {
					role = "admin"
				}
				err := CmdSend(*email, role, settings.EnvironmentName, New(settings), prompts.New())
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "EMAIL [-m | -a]"
		}
	},
}

// IInvites
type IInvites interface {
	Accept(inviteCode string) (string, error)
	List() (*[]models.Invite, error)
	ListRoles() (*[]models.Role, error)
	Rm(inviteID string) error
	Send(email string, role int) error
}

// SInvites is a concrete implementation of IInvites
type SInvites struct {
	Settings *models.Settings
}

// New generates a new instance of IInvites
func New(settings *models.Settings) IInvites {
	return &SInvites{
		Settings: settings,
	}
}
