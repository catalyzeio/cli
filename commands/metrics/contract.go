package metrics

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jault3/mow.cli"
)

type MetricType uint8

const (
	CPU MetricType = iota
	Memory
	NetworkIn
	NetworkOut
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "metrics",
	ShortHelp: "Print service and environment metrics in your local time zone",
	LongHelp: "The `metrics` command gives access to environment metrics or individual service metrics through a variety of formats. " +
		"This is useful for checking on the status and performance of your application or environment as a whole. " +
		"The metrics command cannot be run directly but has sub commands.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.CommandLong(CPUSubCmd.Name, CPUSubCmd.ShortHelp, CPUSubCmd.LongHelp, CPUSubCmd.CmdFunc(settings))
			cmd.CommandLong(MemorySubCmd.Name, MemorySubCmd.ShortHelp, MemorySubCmd.LongHelp, MemorySubCmd.CmdFunc(settings))
			cmd.CommandLong(NetworkInSubCmd.Name, NetworkInSubCmd.ShortHelp, NetworkInSubCmd.LongHelp, NetworkInSubCmd.CmdFunc(settings))
			cmd.CommandLong(NetworkOutSubCmd.Name, NetworkOutSubCmd.ShortHelp, NetworkOutSubCmd.LongHelp, NetworkOutSubCmd.CmdFunc(settings))
		}
	},
}

var CPUSubCmd = models.Command{
	Name:      "cpu",
	ShortHelp: "Print service and environment CPU metrics in your local time zone",
	LongHelp: "`metrics cpu` prints out CPU metrics for your environment or individual services. " +
		"You can print out metrics in csv, json, plain text, or spark lines format. " +
		"If you want plain text format, simply omit the `--json` and `--csv` flags. " +
		"You can only stream metrics using plain text or spark lines formats. " +
		"To print out metrics for every service in your environment, omit the `SERVICE_NAME` argument. " +
		"Otherwise you may choose a service, such as an app service, to retrieve metrics for. " +
		"Here are some sample commands\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" metrics cpu\n" +
		"catalyze -E \"<your_env_alias>\" metrics cpu app01 --stream\n" +
		"catalyze -E \"<your_env_alias>\" metrics cpu --json\n" +
		"catalyze -E \"<your_env_alias>\" metrics cpu db01 --csv -m 60\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to print metrics for")
			json := subCmd.BoolOpt("json", false, "Output the data as json")
			csv := subCmd.BoolOpt("csv", false, "Output the data as csv")
			text := subCmd.BoolOpt("text", true, "Output the data in plain text")
			stream := subCmd.BoolOpt("stream", false, "Repeat calls once per minute until this process is interrupted.")
			mins := subCmd.IntOpt("m mins", 1, "How many minutes worth of metrics to retrieve.")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdMetrics(*serviceName, CPU, *json, *csv, *text, *stream, *mins, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "[SERVICE_NAME] [(--json | --csv | --text)] [--stream] [-m]"
		}
	},
}

var MemorySubCmd = models.Command{
	Name:      "memory",
	ShortHelp: "Print service and environment memory metrics in your local time zone",
	LongHelp: "`metrics memory` prints out memory metrics for your environment or individual services. " +
		"You can print out metrics in csv, json, plain text, or spark lines format. " +
		"If you want plain text format, simply omit the `--json` and `--csv` flags. " +
		"You can only stream metrics using plain text or spark lines formats. " +
		"To print out metrics for every service in your environment, omit the `SERVICE_NAME` argument. " +
		"Otherwise you may choose a service, such as an app service, to retrieve metrics for. " +
		"Here are some sample commands\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" metrics memory\n" +
		"catalyze -E \"<your_env_alias>\" metrics memory app01 --stream\n" +
		"catalyze -E \"<your_env_alias>\" metrics memory --json\n" +
		"catalyze -E \"<your_env_alias>\" metrics memory db01 --csv -m 60\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to print metrics for")
			json := subCmd.BoolOpt("json", false, "Output the data as json")
			csv := subCmd.BoolOpt("csv", false, "Output the data as csv")
			text := subCmd.BoolOpt("text", true, "Output the data in plain text")
			stream := subCmd.BoolOpt("stream", false, "Repeat calls once per minute until this process is interrupted.")
			mins := subCmd.IntOpt("m mins", 1, "How many minutes worth of metrics to retrieve.")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdMetrics(*serviceName, Memory, *json, *csv, *text, *stream, *mins, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "[SERVICE_NAME] [(--json | --csv | --text)] [--stream] [-m]"
		}
	},
}

