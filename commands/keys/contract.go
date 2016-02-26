package keys

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/deploykeys"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd for keys
var Cmd = models.Command{
	Name:      "keys",
	ShortHelp: "Tasks for SSH keys",
	LongHelp:  "Tasks for your own SSH keys",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(AddSubCmd.Name, AddSubCmd.ShortHelp, AddSubCmd.CmdFunc(settings))
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(RemoveSubCmd.Name, RemoveSubCmd.ShortHelp, RemoveSubCmd.CmdFunc(settings))
			cmd.Command(SetSubCmd.Name, SetSubCmd.ShortHelp, SetSubCmd.CmdFunc(settings))
		}
	},
}

var AddSubCmd = models.Command{
	Name:      "add",
	ShortHelp: "Add a public key",
	LongHelp:  "Add a new RSA public key to your own user account",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			name := cmd.StringArg("NAME", "", "The name for the new key, for your own purposes")
			path := cmd.StringArg("PUBLIC_KEY_PATH", "", "Relative path to the public key file")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdAdd(*name, *path, New(settings))
				if err != nil {
					logrus.Fatal(err)
				}
			}
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List your public keys",
	LongHelp:  "List the names of all public keys currently attached to your user",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdList(New(settings), deploykeys.New(settings))
				if err != nil {
					logrus.Fatal(err)
				}
			}
		}
	},
}

var RemoveSubCmd = models.Command{
	Name:      "rm",
	ShortHelp: "Remove a public key",
	LongHelp:  "Remove a public key from your own account, by name.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			name := cmd.StringArg("NAME", "", "The name of the key to remove.")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRemove(*name, New(settings))
				if err != nil {
					logrus.Fatal(err)
				}
			}
		}
	},
}

var SetSubCmd = models.Command{
	Name:      "set",
	ShortHelp: "Set your auth key",
	LongHelp:  "Set the private key used to sign in instead of username and password. This is expected to correspond to an OpenSSH-formatted RSA public key in the same directory.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			path := cmd.StringArg("PRIVATE_KEY_PATH", "", "Relative path to the private key file")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdSet(*path, settings)
				if err != nil {
					logrus.Fatal(err)
				}
			}
		}
	},
}

type IKeys interface {
	List() (*[]models.UserKey, error)
	Add(name, publicKey string) error
	Remove(name string) error
}

type SKeys struct {
	Settings *models.Settings
}

func New(settings *models.Settings) IKeys {
	return &SKeys{Settings: settings}
}
