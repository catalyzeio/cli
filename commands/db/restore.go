package db

// RestoreBackup a database service from an existing backup
/*func RestoreBackup(svcName, backupID string, skipPoll bool, id IDb, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\"\n", svcName)
	}
	task, err := id.Restore(backupID, service)
	if err != nil {
		return err
	}
	fmt.Printf("Restoring (task ID = %s)\n", task.ID)
	if !skipPoll {
		fmt.Print("Polling until restore finishes.")
		ch := make(chan string, 1)
		go helpers.PollTaskStatus(task.ID, ch, settings)
		status := <-ch
		task.Status = status
		fmt.Printf("\nEnded in status '%s'\n", task.Status)
		helpers.DumpLogs(service, task, "restore", settings)
		if task.Status != "finished" {
			return fmt.Errorf("Task finished with invalid status %s\n", task.Status)
		}
	}
}

func (d *SDb) Restore(backupID string, service *models.Service) (*models.Task, error) {
	backup := map[string]string{
		"archiveType":    "cf",
		"encryptionType": "aes",
	}
	b, err := json.Marshal(backup)
	if err != nil {
		return nil, err
	}
	headers := httpclient.GetHeaders(d.Settings.APIKey, d.Settings.SessionToken, d.Settings.Version, d.Settings.Pod)
	resp, statusCode, nil := httpclient.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/restore/%s", d.Settings.PaasHost, d.Settings.PaasHostVersion, d.Settings.EnvironmentID, service.ID, backupID), headers)
	if err != nil {
		return nil, err
	}
	var m map[string]string
	err = httpclient.ConvertResp(resp, statusCode, &m)
	if err != nil {
		return nil, err
	}
	return &models.Task{
		ID: m["taskId"],
	}, nil
}*/
