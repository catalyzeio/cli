package environments

/*const (
	validPod   = "pod01"
	emptyPod   = "pod02"
	invalidPod = "pod03"
)

type MEnvironments models.Environment

func (m *MEnvironments) List() (*[]models.Environment, error) {
	switch m.Pod {
	case validPod:
		return &[]models.Environment{
			models.Environment{
				ID:        m.ID,
				Name:      m.Name,
				Pod:       m.Pod,
				Namespace: m.Namespace,
				OrgID:     m.OrgID,
			},
		}, nil
	case emptyPod:
		return &[]models.Environment{}, nil
	default:
		return nil, errors.New("there was an error")
	}
}

func (m *MEnvironments) Retrieve(envID string) (*models.Environment, error) {
	envs, err := m.List()
	if err == nil && len(*envs) > 0 {
		return &(*envs)[0], nil
	}
	return nil, errors.New("404 not found")
}

var environmentsTests = []struct {
	pod       string
	expectErr bool
}{
	{validPod, false},
	{emptyPod, false},
	{invalidPod, true},
}

func TestEnvironments(t *testing.T) {
	for _, data := range environmentsTests {
		t.Logf("%+v\n", data)
		mockEnvs := MEnvironments{
			Pod: data.pod,
		}
		err := CmdEnvironments(&mockEnvs)
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s\n", err.Error())
		}
	}
}

func getSettings(pod string) *models.Settings {
	return &models.Settings{
		Pod:          pod,
		Environments: make(map[string]models.AssociatedEnv, 0),
		Pods:         &[]models.Pod{},
	}
}*/
