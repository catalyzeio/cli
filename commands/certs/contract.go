package certs

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/commands/ssl"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "certs",
	ShortHelp: "Manage your SSL certificates and domains",
	LongHelp:  "Manage your SSL certificates and domains",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(CreateSubCmd.Name, CreateSubCmd.ShortHelp, CreateSubCmd.CmdFunc(settings))
			cmd.Command(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.CmdFunc(settings))
			cmd.Command(RmSubCmd.Name, RmSubCmd.ShortHelp, RmSubCmd.CmdFunc(settings))
			cmd.Command(UpdateSubCmd.Name, UpdateSubCmd.ShortHelp, UpdateSubCmd.CmdFunc(settings))
		}
	},
}

var CreateSubCmd = models.Command{
	Name:      "create",
	ShortHelp: "Create a new domain with an SSL certificate and private key",
	LongHelp:  "Create a new domain with an SSL certificate and private key",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			name := subCmd.StringArg("HOSTNAME", "", "The hostname of this domain and SSL certificate plus private key pair")
			pubKeyPath := subCmd.StringArg("PUBLIC_KEY_PATH", "", "The path to a public key file in PEM format")
			privKeyPath := subCmd.StringArg("PRIVATE_KEY_PATH", "", "The path to an unencrypted private key file in PEM format")
			selfSigned := subCmd.BoolOpt("s self-signed", false, "Whether or not the given SSL certificate and private key are self signed")
			resolve := subCmd.BoolOpt("r resolve", true, "Whether or not to attempt to automatically resolve incomplete SSL certificate issues")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdCreate(*name, *pubKeyPath, *privKeyPath, *selfSigned, *resolve, New(settings), services.New(settings), ssl.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "HOSTNAME PUBLIC_KEY_PATH PRIVATE_KEY_PATH [-s] [-r]"
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List all existing domains that have SSL certificate and private key pairs",
	LongHelp:  "List all existing domains that have SSL certificate and private key pairs",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdList(New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
		}
	},
}

var RmSubCmd = models.Command{
	Name:      "rm",
	ShortHelp: "Remove an existing domain and its associated SSL certificate and private key pair",
	LongHelp:  "Remove an existing domain and its associated SSL certificate and private key pair",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			name := subCmd.StringArg("HOSTNAME", "", "The hostname of the domain and SSL certificate and private key pair")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdRm(*name, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "HOSTNAME"
		}
	},
}

var UpdateSubCmd = models.Command{
	Name:      "update",
	ShortHelp: "Update the SSL certificate and private key pair for an existing domain",
	LongHelp:  "Update the SSL certificate and private key pair for an existing domain",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			name := subCmd.StringArg("HOSTNAME", "", "The hostname of this domain and SSL certificate and private key pair")
			pubKeyPath := subCmd.StringArg("PUBLIC_KEY_PATH", "", "The path to a public key file in PEM format")
			privKeyPath := subCmd.StringArg("PRIVATE_KEY_PATH", "", "The path to an unencrypted private key file in PEM format")
			selfSigned := subCmd.BoolOpt("s self-signed", false, "Whether or not the given SSL certificate and private key are self signed")
			resolve := subCmd.BoolOpt("r resolve", true, "Whether or not to attempt to automatically resolve incomplete SSL certificate issues")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdUpdate(*name, *pubKeyPath, *privKeyPath, *selfSigned, *resolve, New(settings), services.New(settings), ssl.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "HOSTNAME PUBLIC_KEY_PATH PRIVATE_KEY_PATH [-s] [-r]"
		}
	},
}

// ICerts
type ICerts interface {
	Create(hostname, pubKey, privKey, svcID string) error
	Update(hostname, pubKey, privKey, svcID string) error
	List(svcID string) (*[]models.Cert, error)
	Rm(hostname, svcID string) error
}

// SCerts is a concrete implementation of ICerts
type SCerts struct {
	Settings *models.Settings
}

// New returns an instance of ICerts
func New(settings *models.Settings) ICerts {
	return &SCerts{
		Settings: settings,
	}
}
