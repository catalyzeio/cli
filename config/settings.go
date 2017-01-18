package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
	"github.com/mitchellh/go-homedir"
)

const (
	settingsFormatV1 = "v1"
	settingsFormatV2 = "v2"
)

const (
	// SettingsFile is the name of the catalyze config file.
	SettingsFile  = ".catalyze"
	currentFormat = settingsFormatV2
)

// SettingsRetriever defines an interface for a class responsible for generating
// a settings object used for most commands in the CLI. Some examples might be
// for retrieving settings based on the settings file or generating a settings
// object based on a directly entered environment ID and service ID.
type SettingsRetriever interface {
	GetSettings(string, string, string, string, string, string, string, string, string) (*models.Settings, error)
}

// FileSettingsRetriever reads in data from the SettingsFile and generates a
// settings object.
type FileSettingsRetriever struct{}

// GetSettings returns a Settings object for the current context
func (s FileSettingsRetriever) GetSettings(envName, svcName, accountsHost, authHost, ignoreAuthHostVersion, paasHost, ignorePaasHostVersion, username, password string) (*models.Settings, error) {
	HomeDir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filepath.Join(HomeDir, SettingsFile))
	if os.IsNotExist(err) {
		file, err = os.Create(filepath.Join(HomeDir, SettingsFile))
	}
	defer file.Close()
	if err != nil {
		return nil, err
	}
	var settings models.Settings
	json.NewDecoder(file).Decode(&settings)
	if settings.Format != currentFormat {
		if settings.Format == "" {
			settings.Format = "v1"
		}
		file.Seek(0, 0)
		settings, err = migrateSettings(file, settings.Format, currentFormat)
		if err != nil {
			return nil, err
		}
	}
	if settings.Environments == nil {
		settings.Environments = make(map[string]models.AssociatedEnvV2)
	}

	// try and set the given env first, if it exists
	if envName != "" {
		setGivenEnv(envName, &settings)
		if settings.EnvironmentID == "" {
			logrus.Fatalf("No environment named \"%s\" has been associated. Run \"catalyze associated\" to see what environments have been associated or run \"catalyze associate\" from a local git repo to create a new association", envName)
		}
	}

	settings.AccountsHost = accountsHost
	settings.AuthHost = authHost
	settings.PaasHost = paasHost
	settings.Username = username
	settings.Password = password

	authHostVersion := os.Getenv(AuthHostVersionEnvVar)
	if authHostVersion == "" {
		authHostVersion = AuthHostVersion
	}
	settings.AuthHostVersion = authHostVersion

	paasHostVersion := os.Getenv(PaasHostVersionEnvVar)
	if paasHostVersion == "" {
		paasHostVersion = PaasHostVersion
	}
	settings.PaasHostVersion = paasHostVersion

	logrus.Debugf("Accounts Host: %s", accountsHost)
	logrus.Debugf("Auth Host: %s", authHost)
	logrus.Debugf("Paas Host: %s", paasHost)
	logrus.Debugf("Auth Host Version: %s", authHostVersion)
	logrus.Debugf("Paas Host Version: %s", paasHostVersion)
	logrus.Debugf("Environment ID: %s", settings.EnvironmentID)
	logrus.Debugf("Environment Name: %s", settings.EnvironmentName)
	logrus.Debugf("Pod: %s", settings.Pod)
	logrus.Debugf("Org ID: %s", settings.OrgID)

	settings.Version = VERSION
	return &settings, nil
}

func migrateSettings(file *os.File, oldFormat, newFormat string) (models.Settings, error) {
	if oldFormat == "v1" {
		return migrateFromV1(file)
	}
	return models.Settings{}, fmt.Errorf("Invalid or corrupt settings file. Please fix the %s file in your home directory or contact Catalyze support", SettingsFile)
}

func migrateFromV1(file *os.File) (models.Settings, error) {
	logrus.Debugf("Migrating settings from %s to %s", settingsFormatV1, currentFormat)
	var oldSettings models.SettingsV1
	json.NewDecoder(file).Decode(&oldSettings)
	newSettings := models.Settings{
		PrivateKeyPath:  oldSettings.PrivateKeyPath,
		SessionToken:    oldSettings.SessionToken,
		UsersID:         oldSettings.UsersID,
		Environments:    map[string]models.AssociatedEnvV2{},
		Pods:            oldSettings.Pods,
		PodCheck:        oldSettings.PodCheck,
		Format:          currentFormat,
		AutoUpdateOptIn: models.AutoUpdateUnspecified,
	}
	for envName, env := range oldSettings.Environments {
		newSettings.Environments[envName] = models.AssociatedEnvV2{
			EnvironmentID: env.EnvironmentID,
			Name:          env.Name,
			Pod:           env.Pod,
			OrgID:         env.OrgID,
		}
	}
	return newSettings, nil
}

// SaveSettings persists the settings to disk
func SaveSettings(settings *models.Settings) error {
	HomeDir, err := homedir.Dir()
	if err != nil {
		return err
	}
	b, _ := json.Marshal(&settings)
	return ioutil.WriteFile(filepath.Join(HomeDir, SettingsFile), b, 0644)
}

// DeleteBreadcrumb removes the environment in the  global list
func DeleteBreadcrumb(alias string, settings *models.Settings) error {
	if _, ok := settings.Environments[alias]; !ok {
		return fmt.Errorf("An environment named \"%s\" has not been associated. Run \"catalyze associated\" to see current associations.", alias)
	}

	delete(settings.Environments, alias)
	SaveSettings(settings)
	return nil
}

// setGivenEnv takes the given env name and finds it in the env list
// in the given settings object. It then populates the EnvironmentID and
// ServiceID on the settings object with appropriate values.
func setGivenEnv(envName string, settings *models.Settings) {
	for eName, e := range settings.Environments {
		if eName == envName {
			settings.EnvironmentID = e.EnvironmentID
			settings.Pod = e.Pod
			settings.EnvironmentName = envName
			settings.OrgID = e.OrgID
			break
		}
	}
}

// defaultEnvPrompt asks the user when they dont have a default environment and
// aren't in an associated directory if they would like to proceed with the
// first environment found.
func defaultEnvPrompt(envName string) error {
	var answer string
	for {
		fmt.Printf("No environment was specified and no default environment was found. Falling back to %s\n", envName)
		fmt.Print("Do you wish to proceed? (y/n) ")
		fmt.Scanln(&answer)
		fmt.Println("")
		if answer != "y" && answer != "n" {
			fmt.Printf("%s is not a valid option. Please enter 'y' or 'n'\n", answer)
		} else {
			break
		}
	}
	if answer == "n" {
		return errors.New("Exiting")
	}
	return nil
}

// CheckRequiredAssociation ensures if an association is required for a command to run,
// that an appropriate environment has been picked and values assigned to the
// given settings object before a command is run. This is intended to be called
// before every command.
func CheckRequiredAssociation(required, prompt bool, settings *models.Settings) error {
	if required && settings.EnvironmentID == "" {
		err := ErrEnvRequired
		if prompt {
			for _, e := range settings.Environments {
				err = defaultEnvPrompt(e.Name)
				if err == nil {
					setGivenEnv(e.Name, settings)
				}
				break
			}
		}
		return err
	}
	return nil
}
