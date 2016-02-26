package associate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/config"
	"github.com/catalyzeio/cli/lib/git"
	"github.com/catalyzeio/cli/models"
)

func CmdAssociate(envLabel, svcLabel, alias, remote string, defaultEnv bool, ia IAssociate, ig git.IGit, ie environments.IEnvironments, is services.IServices) error {
	if !ig.Exists() {
		return errors.New("No git repo found in the current directory")
	}
	logrus.Printf("Existing git remotes named \"%s\" will be overwritten", remote)
	envs, err := ie.List()
	if err != nil {
		return err
	}
	var e *models.Environment
	var svcs *[]models.Service
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
		return fmt.Errorf("No environment with label \"%s\" found", envLabel)
	}
	if svcs == nil {
		return fmt.Errorf("No services found for environment with name \"%s\"", envLabel)
	}

	var chosenService *models.Service
	availableCodeServices := []string{}
	for _, service := range *svcs {
		if service.Type == "code" {
			if service.Label == svcLabel {
				chosenService = &service
				break
			}
			availableCodeServices = append(availableCodeServices, service.Label)
		}
	}
	if chosenService == nil {
		return fmt.Errorf("No code service found with name '%s'. Code services found: %s", svcLabel, strings.Join(availableCodeServices, ", "))
	}
	remotes, err := ig.List()
	if err != nil {
		return err
	}
	for _, r := range remotes {
		if r == remote {
			ig.Rm(remote)
			break
		}
	}
	err = ig.Add(remote, chosenService.Source)
	if err != nil {
		return err
	}
	logrus.Printf("\"%s\" remote added.", remote)

	name := alias
	if name == "" {
		name = envLabel
	}
	err = ia.Associate(name, remote, defaultEnv, e, chosenService)
	if err != nil {
		return err
	}
	logrus.Printf("Your git repository \"%s\" has been associated with code service \"%s\" and environment \"%s\"", remote, svcLabel, name)
	return nil
}

// Associate an environment so that commands can be run against it. This command
// no longer adds a git remote. See commands.AddRemote().
func (s *SAssociate) Associate(name, remote string, defaultEnv bool, env *models.Environment, chosenService *models.Service) error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	s.Settings.Environments[name] = models.AssociatedEnv{
		EnvironmentID: env.ID,
		ServiceID:     chosenService.ID,
		Directory:     dir,
		Name:          env.Name,
		Pod:           env.Pod,
		OrgID:         env.OrgID,
	}
	if defaultEnv {
		s.Settings.Default = name
	}
	config.DropBreadcrumb(name, s.Settings)
	if len(s.Settings.Environments) > 1 && s.Settings.Default == "" {
		logrus.Printf("You now have %d environments associated. Consider running \"catalyze default ENV_NAME\" to set a default", len(s.Settings.Environments))
	}

	return nil
}
