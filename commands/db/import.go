package db

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/crypto"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/lib/transfer"
	"github.com/catalyzeio/cli/models"
)

func CmdImport(databaseName, filePath, mongoCollection, mongoDatabase string, skipBackup bool, id IDb, ip prompts.IPrompts, is services.IServices, ij jobs.IJobs) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("A file does not exist at path '%s'", filePath)
	}
	service, err := is.RetrieveByLabel(databaseName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services\" command.", databaseName)
	}
	if !skipBackup {
		logrus.Printf("Backing up \"%s\" before performing the import", databaseName)
		job, err := id.Backup(service)
		if err != nil {
			return err
		}
		logrus.Printf("Backup started (job ID = %s)", job.ID)

		// all because logrus treats print, println, and printf the same
		logrus.StandardLogger().Out.Write([]byte("Polling until backup finishes."))
		status, err := ij.PollTillFinished(job.ID, service.ID)
		if err != nil {
			return err
		}
		job.Status = status
		logrus.Printf("\nEnded in status '%s'", job.Status)
		err = id.DumpLogs("backup", job, service)
		if err != nil {
			return err
		}
		if job.Status != "finished" {
			return fmt.Errorf("Job finished with invalid status %s", job.Status)
		}
	} else {
		err := ip.YesNo("Are you sure you want to import data into your database without backing it up first? (y/n) ")
		if err != nil {
			return err
		}
	}
	logrus.Printf("Importing '%s' into %s (ID = %s)", filePath, databaseName, service.ID)
	job, err := id.Import(filePath, mongoCollection, mongoDatabase, service)
	if err != nil {
		return err
	}
	// all because logrus treats print, println, and printf the same
	logrus.StandardLogger().Out.Write([]byte(fmt.Sprintf("Processing import (job ID = %s).", job.ID)))

	status, err := ij.PollTillFinished(job.ID, service.ID)
	if err != nil {
		return err
	}
	job.Status = status
	logrus.Printf("\nImport complete (end status = '%s')", job.Status)
	err = id.DumpLogs("restore", job, service)
	if err != nil {
		return err
	}
	if job.Status != "finished" {
		return fmt.Errorf("Finished with invalid status %s", job.Status)
	}
	return nil
}

// Import imports data into a database service. The import is accomplished
// by encrypting the file locally, requesting a location that it can be uploaded
// to, then uploads the file. Once uploaded an automated service processes the
// file and acts according to the given parameters.
//
// The type of file that should be imported depends on the database. For
// PostgreSQL and MySQL, this should be a single `.sql` file. For Mongo, this
// should be a single tar'ed, gzipped archive (`.tar.gz`) of the database dump
// that you want to import.
func (d *SDb) Import(filePath, mongoCollection, mongoDatabase string, service *models.Service) (*models.Job, error) {
	key := make([]byte, crypto.KeySize)
	iv := make([]byte, crypto.IVSize)
	rand.Read(key)
	rand.Read(iv)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	encryptFileReader, err := d.Crypto.NewEncryptReader(file, key, iv)
	if err != nil {
		return nil, err
	}
	rt := transfer.NewReaderTransfer(encryptFileReader, encryptFileReader.CalculateTotalSize(int(fi.Size())))
	options := map[string]string{}
	if mongoCollection != "" {
		options["databaseCollection"] = mongoCollection
	}
	if mongoDatabase != "" {
		options["database"] = mongoDatabase
	}
	tmpAuth, err := d.TempUploadAuth(service)
	if err != nil {
		return nil, err
	}
	defer d.revokeAuth(service, tmpAuth)
	// Make sure to revoke temp auth even with an interrupt.
	c := make(chan os.Signal, 1)
	// "done" is used below to cancel printing the download status
	done := make(chan bool)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		done <- false
		rt.KillTransfer()
		d.revokeAuth(service, tmpAuth)
		os.Exit(1)
	}()
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1"), Credentials: credentials.NewStaticCredentials(tmpAuth.AccessKeyID, tmpAuth.SecretAccessKey, tmpAuth.SessionToken)})
	if err != nil {
		done <- false
		return nil, err
	}
	uploader := s3manager.NewUploader(sess)

	go printTransferStatus(false, rt, done)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(tmpAuth.Bucket),
		Key:                  aws.String(tmpAuth.FileName),
		Body:                 rt,
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		done <- false
		return nil, err
	}
	done <- true
	importParams := map[string]interface{}{}
	for key, value := range options {
		importParams[key] = value
	}
	importParams["filename"] = tmpAuth.FileName
	importParams["encryptionKey"] = string(d.Crypto.Hex(key, crypto.KeySize*2))
	importParams["encryptionIV"] = string(d.Crypto.Hex(iv, crypto.IVSize*2))
	importParams["dropDatabase"] = false

	b, err := json.Marshal(importParams)
	if err != nil {
		return nil, err
	}
	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod, d.Settings.UsersID)
	resp, statusCode, err := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/import", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, service.ID), headers)
	if err != nil {
		return nil, err
	}
	var job models.Job
	err = httpclient.ConvertResp(resp, statusCode, &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (d *SDb) revokeAuth(service *models.Service, tmpAuth *models.TempAuth) {
	if err := d.RevokeTempUploadAuth(service, tmpAuth.UserID); err != nil {
		logrus.Printf("Failed to cleanup after uploading your encrypted file: %s", err)
	}
}

func (d *SDb) TempUploadAuth(service *models.Service) (*models.TempAuth, error) {
	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod, d.Settings.UsersID)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/temp-auth?action_type=%s", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, service.ID, "upload"), headers)
	if err != nil {
		return nil, err
	}
	var tempAuth models.TempAuth
	err = httpclient.ConvertResp(resp, statusCode, &tempAuth)
	if err != nil {
		return nil, err
	}
	return &tempAuth, nil
}

func (d *SDb) RevokeTempUploadAuth(service *models.Service, userID string) error {
	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod, d.Settings.UsersID)
	resp, statusCode, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/revoke-temp-auth?user_id=%s", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, service.ID, url.QueryEscape(userID)), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
