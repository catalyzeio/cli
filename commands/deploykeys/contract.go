package deploykeys

import (
	"crypto/rsa"

	"golang.org/x/crypto/ssh"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

var Cmd = models.Command{
	Name:      "deploy-keys",
	ShortHelp: "Tasks for SSH deploy keys",
	LongHelp:  "Tasks for SSH deploy keys",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(AddSubCmd.Name, AddSubCmd.ShortHelp, AddSubCmd.CmdFunc(settings))
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(RmSubCmd.Name, RmSubCmd.ShortHelp, RmSubCmd.CmdFunc(settings))
		}
	},
}

var AddSubCmd = models.Command{
	Name:      "add",
	ShortHelp: "Add a new deploy key",
	LongHelp:  "Add a new deploy key to a code service on the associated environment",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			name := subCmd.StringArg("NAME", "", "The name for the new key, for your own purposes")
			path := subCmd.StringArg("KEY_PATH", "", "Relative path to the SSH key file")
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the code service to add this deploy key to")
			private := subCmd.BoolOpt("p private", false, "Whether or not this is a private key")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdAdd(*name, *path, *serviceName, *private, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "NAME KEY_PATH SERVICE_NAME [-p]"
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all deploy keys",
	LongHelp:  "List all deploy keys for a code service on the associated environment",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the code service to list deploy keys")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdList(*serviceName, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "SERVICE_NAME"
		}
	},
}

var RmSubCmd = models.Command{
	Name:      "rm",
	ShortHelp: "Remove a deploy key",
	LongHelp:  "Remove a deploy key from a code service on the associated environment",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			name := subCmd.StringArg("NAME", "", "The name of the key to remove")
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the code service to remove this deploy key from")
			private := subCmd.BoolOpt("p private", false, "Whether or not this is a private key")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRm(*name, *serviceName, *private, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "NAME SERVICE_NAME [-p]"
		}
	},
}

type IDeployKeys interface {
	Add(name, keyType, key, svcID string) error
	ExtractPublicKey(privKey *rsa.PrivateKey) (ssh.PublicKey, error)
	List(svcID string) (*[]models.DeployKey, error)
	ParsePrivateKey(b []byte) (*rsa.PrivateKey, error)
	ParsePublicKey(b []byte) (ssh.PublicKey, error)
	Rm(name, keyType, svcID string) error
}

type SDeployKeys struct {
	Settings *models.Settings
}

func New(settings *models.Settings) IDeployKeys {
	return &SDeployKeys{
		Settings: settings,
	}
}
