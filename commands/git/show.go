package git

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
)

func CmdShow(svcName string, is services.IServices) error {
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
	logrus.Println(service.Source)
	return nil
}
