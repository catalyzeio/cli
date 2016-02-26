package files

/*const (
	validSvcID   = "1"
	invalidSvcID = "2"
)

type MFiles models.ServiceFile

// To trigger an error, specify the invalidSvcID
func (m *MFiles) Create(svcID, filePath, name, mode string) (*models.ServiceFile, error) {
	if svcID == invalidSvcID {
		return nil, errors.New("There was an error")
	}
	return &models.ServiceFile{
		ID:             1,
		Contents:       "echo 'test'",
		GID:            0,
		Mode:           mode,
		Name:           name,
		UID:            0,
		EnableDownload: true,
	}, nil
}

// To trigger an error, specify the invalidSvcID
func (m *MFiles) List(svcID string) (*[]models.ServiceFile, error) {
	if svcID == invalidSvcID {
		return nil, errors.New("There was an error")
	}
	svcFile, _ := m.Create(svcID, "/tmp/file", "/my/file", "0644")
	return &[]models.ServiceFile{
		*svcFile,
	}, nil
}

// To trigger an error, specify the invalidSvcID
func (m *MFiles) Retrieve(fileName string, svcID string) (*models.ServiceFile, error) {
	if svcID == invalidSvcID {
		return nil, errors.New("There was an error")
	}
	svcFile, _ := m.Create(svcID, "/tmp/file", "/my/file", "0644")
	return svcFile, nil
}

func (m *MFiles) Rm(fileID int, svcID string) error {
	if svcID == invalidSvcID {
		return errors.New("There was an error")
	}
	return nil
}

// To trigger an error, specify an output that is not an empty string
func (m *MFiles) Save(output string, force bool, file *models.ServiceFile) error {
	if output != "" {
		return errors.New("There was an error")
	}
	return nil
}

var downloadTests = []struct {
	svcName   string
	fileName  string
	output    string
	force     bool
	expectErr bool
}{
	{"", "", "", false, false},
}

func TestDownload(t *testing.T) {
	for _, data := range downloadTests {
		t.Logf("%+v\n", data)
		err := CmdDownload(data.svcName, data.fileName, data.output, data.force, &MFiles{}, &MServices{})
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s\n", err.Error())
		}
	}
}*/
