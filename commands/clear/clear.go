package clear

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/models"
)

func CmdClear(privateKey, session, environments, pods bool, settings *models.Settings) error {
	if privateKey {
		settings.PrivateKeyPath = ""
	}
	if session {
		settings.SessionToken = ""
		settings.UsersID = ""
	}
	if environments {
		settings.Environments = map[string]models.AssociatedEnvV2{}
	}
	if pods {
		settings.Pods = &[]models.Pod{}
	}
	config.SaveSettings(settings)
	if !privateKey && !session && !environments && !pods {
		logrus.Println("No settings were specified. To see available options, run \"catalyze clear --help\"")
	} else {
		logrus.Println("All specified settings have been cleared")
	}
	return nil
}
