package associated

import (
	"testing"

	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/test"
)

func TestAssociated(t *testing.T) {
	settings := &models.Settings{
		Environments: map[string]models.AssociatedEnvV2{
			test.Alias: models.AssociatedEnvV2{
				Name:          test.EnvName,
				EnvironmentID: test.EnvID,
				Pod:           test.Pod,
				OrgID:         test.OrgID,
			},
		},
	}
	err := CmdAssociated(New(settings))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func TestAssociatedNoAssociations(t *testing.T) {
	settings := &models.Settings{
		Environments: map[string]models.AssociatedEnvV2{},
	}
	err := CmdAssociated(New(settings))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}
