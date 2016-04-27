package files

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
)

func CmdEdit(svcName, fileName string, ifiles IFiles, is services.IServices) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return errors.New("\"EDITOR\" is a required environment variable for the \"catalyze files edit\" command")
	}
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services\" command.", svcName)
	}
	file, err := ifiles.Retrieve(fileName, service.ID)
	if err != nil {
		return err
	}
	if file == nil {
		return fmt.Errorf("File with name %s does not exist. Try listing files again by running \"catalyze files list %s\"", fileName, svcName)
	}
	tmpfile, err := ioutil.TempFile("", fmt.Sprintf("%d", file.ID))
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	err = ifiles.Save(tmpfile.Name(), true, file)
	if err != nil {
		return err
	}
	logrus.Debugf("File saved to %s", tmpfile.Name())

	logrus.Debugf("Opening with %s %s", editor, tmpfile.Name())
	cmd := exec.Command(editor, tmpfile.Name())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	err = cmd.Run()
	if err != nil {
		return err
	}
	logrus.Debugln("Editor returned, uploading...")

	err = ifiles.Update(service.ID, tmpfile.Name(), file.ID)
	if err != nil {
		return err
	}
	logrus.Printf("Successfully updated %s", fileName)
	logrus.Printf("To make your changes go live, you must redeploy your service with the \"catalyze redeploy %s\" command", service.Label)
	return nil
}