var NetworkInSubCmd = models.Command{
	Name:      "network-in",
	ShortHelp: "Print service and environment received network data metrics in your local time zone",
	LongHelp: "`metrics network-in` prints out received network metrics for your environment or individual services. " +
		"You can print out metrics in csv, json, plain text, or spark lines format. " +
		"If you want plain text format, simply omit the `--json` and `--csv` flags. " +
		"You can only stream metrics using plain text or spark lines formats. " +
		"To print out metrics for every service in your environment, omit the `SERVICE_NAME` argument. " +
		"Otherwise you may choose a service, such as an app service, to retrieve metrics for. Here are some sample commands\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" metrics network-in\n" +
		"catalyze -E \"<your_env_alias>\" metrics network-in app01 --stream\n" +
		"catalyze -E \"<your_env_alias>\" metrics network-in --json\n" +
		"catalyze -E \"<your_env_alias>\" metrics network-in db01 --csv -m 60\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to print metrics for")
			json := subCmd.BoolOpt("json", false, "Output the data as json")
			csv := subCmd.BoolOpt("csv", false, "Output the data as csv")
			text := subCmd.BoolOpt("text", true, "Output the data in plain text")
			stream := subCmd.BoolOpt("stream", false, "Repeat calls once per minute until this process is interrupted.")
			mins := subCmd.IntOpt("m mins", 1, "How many minutes worth of metrics to retrieve.")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdMetrics(*serviceName, NetworkIn, *json, *csv, *text, *stream, *mins, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "[SERVICE_NAME] [(--json | --csv | --text)] [--stream] [-m]"
		}
	},
}

var NetworkOutSubCmd = models.Command{
	Name:      "network-out",
	ShortHelp: "Print service and environment transmitted network data metrics in your local time zone",
	LongHelp: "`metrics network-out` prints out transmitted network metrics for your environment or individual services. " +
		"You can print out metrics in csv, json, plain text, or spark lines format. " +
		"If you want plain text format, simply omit the `--json` and `--csv` flags. " +
		"You can only stream metrics using plain text or spark lines formats. " +
		"To print out metrics for every service in your environment, omit the `SERVICE_NAME` argument. " +
		"Otherwise you may choose a service, such as an app service, to retrieve metrics for. " +
		"Here are some sample commands\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" metrics network-out\n" +
		"catalyze -E \"<your_env_alias>\" metrics network-out app01 --stream\n" +
		"catalyze -E \"<your_env_alias>\" metrics network-out --json\n" +
		"catalyze -E \"<your_env_alias>\" metrics network-out db01 --csv -m 60\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the service to print metrics for")
			json := subCmd.BoolOpt("json", false, "Output the data as json")
			csv := subCmd.BoolOpt("csv", false, "Output the data as csv")
			text := subCmd.BoolOpt("text", true, "Output the data in plain text")
			stream := subCmd.BoolOpt("stream", false, "Repeat calls once per minute until this process is interrupted.")
			mins := subCmd.IntOpt("m mins", 1, "How many minutes worth of metrics to retrieve.")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdMetrics(*serviceName, NetworkOut, *json, *csv, *text, *stream, *mins, New(settings), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "[SERVICE_NAME] [(--json | --csv | --text)] [--stream] [-m]"
		}
	},
}

// IMetrics
type IMetrics interface {
	RetrieveEnvironmentMetrics(mins int) (*[]models.Metrics, error)
	RetrieveServiceMetrics(mins int, svcID string) (*models.Metrics, error)
}

// SMetrics is a concrete implementation of IMetrics
type SMetrics struct {
	Settings *models.Settings
}

// New returns an instance of IMetrics
func New(settings *models.Settings) IMetrics {
	return &SMetrics{
		Settings: settings,
	}
}
