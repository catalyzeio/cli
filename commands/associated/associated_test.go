package associated

/*var associatedTests = []struct {
	envs       map[string]models.AssociatedEnv
	defaultEnv string
	expectErr  bool
}{
	{
		map[string]models.AssociatedEnv{
			"test": models.AssociatedEnv{
				EnvironmentID: "env1",
				ServiceID:     "svc1",
				Directory:     "~/test",
				Name:          "test",
				Pod:           "pod01",
				OrgID:         "1234",
			},
			"test2": models.AssociatedEnv{
				EnvironmentID: "env2",
				ServiceID:     "svc2",
				Directory:     "~/test2",
				Name:          "long-env-name",
				Pod:           "pod02",
				OrgID:         "5678",
			},
		},
		"test2",
		false,
	},
}

func TestAssociated(t *testing.T) {
	for _, data := range associatedTests {
		t.Logf("%+v\n", data)
		err := CmdAssociated(New(getSettings(data.envs)))
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s\n", err.Error())
		}
	}
}

func getSettings(associatedEnvs map[string]models.AssociatedEnv) *models.Settings {
	return &models.Settings{
		Environments: associatedEnvs,
		Pods:         &[]models.Pod{},
	}
}*/
