package services

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
)

func CmdRename(svcName, label string, is IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}
	data := map[string]string{}
	data["label"] = label
	err = is.Update(service.ID, data)
	if err != nil {
		return err
	}
	logrus.Printf("Successfully renamed your service to %s", label)
	return nil
}

func (s *SServices) Update(svcID string, updates map[string]string) error {
	b, err := json.Marshal(updates)
	if err != nil {
		return err
	}
	headers := s.Settings.HTTPManager.GetHeaders(s.Settings.SessionToken, s.Settings.Version, s.Settings.Pod, s.Settings.UsersID)
	resp, statusCode, err := s.Settings.HTTPManager.Put(b, fmt.Sprintf("%s%s/environments/%s/services/%s", s.Settings.PaasHost, s.Settings.PaasHostVersion, s.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return err
	}
	return s.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
