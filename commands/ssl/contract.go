package ssl

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
	"github.com/jawher/mow.cli"
)

// IncompleteChainError represents an error thrown when an SSL certificate is
// signed by an invalid CA or has an incomplete certificate chain.
type IncompleteChainError struct {
	Message string
	Err     error
}

func (i *IncompleteChainError) Error() string {
	return fmt.Sprintf("%s: %s", i.Message, i.Err.Error())
}

// HostnameMismatchError represents an error thrown when the hostname listed in
// an SSL certificate is different from the user input.
type HostnameMismatchError struct {
	Message string
	Err     error
}

func (h *HostnameMismatchError) Error() string {
	return fmt.Sprintf("%s: %s", h.Message, h.Err.Error())
}

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "ssl",
	ShortHelp: "Perform operations on local certificates to verify their validity",
	LongHelp:  "Perform operations on local certificates to verify their validity",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.Command(ResolveSubCmd.Name, ResolveSubCmd.ShortHelp, ResolveSubCmd.CmdFunc(settings))
			cmd.Command(VerifySubCmd.Name, VerifySubCmd.ShortHelp, VerifySubCmd.CmdFunc(settings))
		}
	},
}

var ResolveSubCmd = models.Command{
	Name:      "resolve",
	ShortHelp: "Verify that an SSL certificate is signed by a valid CA and attempt to resolve any incomplete certificate chains that are found",
	LongHelp:  "Verify that an SSL certificate is signed by a valid CA and attempt to resolve any incomplete certificate chains that are found",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			chain := subCmd.StringArg("CHAIN", "", "The path to your full certificate chain in PEM format")
			privateKey := subCmd.StringArg("PRIVATE_KEY", "", "The path to your private key in PEM format")
			hostname := subCmd.StringArg("HOSTNAME", "", "The hostname that should match your certificate (i.e. \"*.catalyze.io\")")
			output := subCmd.StringArg("OUTPUT", "", "The path of a file to save your properly resolved certificate chain (defaults to STDOUT)")
			force := subCmd.BoolOpt("f force", false, "If an output file is specified and already exists, setting force to true will overwrite the existing output file")
			subCmd.Action = func() {
				err := CmdResolve(*chain, *privateKey, *hostname, *output, *force, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "CHAIN PRIVATE_KEY HOSTNAME [OUTPUT] [-f]"
		}
	},
}

var VerifySubCmd = models.Command{
	Name:      "verify",
	ShortHelp: "Verify whether a certificate chain is complete and if it matches the given private key",
	LongHelp:  "Verify whether a certificate chain is complete and if it matches the given private key",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			chain := subCmd.StringArg("CHAIN", "", "The path to your full certificate chain in PEM format")
			privateKey := subCmd.StringArg("PRIVATE_KEY", "", "The path to your private key in PEM format")
			hostname := subCmd.StringArg("HOSTNAME", "", "The hostname that should match your certificate (i.e. \"*.catalyze.io\")")
			selfSigned := subCmd.BoolOpt("s self-signed", false, "Whether or not the certificate is self signed. If set, chain verification is skipped")
			subCmd.Action = func() {
				err := CmdVerify(*chain, *privateKey, *hostname, *selfSigned, New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "CHAIN PRIVATE_KEY HOSTNAME [-s]"
		}
	},
}

// ISSL
type ISSL interface {
	Verify(chainPath, privateKeyPath, hostname string, selfSigned bool) error
	Resolve(chainPath string) ([]byte, error)
}

// SSSL is a concrete implementation of ISSL
type SSSL struct {
	Settings *models.Settings
}

// New generates a new instance of ISSL
func New(settings *models.Settings) ISSL {
	return &SSSL{
		Settings: settings,
	}
}
