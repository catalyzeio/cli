package associate

/*const validEnvName = "env"
const invalidEnvName = "badEnv"

const validAppName = "app01"
const invalidAppName = "badApp"

type MEnvironments models.Environment

func (m *MEnvironments) List() (*[]models.Environment, error) {
	return &[]models.Environment{
		models.Environment{
			ID:        m.ID,
			Name:      m.Name,
			Pod:       m.Pod,
			Namespace: m.Namespace,
			OrgID:     m.OrgID,
		},
	}, nil
}

func (m *MEnvironments) Retrieve(envID string) (*models.Environment, error) {
	envs, _ := m.List()
	return &(*envs)[0], nil
}

type MServices models.Service

func (m *MServices) List() (*[]models.Service, error) {
	return &[]models.Service{
		models.Service{
			ID:      m.ID,
			Type:    m.Type,
			Label:   m.Label,
			Size:    m.Size,
			Name:    m.Name,
			EnvVars: m.EnvVars,
			Source:  m.Source,
			LBIP:    m.LBIP,
		},
	}, nil
}

func (m *MServices) ListByEnvID(envID, podID string) (*[]models.Service, error) {
	return m.List()
}

func (m *MServices) Retrieve(svcID string) (*models.Service, error) {
	svcs, _ := m.List()
	svc := (*svcs)[0]
	svc.ID = svcID
	return &svc, nil
}

func (m *MServices) RetrieveByLabel(label string) (*models.Service, error) {
	svcs, _ := m.List()
	svc := (*svcs)[0]
	svc.Label = label
	return &svc, nil
}

var associateTests = []struct {
	envLabel      string
	svcLabel      string
	alias         string
	remote        string
	defaultEnv    bool
	createGitRepo bool
	expectErr     bool
}{
	{validEnvName, validAppName, "e", "", false, true, false},
	{validEnvName, validAppName, "e", "", true, true, false},
	{validEnvName, validAppName, "e", "ctlyz", false, true, false},
	{validEnvName, validAppName, "", "", false, true, false},
	{validEnvName, invalidAppName, "e", "", false, true, true},
	{invalidEnvName, validAppName, "e", "", false, true, true},
	{validEnvName, validAppName, "e", "", false, false, true},
}

func TestAssociate(t *testing.T) {
	for _, data := range associateTests {
		t.Logf("%+v\n", data)
		createGitRepo(data.envLabel, data.createGitRepo)

		settings := getSettings()
		mockEnvs := MEnvironments{
			Name:  data.envLabel,
			Pod:   settings.Pod,
			OrgID: settings.OrgID,
		}
		mockSvcs := MServices{
			Type:  "code",
			Label: data.svcLabel,
			Name:  "code",
		}
		err := CmdAssociate(data.envLabel, data.svcLabel, data.alias, data.remote, data.defaultEnv, New(settings), git.New(), &mockEnvs, &mockSvcs)
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s\n", err.Error())
		}
		name := data.alias
		if name == "" {
			name = data.envLabel
		}

		found := false
		for _, env := range settings.Environments {
			if env.Name == name {
				found = true
				break
			}
		}
		if !found {
			t.Error("Environment not added to the settings list of environments or the alias was not used")
		}
		if data.defaultEnv && settings.Default != name {
			t.Error("Default environment specified but was not stored in the settings")
		}

		expectedRemote := data.remote
		if expectedRemote == "" {
			expectedRemote = "catalyze"
		}
		remotes, err := git.New().List()
		if err != nil {
			t.Errorf("Error listing git remotes: %s\n", err.Error())
		}
		found = false
		for _, r := range remotes {
			if r == expectedRemote {
				found = true
			}
		}
		if !found {
			t.Errorf("Proper git remote not listed. Found '%+v' instead\n", remotes)
		}

		destroyGitRepo(data.envLabel)
	}
}

func getSettings() *models.Settings {
	return &models.Settings{
		AuthHost:        config.AuthHost,
		PaasHost:        config.PaasHost,
		AuthHostVersion: "",
		PaasHostVersion: "",
		Version:         "dev",
		Username:        "test",
		Password:        "test",
		EnvironmentID:   "1234",
		ServiceID:       "5678",
		Pod:             "pod01",
		EnvironmentName: validEnvName,
		OrgID:           "192837465",
		SessionToken:    "1234567890",
		UsersID:         "0987654321",
		Environments:    make(map[string]models.AssociatedEnv, 0),
		Default:         "",
		Pods:            &[]models.Pod{},
	}
}

func createGitRepo(name string, create bool) {

}

func destroyGitRepo(name string) {

}*/
