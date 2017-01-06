package db

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/crypto"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/jault3/mow.cli"
)

// Cmd is the contract between the user and the CLI. This specifies the command
// name, arguments, and required/optional arguments and flags for the command.
var Cmd = models.Command{
	Name:      "db",
	ShortHelp: "Tasks for databases",
	LongHelp:  "The `db` command gives access to backup, import, and export services for databases. The db command can not be run directly but has sub commands.",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(cmd *cli.Cmd) {
			cmd.CommandLong(BackupSubCmd.Name, BackupSubCmd.ShortHelp, BackupSubCmd.LongHelp, BackupSubCmd.CmdFunc(settings))
			cmd.CommandLong(DownloadSubCmd.Name, DownloadSubCmd.ShortHelp, DownloadSubCmd.LongHelp, DownloadSubCmd.CmdFunc(settings))
			cmd.CommandLong(ExportSubCmd.Name, ExportSubCmd.ShortHelp, ExportSubCmd.LongHelp, ExportSubCmd.CmdFunc(settings))
			cmd.CommandLong(ImportSubCmd.Name, ImportSubCmd.ShortHelp, ImportSubCmd.LongHelp, ImportSubCmd.CmdFunc(settings))
			cmd.CommandLong(ListSubCmd.Name, ListSubCmd.ShortHelp, ListSubCmd.LongHelp, ListSubCmd.CmdFunc(settings))
			cmd.CommandLong(LogsSubCmd.Name, LogsSubCmd.ShortHelp, LogsSubCmd.LongHelp, LogsSubCmd.CmdFunc(settings))
		}
	},
}

var BackupSubCmd = models.Command{
	Name:      "backup",
	ShortHelp: "Create a new backup",
	LongHelp: "`db backup` creates a new backup for the given database service. " +
		"The backup is started and unless `-s` is specified, the CLI will poll every few seconds until it finishes. " +
		"Regardless of a successful backup or not, the logs for the backup will be printed to the console when the backup is finished. " +
		"If an error occurs and the logs are not printed, you can use the [db logs](#db-logs) command to print out historical backup job logs. Here is a sample command\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" db backup db01\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database service to create a backup for (i.e. 'db01')")
			skipPoll := subCmd.BoolOpt("s skip-poll", false, "Whether or not to wait for the backup to finish")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdBackup(*databaseName, *skipPoll, New(settings, crypto.New(), jobs.New(settings)), services.New(settings), jobs.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "DATABASE_NAME [-s]"
		}
	},
}

var DownloadSubCmd = models.Command{
	Name:      "download",
	ShortHelp: "Download a previously created backup",
	LongHelp: "`db download` downloads a previously created backup to your local hard drive. " +
		"Be careful using this command as it could download PHI. " +
		"Be sure that all hard drive encryption and necessary precautions have been taken before performing a download. " +
		"The ID of the backup is found by first running the [db list](#db-list) command. Here is a sample command\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" db download db01 cd2b4bce-2727-42d1-89e0-027bf3f1a203 ./db.sql\n```\n\n" +
		"This assumes you are downloading a MySQL or PostgreSQL backup which takes the `.sql` file format. If you are downloading a mongo backup, the command might look like this\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" db download db01 cd2b4bce-2727-42d1-89e0-027bf3f1a203 ./db.tar.gz\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database service which was backed up (i.e. 'db01')")
			backupID := subCmd.StringArg("BACKUP_ID", "", "The ID of the backup to download (found from \"catalyze backup list\")")
			filePath := subCmd.StringArg("FILEPATH", "", "The location to save the downloaded backup to. This location must NOT already exist unless -f is specified")
			force := subCmd.BoolOpt("f force", false, "If a file previously exists at \"filepath\", overwrite it and download the backup")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdDownload(*databaseName, *backupID, *filePath, *force, New(settings, crypto.New(), jobs.New(settings)), prompts.New(), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "DATABASE_NAME BACKUP_ID FILEPATH [-f]"
		}
	},
}

var ExportSubCmd = models.Command{
	Name:      "export",
	ShortHelp: "Export data from a database",
	LongHelp: "`db export` is a simple wrapper around the `db backup` and `db download` commands. " +
		"When you request an export, a backup is created that will be added to the list of backups shown when you perform the [db list](#db-list) command. " +
		"Then that backup is immediately downloaded. Regardless of a successful export or not, the logs for the backup will be printed to the console when the export is finished. " +
		"If an error occurs and the logs are not printed, you can use the [db logs](#db-logs) command to print out historical backup job logs. Here is a sample command\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" db export db01 ./dbexport.sql\n```\n\n" +
		"This assumes you are exporting a MySQL or PostgreSQL database which takes the `.sql` file format. If you are exporting a mongo database, the command might look like this\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" db export db01 ./dbexport.tar.gz\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database to export data from (i.e. 'db01')")
			filePath := subCmd.StringArg("FILEPATH", "", "The location to save the exported data. This location must NOT already exist unless -f is specified")
			force := subCmd.BoolOpt("f force", false, "If a file previously exists at `filepath`, overwrite it and export data")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdExport(*databaseName, *filePath, *force, New(settings, crypto.New(), jobs.New(settings)), prompts.New(), services.New(settings), jobs.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "DATABASE_NAME FILEPATH [-f]"
		}
	},
}

