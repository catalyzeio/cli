package associated

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/models"
)

func CmdAssociated(ia IAssociated) error {
	envs, err := ia.Associated()
	if err != nil {
		return err
	}
	for envAlias, env := range envs {
		logrus.Printf(`%s:
    Environment ID:   %s
    Environment Name: %s
    Pod:              %s
    Organization ID:  %s
`, envAlias, env.EnvironmentID, env.Name, env.Pod, env.OrgID)
	}
	if len(envs) == 0 {
		logrus.Println("No environments have been associated. Run \"catalyze associate\" to get started.")
	}
	return nil
}

// Associated lists all currently associated environments.
func (a *SAssociated) Associated() (map[string]models.AssociatedEnvV2, error) {
	return a.Settings.Environments, nil
}
