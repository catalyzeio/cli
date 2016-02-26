package db

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/lib/jobs"
	"github.com/catalyzeio/cli/models"
)

func CmdImport(databaseName, filePath, mongoCollection, mongoDatabase string, id IDb, is services.IServices, ij jobs.IJobs) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("A file does not exist at path '%s'", filePath)
	}
	service, err := is.RetrieveByLabel(databaseName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"", databaseName)
	}
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
	logrus.Printf("Importing '%s' into %s (ID = %s)", filePath, databaseName, service.ID)
	job, err = id.Import(filePath, mongoCollection, mongoDatabase, service)
	if err != nil {
		return err
	}
	// all because logrus treats print, println, and printf the same
	logrus.StandardLogger().Out.Write([]byte(fmt.Sprintf("Processing import (job ID = %s).", job.ID)))

	status, err = ij.PollTillFinished(job.ID, service.ID)
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
	key := make([]byte, 32)
	iv := make([]byte, aes.BlockSize)
	rand.Read(key)
	rand.Read(iv)
	logrus.Println("Encrypting...")
	encrFilePath, err := d.Crypto.EncryptFile(filePath, key, iv)
	if err != nil {
		return nil, err
	}
	defer os.Remove(encrFilePath)
	options := map[string]string{}
	if mongoCollection != "" {
		options["mongoCollection"] = mongoCollection
	}
	if mongoDatabase != "" {
		options["mongoDatabase"] = mongoDatabase
	}
	logrus.Println("Uploading...")
	tempURL, err := d.TempUploadURL(service)
	if err != nil {
		return nil, err
	}

	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod)
	resp, statusCode, err := httpclient.PutFile(encrFilePath, tempURL.URL, headers)
	if err != nil {
		return nil, err
	}
	err = httpclient.ConvertResp(resp, statusCode, nil)
	if err != nil {
		return nil, err
	}
	importParams := map[string]interface{}{}
	for key, value := range options {
		importParams[key] = value
	}
	importParams["filename"] = tempURL.URL
	importParams["encryptionKey"] = string(d.Crypto.Base64Encode(d.Crypto.Hex(key)))
	importParams["encryptionIV"] = string(d.Crypto.Base64Encode(d.Crypto.Hex(iv)))
	importParams["dropDatabase"] = false

	b, err := json.Marshal(importParams)
	if err != nil {
		return nil, err
	}
	resp, statusCode, err = httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/import", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, service.ID), headers)
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

func (d *SDb) TempUploadURL(service *models.Service) (*models.TempURL, error) {
	headers := httpclient.GetHeaders(d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod)
	resp, statusCode, err := httpclient.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/restore-url", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, service.ID), headers)
	if err != nil {
		return nil, err
	}
	var tempURL models.TempURL
	err = httpclient.ConvertResp(resp, statusCode, &tempURL)
	if err != nil {
		return nil, err
	}
	return &tempURL, nil
}
