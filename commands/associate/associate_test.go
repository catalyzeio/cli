package associate

import (
	"errors"
	"reflect"
	"testing"

	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/test"
)

var associateTests = []struct {
	envName   string
	alias     string
	expectErr bool
}{
	{test.EnvName, test.Alias, false},
	{test.EnvName, "", false},
	{"", test.Alias, true},
	{"", "", false},
	{"bad-env", test.Alias, true},
}

func TestAssociate(t *testing.T) {
	mockHTTPManager := new(test.MockHTTPManager)
	mockHTTPManager.TLSManager = httpclient.NewTLSHTTPManager(false)
	settings := &models.Settings{
		HTTPManager: mockHTTPManager,
		Pods: &[]models.Pod{
			models.Pod{
				Name: test.Pod,
			},
		},
	}
	for _, data := range associateTests {
		t.Logf("Data: %+v", data)

		// reset
		settings.Environments = map[string]models.AssociatedEnvV2{}

		// expectations
		//var envs []models.Environment
		body := []byte("[{\"name\": \"" + test.EnvName + "\",\"id\": \"" + test.EnvID + "\",\"namespace\":\"" + test.Namespace + "\",\"organizationId\":\"" + test.OrgID + "\"}]")
		mockHTTPManager.On("GetHeaders", settings.SessionToken, settings.Version, test.Pod, settings.UsersID).Return(map[string][]string{})
		mockHTTPManager.On("Get", []byte(nil), "/environments", map[string][]string{}).Return(body, 200, nil)
		//mockHTTPManager.On("ConvertResp", body, 200, &envs).Return(nil)

		// test
		err := CmdAssociate(data.envName, data.alias, New(settings), environments.New(settings))

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
		expectedEnvs := map[string]models.AssociatedEnvV2{}
		if !data.expectErr {
			name := data.alias
			if name == "" {
				name = data.envName
				if name == "" {
					name = test.EnvName
				}
			}
			expectedEnvs = map[string]models.AssociatedEnvV2{
				name: models.AssociatedEnvV2{
					Name:          test.EnvName,
					EnvironmentID: test.EnvID,
					OrgID:         test.OrgID,
					Pod:           test.Pod,
				},
			}
		}
		if !reflect.DeepEqual(expectedEnvs, settings.Environments) {
			t.Errorf("Associated environment not added to settings object correctly. Expected: %+v. Found: %+v", expectedEnvs, settings.Environments)
		}
	}
}

func TestAssociateWithPodErrors(t *testing.T) {
	mockHTTPManager := new(test.MockHTTPManager)
	mockHTTPManager.TLSManager = httpclient.NewTLSHTTPManager(false)
	settings := &models.Settings{
		HTTPManager:  mockHTTPManager,
		Environments: map[string]models.AssociatedEnvV2{},
		Pods: &[]models.Pod{
			models.Pod{
				Name: test.Pod,
			},
		},
	}

	// expectations
	body := []byte("{\"message\": \"error\",\"id\": \"1\"}")
	mockHTTPManager.On("GetHeaders", settings.SessionToken, settings.Version, test.Pod, settings.UsersID).Return(map[string][]string{})
	mockHTTPManager.On("Get", []byte(nil), "/environments", map[string][]string{}).Return(body, 400, errors.New("error"))

	// test
	err := CmdAssociate("", "", New(settings), environments.New(settings))

	// assert
	if err == nil {
		t.Fatalf("Expected error but no error returned")
	}
	expectedEnvs := map[string]models.AssociatedEnvV2{}
	if !reflect.DeepEqual(expectedEnvs, settings.Environments) {
		t.Errorf("Associated environment not added to settings object correctly. Expected: %+v. Found: %+v", expectedEnvs, settings.Environments)
	}
}
