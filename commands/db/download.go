package db

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
	"github.com/tj/go-spin"
)

func CmdDownload(databaseName, backupID, filePath string, force bool, id IDb, ip prompts.IPrompts, is services.IServices) error {
	err := ip.PHI()
	if err != nil {
		return err
	}
	if !force {
		if _, err := os.Stat(filePath); err == nil {
			return fmt.Errorf("File already exists at path '%s'. Specify `--force` to overwrite", filePath)
		}
	} else {
		os.Remove(filePath)
	}
	service, err := is.RetrieveByLabel(databaseName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", databaseName)
	}
	err = id.Download(backupID, filePath, service)
	if err != nil {
		return err
	}
	logrus.Printf("%s backup downloaded successfully to %s", databaseName, filePath)
	logrus.Printf("You can also view logs for this backup with the \"catalyze db logs %s %s\" command", databaseName, backupID)
	return nil
}

// Download an existing backup to the local machine. The backup is encrypted
// throughout the entire journey and then decrypted once it is stored locally.
func (d *SDb) Download(backupID, filePath string, service *models.Service) error {
	job, err := d.Jobs.Retrieve(backupID, service.ID, false)
	if err != nil {
		return err
	}
	if job.Type != "backup" || (job.Status != "finished" && job.Status != "disappeared") {
		return errors.New("Only 'finished' 'backup' jobs may be downloaded")
	}
	spinner := spin.New()
	ticker := time.Tick(100 * time.Millisecond)
	done := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-ticker:
				fmt.Printf("\r\033[mDownloading backup %s. This may take awhile %s\033[m ", backupID, spinner.Next())
			case <-done:
				return
			}
		}
	}()
	tempURL, err := d.TempDownloadURL(backupID, service)
	if err != nil {
		done <- struct{}{}
		return err
	}
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		done <- struct{}{}
		return err
	}
	defer os.Remove(dir)
	tmpFile, err := ioutil.TempFile(dir, "")
	if err != nil {
		done <- struct{}{}
		return err
	}
	logrus.Debugf("Downloading to temporary file at %s", tmpFile.Name())

	u, _ := url.Parse(tempURL.URL)
	svc := s3.New(session.New(&aws.Config{Region: aws.String("us-east-1"), Credentials: credentials.AnonymousCredentials}))
	req, resp := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(strings.Split(u.Host, ".")[0]),
		Key:    aws.String(strings.TrimLeft(u.Path, "/")),
	})
	req.HTTPRequest.URL.RawQuery = u.RawQuery
	err = req.Send()
	if err != nil {
		done <- struct{}{}
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		done <- struct{}{}
		return err
	}
	tmpFile.Close()
	done <- struct{}{}

	logrus.Println("\nDecrypting...")
	err = d.Crypto.DecryptFile(tmpFile.Name(), job.Backup.Key, job.Backup.IV, filePath)
	if err != nil {
		return err
	}

	return nil
}

func (d *SDb) TempDownloadURL(jobID string, service *models.Service) (*models.TempURL, error) {
	headers := d.Settings.HTTPManager.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod, d.Settings.UsersID)
	resp, statusCode, err := d.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/backup-url/%s", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, service.ID, jobID), headers)
	if err != nil {
		return nil, err
	}
	var tempURL models.TempURL
	err = d.Settings.HTTPManager.ConvertResp(resp, statusCode, &tempURL)
	if err != nil {
		return nil, err
	}
	return &tempURL, nil
}
