package git

import (
	"fmt"
	"os/exec"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
)

func CmdAdd(svcName, remote string, ig IGit, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	if service.Source == "" {
		return fmt.Errorf("No git remote found for the \"%s\" service.", svcName)
	}
	err = ig.Add(remote, service.Source)
	if err != nil {
		return err
	}
	logrus.Printf("\"%s\" remote added.", remote)
	return nil
}

// Add a git remote to a git repo in the current working directory. If the
// current working directory is not yet a git repo, this command will fail.
func (g *SGit) Add(remote, gitURL string) error {
	_, err := exec.Command("git", "remote", "add", remote, gitURL).Output()
	return err
}
