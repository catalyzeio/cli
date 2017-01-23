package config

import (
	"errors"
	"runtime"

	"github.com/Sirupsen/logrus"
)

const (
	// VERSION is the current cli version
	VERSION = "4.0.0"
	// Beta determines whether or not this is a beta build of the CLI
	Beta = true
	// AccountsHost is the production accounts URL
	AccountsHost = "https://product.catalyze.io/stratum"
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
	// SkipVerifyEnvVar is the env variable used to accept invalid SSL certificates
	SkipVerifyEnvVar = "SKIP_VERIFY"

	// InvalidChars is a string containing all invalid characters for naming
	InvalidChars = "/?%"
)

// ErrEnvRequired is thrown when a command is run that requires an environment to be associated first
var ErrEnvRequired = errors.New("No Catalyze environment has been associated. Run \"catalyze associate\" from a local git repo first")

// ArchString translates the current architecture into an easier to read value.
// amd64 becomes 64-bit, 386 becomes 32-bit, etc.
func ArchString() string {
	archString := "other"
	switch runtime.GOARCH {
	case "386":
		archString = "32-bit"
	case "amd64":
		archString = "64-bit"
	case "arm":
		archString = "arm"
	}
	return archString
}
