package associate

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/models"
)

func CmdAssociate(envName, alias string, ia IAssociate, ie environments.IEnvironments) error {
	if envName == "" && alias != "" {
		return errors.New("An environment name is required When specifying an alias")
	}
	envs, errs := ie.List()
	if errs != nil && len(errs) > 0 {
		for pod, err := range errs {
			logrus.Debugf("Failed to list environments for pod \"%s\": %s", pod, err)
		}
	}
	found := false
	for _, env := range *envs {
		if envName == "" || env.Name == envName {
			found = true
			name := alias
			if envName == "" || name == "" {
				name = env.Name
			}
			err := ia.Associate(name, env.ID, env.Name, env.Pod, env.OrgID)
			if err != nil {
				return err
			}
			if envName != "" {
				break
			}
		}
	}
	if !found {
		return fmt.Errorf("No environment with name \"%s\" found", envName)
	}
	logrus.Println("After associating to an environment, you need to add a git repository with the \"catalyze git-remote add\" command and add a cert with the \"catalyze certs create\" command, if you have not done so already")
	return nil
}

// Associate an environment so that commands can be run against it. This command
// no longer adds a git remote. See commands.AddRemote().
func (s *SAssociate) Associate(alias, envID, envName, pod, orgID string) error {
	s.Settings.Environments[alias] = models.AssociatedEnvV2{
		EnvironmentID: envID,
		Name:          envName,
		Pod:           pod,
		OrgID:         orgID,
	}

	return nil
}
