package associate

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/commands/git"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/models"
)

func CmdAssociate(envLabel, alias, remote string, ia IAssociate, ig git.IGit, ie environments.IEnvironments, is services.IServices) error {
	if !ig.Exists() {
		return errors.New("No git repo found in the current directory")
	}
	logrus.Printf("Existing git remotes named \"%s\" will be overwritten", remote)
	envs, errs := ie.List()
	if errs != nil && len(errs) > 0 {
		for pod, err := range errs {
			logrus.Debugf("Failed to list environments for pod \"%s\": %s", pod, err)
		}
	}
	var e *models.Environment
	var svcs *[]models.Service
	var err error
	for _, env := range *envs {
		if env.Name == envLabel {
			e = &env
			svcs, err = is.ListByEnvID(env.ID, env.Pod)
			if err != nil {
				return err
			}
			break
		}
	}
	if e == nil {
		return fmt.Errorf("No environment with name \"%s\" found", envLabel)
	}
	if svcs == nil {
		return fmt.Errorf("No services found for environment with name \"%s\"", envLabel)
	}

	name := alias
	if name == "" {
		name = envLabel
	}
	err = ia.Associate(name, remote, e)
	if err != nil {
		return err
	}
	logrus.Println("After associating to an environment, you need to add a git repository with the \"catalyze git-remote add\" command and add a cert with the \"catalyze certs create\" command, if you have not done so already")
	return nil
}

// Associate an environment so that commands can be run against it. This command
// no longer adds a git remote. See commands.AddRemote().
func (s *SAssociate) Associate(name, remote string, env *models.Environment) error {
	s.Settings.Environments[name] = models.AssociatedEnvV2{
		EnvironmentID: env.ID,
		Name:          env.Name,
		Pod:           env.Pod,
		OrgID:         env.OrgID,
	}

	return nil
}
