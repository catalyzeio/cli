package clear

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "clear",
	ShortHelp: "Clear out information in the global settings file to fix a misconfigured CLI.",
	LongHelp: "`clear` allows you to manage your global settings file in case your CLI becomes misconfigured. " +
		"The global settings file is stored in your home directory at `~/.catalyze`. " +
		"You can clear out all settings or pick and choose which ones need to be removed. " +
		"After running the `clear` command, any other CLI command will reset the removed settings to their appropriate values. Here are some sample commands\n\n" +
		"```\ncatalyze clear --all\n" +
		"catalyze clear --environments # removes your associated environments\n" +
		"catalyze clear --session --private-key # removes all session and private key authentication information\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			privateKey := cmd.BoolOpt("private-key", false, "Clear out the saved private key information")
			session := cmd.BoolOpt("session", false, "Clear out all session information")
			envs := cmd.BoolOpt("environments", false, "Clear out all associated environments")
			pods := cmd.BoolOpt("pods", false, "Clear out all saved pods")
			all := cmd.BoolOpt("all", false, "Clear out all settings")
			cmd.Action = func() {
				if *all {
					*privateKey = true
					*session = true
					*envs = true
					*pods = true
				}
				err := CmdClear(*privateKey, *session, *envs, *pods, settings)
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "[--private-key] [--session] [--environments] [--pods] [--all]"
		}
	},
}
