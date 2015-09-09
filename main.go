package main

import (
	"fmt"
	"os"

	"github.com/catalyzeio/catalyze/commands"
	"github.com/catalyzeio/catalyze/config"
	"github.com/jawher/mow.cli"
)

// current cli version
const VERSION = "2.0.0"

// alternate BaaS URL
const BaasHost = "https://api.catalyze.io"

// alternate PaaS URL
const PaasHost = "https://paas-api.catalyze.io"

func main() {
	var app = cli.App("catalyze", fmt.Sprintf("Catalyze CLI. Version %s", VERSION))
	InitCLI(app)

	versionFlag := app.Bool(cli.BoolOpt{
		Name:      "version",
		Desc:      "CLI Version",
		HideValue: true,
	})

	app.Action = func() {
		// just to make this function like a normal CLI, if specifying
		// catalyze --version, output version and quit
		if *versionFlag {
			version()
		} else {
			app.PrintHelp()
		}
	}

	app.Command("version", "Output the version and quit", func(cmd *cli.Cmd) {
		cmd.Action = version
	})

	app.Run(os.Args)
}

// InitCLI adds arguments and commands to the given cli instance
func InitCLI(app *cli.Cli) {
	baasHost := os.Getenv("BAAS_HOST")
	if baasHost == "" {
		baasHost = BaasHost
	}
	paasHost := os.Getenv("PAAS_HOST")
	if paasHost == "" {
		paasHost = PaasHost
	}
	username := app.String(cli.StringOpt{
		Name:      "U username",
		Desc:      "Catalyze Username",
		EnvVar:    "CATALYZE_USERNAME",
		HideValue: true,
	})
	password := app.String(cli.StringOpt{
		Name:      "P password",
		Desc:      "Catalyze Password",
		EnvVar:    "CATALYZE_PASSWORD",
		HideValue: true,
	})
	givenEnvName := app.String(cli.StringOpt{
		Name:      "E env",
		Desc:      "The local alias of the environment in which this command will be run",
		EnvVar:    "CATALYZE_ENV",
		HideValue: true,
	})

	app.Command("associate", "Associates an environment", func(cmd *cli.Cmd) {
		envName := cmd.StringArg("ENV_NAME", "", "The name of your environment")
		serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the primary code service to associate with this environment (i.e. 'app01')")
		alias := cmd.StringOpt("a alias", "", "A shorter name to reference your environment by for local commands")
		remote := cmd.StringOpt("r remote", "catalyze", "The name of the remote")
		defaultEnv := cmd.BoolOpt("d default", false, "Specifies whether or not the associated environment will be the default")
		cmd.Action = func() {
			settings := config.GetSettings(false, false, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.Associate(*envName, *serviceName, *alias, *remote, *defaultEnv, settings)
		}
		cmd.Spec = "ENV_NAME [SERVICE_NAME] [-a] [-r] [-d]"
	})
	app.Command("associated", "Lists all associated environments", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			settings := config.GetSettings(false, false, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.Associated(settings)
		}
	})
	app.Command("backup", "Backup and restore services for databases", func(cmd *cli.Cmd) {
		cmd.Command("create", "Create a new backup", func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the database service to create a backup for (i.e. 'db01')")
			skipPoll := subCmd.BoolOpt("s skip-poll", false, "Whether or not to wait for the backup to finish")
			subCmd.Action = func() {
				settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
				commands.CreateBackup(*serviceName, *skipPoll, settings)
			}
			subCmd.Spec = "SERVICE_NAME [-s]"
		})
		cmd.Command("download", "Download a previously created backup", func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the database service which was backed up (i.e. 'db01')")
			backupID := subCmd.StringArg("BACKUP_ID", "", "The ID of the backup to download (found from \"catalyze backup list\")")
			filePath := subCmd.StringArg("FILEPATH", "", "The location to save the downloaded backup to. This location must NOT already exist unless -f is specified")
			force := subCmd.BoolOpt("f force", false, "If a file previously exists at \"filepath\", overwrite it and download the backup")
			subCmd.Action = func() {
				settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
				commands.DownloadBackup(*serviceName, *backupID, *filePath, *force, settings)
			}
			subCmd.Spec = "SERVICE_NAME BACKUP_ID FILEPATH [-f]"
		})
		cmd.Command("list", "List created backups", func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the database service to list backups for (i.e. 'db01')")
			page := subCmd.IntOpt("p page", 1, "The page to view")
			pageSize := subCmd.IntOpt("n page-size", 10, "The number of items to show per page")
			subCmd.Action = func() {
				settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
				commands.ListBackups(*serviceName, *page, *pageSize, settings)
			}
			subCmd.Spec = "SERVICE_NAME [-p] [-n]"
		})
		cmd.Command("restore", "Restore from a previously created backup", func(subCmd *cli.Cmd) {
			serviceName := subCmd.StringArg("SERVICE_NAME", "", "The name of the database service to restore (i.e. 'db01')")
			backupID := subCmd.StringArg("BACKUP_ID", "", "The ID of the backup to restore (found from `catalyze backup list`)")
			skipPoll := subCmd.BoolOpt("s skip-poll", false, "Whether or not to wait for the restore to finish")
			subCmd.Action = func() {
				settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
				commands.RestoreBackup(*serviceName, *backupID, *skipPoll, settings)
			}
			subCmd.Spec = "SERVICE_NAME BACKUP_ID [-s]"
		})
	})
	app.Command("console", "Open a secure console to a service", func(cmd *cli.Cmd) {
		serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to open up a console for")
		command := cmd.StringArg("COMMAND", "", "An optional command to run when the console becomes available")
		cmd.Action = func() {
			settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.Console(*serviceName, *command, settings)
		}
		cmd.Spec = "SERVICE_NAME [COMMAND]"
	})
	app.Command("dashboard", "Open the Catalyze Dashboard in your default browser", func(cmd *cli.Cmd) {
		cmd.Action = commands.Dashboard
	})
	app.Command("db", "Import and export services for databases", func(cmd *cli.Cmd) {
		cmd.Command("import", "Import data to a database", func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database to import data to (i.e. 'db01')")
			filePath := subCmd.StringArg("FILEPATH", "", "The location of the file to import to the database")
			mongoCollection := subCmd.StringOpt("c mongo-collection", "", "If importing into a mongo service, the name of the collection to import into")
			mongoDatabase := subCmd.StringOpt("d mongo-database", "", "If importing into a mongo service, the name of the database to import into")
			wipeFirst := subCmd.BoolOpt("w wipe-first", false, "Whether or not to wipe the database before processing the import file")
			subCmd.Action = func() {
				settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
				commands.Import(*databaseName, *filePath, *mongoCollection, *mongoDatabase, *wipeFirst, settings)
			}
			subCmd.Spec = "DATABASE_NAME FILEPATH [-w] [-d [-c]]"
		})
		cmd.Command("export", "Export data from a database", func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database to export data from (i.e. 'db01')")
			filePath := subCmd.StringArg("FILEPATH", "", "The location to save the exported data. This location must NOT already exist unless -f is specified")
			force := subCmd.BoolOpt("f force", false, "If a file previously exists at `filepath`, overwrite it and export data")
			subCmd.Action = func() {
				settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
				commands.Export(*databaseName, *filePath, *force, settings)
			}
			subCmd.Spec = "DATABASE_NAME FILEPATH [-f]"
		})
	})
	app.Command("default", "Set the default associated environment", func(cmd *cli.Cmd) {
		envAlias := cmd.StringArg("ENV_ALIAS", "", "The alias of an already associated environment to set as the default")
		cmd.Action = func() {
			settings := config.GetSettings(true, false, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.SetDefault(*envAlias, settings)
		}
		cmd.Spec = "ENV_ALIAS"
	})
	app.Command("disassociate", "Remove the association with an environment", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			settings := config.GetSettings(true, false, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.Disassociate(settings)
		}
	})
	app.Command("environments", "List all environments you have access to", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			settings := config.GetSettings(false, false, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.Environments(settings)
		}
	})
	app.Command("logs", "Show the logs in your terminal streamed from your logging dashboard", func(cmd *cli.Cmd) {
		query := cmd.StringArg("QUERY", "app*", "The query to send to your logging dashboard's elastic search (regex is supported)")
		follow := cmd.BoolOpt("f follow", false, "Tail/follow the logs (Equivalent to -t)")
		tail := cmd.BoolOpt("t tail", false, "Tail/follow the logs (Equivalent to -f)")
		hours := cmd.IntOpt("h hours", 4, "The number of hours to retrieve logs")
		cmd.Action = func() {
			settings := config.GetSettings(false, false, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.Logs(*query, *tail || *follow, *hours, settings)
		}
		cmd.Spec = "[QUERY] [(-f | -t)] [-h]"
	})
	app.Command("logout", "Clear the stored user information from your local machine", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			settings := config.GetSettings(false, false, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.Logout(settings)
		}
	})
	app.Command("metrics", "Print service and environment metrics in your local time zone", func(cmd *cli.Cmd) {
		serviceName := cmd.StringArg("SERVICE_NAME", "", "The name of the service to print metrics for")
		json := cmd.BoolOpt("json", false, "Output the data as json")
		csv := cmd.BoolOpt("csv", false, "Output the data as csv")
		spark := cmd.BoolOpt("spark", false, "Output the data using spark lines")
		stream := cmd.BoolOpt("stream", false, "Repeat calls once per minute until this process is interrupted.")
		mins := cmd.IntOpt("m mins", 1, "How many minutes worth of logs to retrieve.")
		cmd.Action = func() {
			settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.Metrics(*serviceName, *json, *csv, *spark, *stream, *mins, settings)
		}
		cmd.Spec = "[SERVICE_NAME] [(--json | --csv | --spark)] [--stream] [-m]"
	})
	app.Command("rake", "Execute a rake task", func(cmd *cli.Cmd) {
		taskName := cmd.StringArg("TASK_NAME", "", "The name of the rake task to run")
		cmd.Action = func() {
			settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.Rake(*taskName, settings)
		}
		cmd.Spec = "TASK_NAME"
	})
	app.Command("redeploy", "Redeploy a service without having to do a git push", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.Redeploy(settings)
		}
	})
	app.Command("ssl", "Perform operations on local certificates to verify their validity", func(cmd *cli.Cmd) {
		cmd.Command("verify", "Verify whether a certificate chain is complete and if it matches the given private key", func(subCmd *cli.Cmd) {
			chain := subCmd.StringArg("CHAIN", "", "The path to your full certificate chain in PEM format")
			privateKey := subCmd.StringArg("PRIVATE_KEY", "", "The path to your private key in PEM format")
			hostname := subCmd.StringArg("HOSTNAME", "", "The hostname that should match your certificate (i.e. \"*.catalyze.io\")")
			selfSigned := subCmd.BoolOpt("s self-signed", false, "Whether or not the certificate is self signed. If set, chain verification is skipped")
			subCmd.Action = func() {
				commands.VerifyChain(*chain, *privateKey, *hostname, *selfSigned)
			}
			subCmd.Spec = "CHAIN PRIVATE_KEY HOSTNAME [-s]"
		})
	})
	app.Command("status", "Get quick readout of the current status of your associated environment and all of its services", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.Status(settings)
		}
	})
	app.Command("support-ids", "Print out various IDs related to your associated environment to be used when contacting Catalyze support", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.SupportIds(settings)
		}
	})
	app.Command("users", "Manage users who have access to the associated environment", func(cmd *cli.Cmd) {
		cmd.Command("add", "Grant access to the associated environment for the given user", func(subCmd *cli.Cmd) {
			usersID := subCmd.StringArg("USER_ID", "", "The Users ID to give access to the associated environment")
			subCmd.Action = func() {
				settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
				commands.AddUser(*usersID, settings)
			}
			subCmd.Spec = "USER_ID"
		})
		cmd.Command("list", "List all users who have access to the associated environment", func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
				commands.ListUsers(settings)
			}
		})
		cmd.Command("rm", "Revoke access to the associated environment for the given user", func(subCmd *cli.Cmd) {
			usersID := subCmd.StringArg("USER_ID", "", "The Users ID to revoke access from for the associated environment")
			subCmd.Action = func() {
				settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
				commands.RmUser(*usersID, settings)
			}
			subCmd.Spec = "USER_ID"
		})
	})
	app.Command("vars", "Interaction with environment variables for the associated environment", func(cmd *cli.Cmd) {
		cmd.Command("list", "List all environment variables", func(subCmd *cli.Cmd) {
			subCmd.Action = func() {
				settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
				commands.ListVars(settings)
			}
		})
		cmd.Command("set", "Set one or more new environment variables or update the values of existing ones", func(subCmd *cli.Cmd) {
			variables := subCmd.Strings(cli.StringsOpt{
				Name:      "v variable",
				Value:     []string{},
				Desc:      "The env variable to set or update in the form \"<key>=<value>\"",
				HideValue: true,
			})
			subCmd.Action = func() {
				settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
				commands.SetVar(*variables, settings)
			}
			subCmd.Spec = "-v..."
		})
		cmd.Command("unset", "Unset (delete) an existing environment variable", func(subCmd *cli.Cmd) {
			variable := subCmd.StringArg("VARIABLE", "", "The name of the environment variable to unset")
			subCmd.Action = func() {
				settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
				commands.UnsetVar(*variable, settings)
			}
			subCmd.Spec = "VARIABLE"
		})
	})
	app.Command("whoami", "Retrieve your user ID", func(cmd *cli.Cmd) {
		cmd.Action = func() {
			settings := config.GetSettings(false, false, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.WhoAmI(settings)
		}
	})
	app.Command("worker", "Start a background worker", func(cmd *cli.Cmd) {
		target := cmd.StringArg("TARGET", "", "The name of the Procfile target to invoke as a worker")
		cmd.Action = func() {
			settings := config.GetSettings(true, true, *givenEnvName, baasHost, paasHost, *username, *password)
			commands.Worker(*target, settings)
		}
		cmd.Spec = "TARGET"
	})
}

func version() {
	fmt.Printf("version %s\n", VERSION)
}
