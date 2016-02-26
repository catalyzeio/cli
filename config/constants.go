package config

import (
	"errors"

	"github.com/Sirupsen/logrus"
)

const (
	// VERSION is the current cli version
	VERSION = "3.0.0"
	// Beta determines whether or not this is a beta build of the CLI
	Beta = false
	// AccountsHost is the production accounts URL
	AccountsHost = "https://stratum.catalyze.io"
	// AuthHost is the production auth URL
	AuthHost = "https://auth.catalyze.io"
	// AuthHostVersion is the version path for the auth host
	AuthHostVersion = ""
	// PaasHost is the production PaaS URL
	PaasHost = "https://paas-api.catalyze.io"
	// PaasHostVersion is the version path for the PaaS host
	PaasHostVersion = ""
	// LogLevel is the amount of logging to enable
	LogLevel = logrus.InfoLevel
	// JobPollTime is the amount of time in seconds to wait between polls for a job status
	JobPollTime = 5
	// LogPollTime is the amount of time in seconds to wait between polls for new logs
	LogPollTime = 3

	// AccountsHostEnvVar is the env variable used to override AccountsHost
	AccountsHostEnvVar = "ACCOUNTS_HOST"
	// AuthHostEnvVar is the env variable used to override AuthHost
	AuthHostEnvVar = "AUTH_HOST"
	// PaasHostEnvVar is the env variable used to override PaasHost
	PaasHostEnvVar = "PAAS_HOST"
	// AuthHostVersionEnvVar is the env variable used to override AuthHostVersion
	AuthHostVersionEnvVar = "AUTH_HOST_VERSION"
	// PaasHostVersionEnvVar is the env variable used to override PaasHostVersion
	PaasHostVersionEnvVar = "PAAS_HOST_VERSION"
	// CatalyzeUsernameEnvVar is the env variable used to override the username
	CatalyzeUsernameEnvVar = "CATALYZE_USERNAME"
	// CatalyzePasswordEnvVar is the env variable used to override the passowrd
	CatalyzePasswordEnvVar = "CATALYZE_PASSWORD"
	// CatalyzeEnvironmentEnvVar is the env variable used to override the environment used in the current command
	CatalyzeEnvironmentEnvVar = "CATALYZE_ENV"
	// LogLevelEnvVar is the env variable used to override the logging level used
	LogLevelEnvVar = "CATALYZE_LOG_LEVEL"
)

// ErrEnvRequired is thrown when a command is run that requires an environment to be associated first
var ErrEnvRequired = errors.New("No Catalyze environment has been associated. Run \"catalyze associate\" from a local git repo first")
