package updater

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/pods"
	"github.com/catalyzeio/cli/models"
)

func updatePods() {
	settings := &models.Settings{}
	r := config.FileSettingsRetriever{}
	*settings = *r.GetSettings("", "", config.AccountsHost, config.AuthHost, "", config.PaasHost, "", config.CatalyzeUsernameEnvVar, config.CatalyzePasswordEnvVar)
	p := pods.New(settings)
	pods, err := p.List()
	if err == nil {
		settings.Pods = pods
		logrus.Debugf("%+v", settings.Pods)
		config.SaveSettings(settings)
		fmt.Println("hiiiii")
	} else {
		settings.Pods = &[]models.Pod{}
		logrus.Debugf("Error listing pods: %s", err.Error())
	}
}