var ImportSubCmd = models.Command{
	Name:      "import",
	ShortHelp: "Import data into a database",
	LongHelp: "`db import` allows you to inject new data into your database service. For example, if you wrote a simple SQL file\n\n" +
		"```\nCREATE TABLE mytable (\n" +
		"id TEXT PRIMARY KEY,\n" +
		"val TEXT\n" +
		");\n\n" +
		"INSERT INTO mytable (id, val) values ('1', 'test');\n```\n\n" +
		"and stored it at `./db.sql` you could import this into your database service. " +
		"When importing data into mongo, you may specify the database and collection to import into using the `-d` and `-c` flags respectively. " +
		"Regardless of a successful import or not, the logs for the import will be printed to the console when the import is finished. " +
		"Before an import takes place, your database is backed up automatically in case any issues arise. Here is a sample command\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" db import db01 ./db.sql\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database to import data to (i.e. 'db01')")
			filePath := subCmd.StringArg("FILEPATH", "", "The location of the file to import to the database")
			mongoCollection := subCmd.StringOpt("c mongo-collection", "", "If importing into a mongo service, the name of the collection to import into")
			mongoDatabase := subCmd.StringOpt("d mongo-database", "", "If importing into a mongo service, the name of the database to import into")
			skipBackup := subCmd.BoolOpt("s skip-backup", false, "Skip backing up database. Useful for large databases, which can have long backup times.")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdImport(*databaseName, *filePath, *mongoCollection, *mongoDatabase, *skipBackup, New(settings, crypto.New(), jobs.New(settings)), prompts.New(), services.New(settings), jobs.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "DATABASE_NAME FILEPATH [-s][-d [-c]]"
		}
	},
}

var ListSubCmd = models.Command{
	Name:      "list",
	ShortHelp: "List created backups",
	LongHelp: "`db list` lists all previously created backups. " +
		"After listing backups you can copy the backup ID and use it to [download](#db-download) that backup or [view the logs](#db-logs) from that backup. Here is a sample command\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" db list db01\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database service to list backups for (i.e. 'db01')")
			page := subCmd.IntOpt("p page", 1, "The page to view")
			pageSize := subCmd.IntOpt("n page-size", 10, "The number of items to show per page")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdList(*databaseName, *page, *pageSize, New(settings, crypto.New(), jobs.New(settings)), services.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "DATABASE_NAME [-p] [-n]"
		}
	},
}

var LogsSubCmd = models.Command{
	Name:      "logs",
	ShortHelp: "Print out the logs from a previous database backup job",
	LongHelp: "`db logs` allows you to view backup logs from historical backup jobs. " +
		"You can find the backup ID from using the `db list` command. Here is a sample command\n\n" +
		"```\ncatalyze -E \"<your_env_alias>\" db logs db01 cd2b4bce-2727-42d1-89e0-027bf3f1a203\n```",
	CmdFunc: func(settings *models.Settings) func(cmd *cli.Cmd) {
		return func(subCmd *cli.Cmd) {
			databaseName := subCmd.StringArg("DATABASE_NAME", "", "The name of the database service (i.e. 'db01')")
			backupID := subCmd.StringArg("BACKUP_ID", "", "The ID of the backup to download logs from (found from \"catalyze backup list\")")
			subCmd.Action = func() {
				if _, err := auth.New(settings, prompts.New()).Signin(); err != nil {
					logrus.Fatal(err.Error())
				}
				if err := config.CheckRequiredAssociation(true, true, settings); err != nil {
					logrus.Fatal(err.Error())
				}
				err := CmdLogs(*databaseName, *backupID, New(settings, crypto.New(), jobs.New(settings)), services.New(settings), jobs.New(settings))
				if err != nil {
					logrus.Fatal(err.Error())
				}
			}
			subCmd.Spec = "DATABASE_NAME BACKUP_ID"
		}
	},
}

// IDb
type IDb interface {
	Backup(service *models.Service) (*models.Job, error)
	Download(backupID, filePath string, service *models.Service) error
	Export(filePath string, job *models.Job, service *models.Service) error
	Import(filePath, mongoCollection, mongoDatabase string, service *models.Service) (*models.Job, error)
	List(page, pageSize int, service *models.Service) (*[]models.Job, error)
	TempUploadAuth(service *models.Service) (*models.TempAuth, error)
	TempDownloadURL(jobID string, service *models.Service) (*models.TempURL, error)
	TempLogsURL(jobID string, serviceID string) (*models.TempURL, error)
	DumpLogs(taskType string, job *models.Job, service *models.Service) error
}

// SDb is a concrete implementation of IDb
type SDb struct {
	Settings *models.Settings
	Crypto   crypto.ICrypto
	Jobs     jobs.IJobs
}

// New returns an instance of IDb
func New(settings *models.Settings, crypto crypto.ICrypto, jobs jobs.IJobs) IDb {
	return &SDb{
		Settings: settings,
		Crypto:   crypto,
		Jobs:     jobs,
	}
}
