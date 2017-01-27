package logs

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/commands/sites"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "logs",
	ShortHelp: "Show the logs in your terminal streamed from your logging dashboard",
	LongHelp: "`logs` prints out your application logs directly from your logging Dashboard. " +
		"If you do not see your logs, try adjusting the number of hours, minutes, or seconds of logs that are retrieved with the `--hours`, `--minutes`, and `--seconds` options respectively. " +
		"You can also follow the logs with the `-f` option. " +
		"When using `-f` all logs will be printed to the console within the given time frame as well as any new logs that are sent to the logging Dashboard for the duration of the command. " +
		"When using the `-f` option, hit ctrl-c to stop. Here are some sample commands\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" logs --hours=6 --minutes=30\n" +
		"catalyze -E \"<your_env_alias>\" logs -f\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			query := cmd.StringArg("QUERY", "*", "The query to send to your logging dashboard's elastic search (regex is supported)")
			follow := cmd.BoolOpt("f follow", false, "Tail/follow the logs (Equivalent to -t)")
			tail := cmd.BoolOpt("t tail", false, "Tail/follow the logs (Equivalent to -f)")
			hours := cmd.IntOpt("hours", 0, "The number of hours before now (in combination with minutes and seconds) to retrieve logs")
			mins := cmd.IntOpt("minutes", 0, "The number of minutes before now (in combination with hours and seconds) to retrieve logs")
			secs := cmd.IntOpt("seconds", 0, "The number of seconds before now (in combination with hours and minutes) to retrieve logs")
			cmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdLogs(*query, *follow || *tail, *hours, *mins, *secs, settings.EnvironmentID, settings, New(settings), prompts.New(), environments.New(settings), services.New(settings), sites.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			cmd.Spec = "[QUERY] [(-f | -t)] [--hours] [--minutes] [--seconds]"
		}
	},
}

// ILogs ...
type ILogs interface {
	Output(queryString, sessionToken, domain string, from int, startTimestamp, endTimestamp time.Time) (int, error)
	Stream(queryString, sessionToken, domain string, from int, timestamp time.Time) error
	Watch(queryString, domain, sessionToken string) error
}

// SLogs is a concrete implementation of ILogs
type SLogs struct {
	Settings *models.Settings
}

// New returns an instance of ILogs
func New(settings *models.Settings) ILogs {
	return &SLogs{
		Settings: settings,
	}
}
