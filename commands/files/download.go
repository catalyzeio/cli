package files

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/models"
)

// CmdDownload downloads a service file by its name (taken from listing
// service files) to the local machine matching the file's assigned permissions.
// If those permissions cannot be applied, the default 0644 permissions are
// applied. If not output file is specified, the file and permissions are
// printed to stdout.
func CmdDownload(svcName, fileName, output string, force bool, ifiles IFiles, is services.IServices) error {
	if output != "" && !force {
		if _, err := os.Stat(output); err == nil {
			return fmt.Errorf("File already exists at path '%s'. Specify '--force' to overwrite", output)
		}
	}
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	file, err := ifiles.Retrieve(fileName, service.ID)
	if err != nil {
		return err
	}
	if file == nil {
		return fmt.Errorf("File with name %s does not exist. Try listing files again by running \"catalyze files list %s\"", fileName, svcName)
	}
	return ifiles.Save(output, force, file)
}

func (f *SFiles) Retrieve(fileName string, svcID string) (*models.ServiceFile, error) {
	files, err := f.List(svcID)
	if err != nil {
		return nil, err
	}
	for _, ff := range *files {
		if ff.Name == fileName {
			headers := f.Settings.HTTPManager.GetHeaders(f.Settings.SessionToken, f.Settings.Version, f.Settings.Pod, f.Settings.UsersID)
			resp, statusCode, err := f.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/files/%d", f.Settings.PaasHost, f.Settings.PaasHostVersion, f.Settings.EnvironmentID, svcID, ff.ID), headers)
			if err != nil {
				return nil, err
			}
			var file models.ServiceFile
			err = f.Settings.HTTPManager.ConvertResp(resp, statusCode, &file)
			if err != nil {
				return nil, err
			}
			return &file, nil
		}
	}
	return nil, nil
}

func (f *SFiles) Save(output string, force bool, file *models.ServiceFile) error {
	filePerms, err := strconv.ParseUint(file.Mode, 8, 32)
	if err != nil {
		filePerms = 0644
	}

	var wr io.Writer
	if output != "" {
		if force {
			os.Remove(output)
		}
		outFile, err := os.OpenFile(output, os.O_CREATE|os.O_RDWR, os.FileMode(filePerms))
		if err != nil {
			logrus.Printf("Warning! Could not apply %s file permissions. Attempting to apply defaults %s", fileModeToRWXString(filePerms), fileModeToRWXString(uint64(0644)))
			outFile, err = os.OpenFile(output, os.O_CREATE|os.O_RDWR, 0644)
			if err != nil {
				return fmt.Errorf("Could not open %s for writing: %s", output, err.Error())
			}
		}
		defer outFile.Close()
		wr = outFile
	} else {
		logrus.Printf("Mode: %s\n\nContent:", fileModeToRWXString(filePerms))
		wr = os.Stdout
	}
	wr.Write([]byte(file.Contents))
	return nil
}
