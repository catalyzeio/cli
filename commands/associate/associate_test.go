package associate

import (
	"reflect"
	"testing"

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
	settings := &models.Settings{}
	for _, data := range associateTests {
		t.Logf("Data: %+v", data)

		// reset
		settings.Environments = map[string]models.AssociatedEnvV2{}

		// test
		err := CmdAssociate(data.envName, data.alias, New(settings), &test.MockEnvironments{})

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
	settings := &models.Settings{
		Environments: map[string]models.AssociatedEnvV2{},
	}

	// test
	err := CmdAssociate("", "", New(settings), &test.MockEnvironments{Fail: true})

	// assert
	if err == nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	expectedEnvs := map[string]models.AssociatedEnvV2{}
	if !reflect.DeepEqual(expectedEnvs, settings.Environments) {
		t.Errorf("Associated environment not added to settings object correctly. Expected: %+v. Found: %+v", expectedEnvs, settings.Environments)
	}
}
