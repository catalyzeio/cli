package test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/catalyzeio/cli/models"
	"github.com/stretchr/testify/mock"
)

// SetUpGitRepo runs git init in the current directory.
func SetUpGitRepo() error {
	output, err := RunCommand("git", []string{"init"})
	if err != nil {
		return fmt.Errorf("Unexpected error setting up git repo: %s", output)
	}
	return nil
}

// SetUpAssociation runs the associate command with the appropriate arguments to
// successfully associate to the test environment.
func SetUpAssociation() error {
	output, err := RunCommand(BinaryName, []string{"associate", EnvName, SvcLabel, "-a", Alias})
	if err != nil {
		return fmt.Errorf("Unexpected error setting up association: %s", output)
	}
	return nil
}

// ClearAssociations runs the clear --environments command.
func ClearAssociations() error {
	output, err := RunCommand(BinaryName, []string{"clear", "--environments"})
	if err != nil {
		return fmt.Errorf("Unexpected error clearing associations: %s", output)
	}
	return nil
}

// RunCommand runs the given command and arguments with the current os ENV.
func RunCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = strings.NewReader("n\n")
	output, err := cmd.CombinedOutput()
	return string(output), err
}

type MockHTTPManager struct {
	mock.Mock
	TLSManager models.HTTPManager // used for passing through mocks to the real method invocation
}

// GetHeaders mocked
func (m *MockHTTPManager) GetHeaders(sessionToken, version, pod, userID string) map[string][]string {
	m.Called(sessionToken, version, pod, userID)
	return map[string][]string{}
}

// ConvertResp mocked
func (m *MockHTTPManager) ConvertResp(b []byte, statusCode int, s interface{}) error {
	return m.TLSManager.ConvertResp(b, statusCode, s)
}

// Get mocked
func (m *MockHTTPManager) Get(body []byte, url string, headers map[string][]string) ([]byte, int, error) {
	args := m.Called(body, url, headers)
	return args.Get(0).([]byte), args.Int(1), args.Error(2)
}

// Post mocked
func (m *MockHTTPManager) Post(body []byte, url string, headers map[string][]string) ([]byte, int, error) {
	args := m.Called(body, url, headers)
	return args.Get(0).([]byte), args.Int(1), args.Error(2)
}

// PostFile mocked
func (m *MockHTTPManager) PostFile(filepath string, url string, headers map[string][]string) ([]byte, int, error) {
	args := m.Called(filepath, url, headers)
	return args.Get(0).([]byte), args.Int(1), args.Error(2)
}

// PutFile mocked
func (m *MockHTTPManager) PutFile(filepath string, url string, headers map[string][]string) ([]byte, int, error) {
	args := m.Called(filepath, url, headers)
	return args.Get(0).([]byte), args.Int(1), args.Error(2)
}

// Put mocked
func (m *MockHTTPManager) Put(body []byte, url string, headers map[string][]string) ([]byte, int, error) {
	args := m.Called(body, url, headers)
	return args.Get(0).([]byte), args.Int(1), args.Error(2)
}

// Delete mocked
func (m *MockHTTPManager) Delete(body []byte, url string, headers map[string][]string) ([]byte, int, error) {
	args := m.Called(body, url, headers)
	return args.Get(0).([]byte), args.Int(1), args.Error(2)
}

/*type MockEnvironments struct {
	Fail bool
}

func (m *MockEnvironments) List() (*[]models.Environment, map[string]error) {
	if m.Fail {
		return &[]models.Environment{}, map[string]error{
			Pod: errors.New("Failed to list environments"),
		}
	}
	return &[]models.Environment{
		models.Environment{
			ID:        EnvID,
			Name:      EnvName,
			Pod:       Pod,
			Namespace: Namespace,
			OrgID:     OrgID,
		},
	}, nil
}

func (m *MockEnvironments) Retrieve(envID string) (*models.Environment, error) {
	if m.Fail {
		return nil, fmt.Errorf("Failed to retrieve environment by ID %s", envID)
	}
	return &models.Environment{
		ID:        EnvID,
		Name:      EnvName,
		Pod:       Pod,
		Namespace: Namespace,
		OrgID:     OrgID,
	}, nil
}

func (m *MockEnvironments) Update(envID string, updates map[string]string) error {
	if m.Fail {
		return fmt.Errorf("Failed to update environment by ID %s", envID)
	}
	return nil
}

var DefaultService = models.Service{
	ID:         SvcID,
	Identifier: "code-1234",
	DNS:        "code-1234.internal",
	Type:       "code",
	Label:      SvcLabel,
	Size: models.ServiceSize{
		RAM:      1,
		Storage:  0,
		Behavior: "good",
		Type:     "code",
		CPU:      1,
	},
	Name:           "code",
	EnvVars:        map[string]string{},
	Source:         "git@catalyze-git.com",
	LBIP:           "127.0.0.1",
	Scale:          1,
	WorkerScale:    1,
	ReleaseVersion: "0",
	Redeployable:   true,
}

type MockServices struct {
	Fail   bool
	Return *models.Service
}

func (m *MockServices) List() (*[]models.Service, error) {
	if m.Fail {
		return nil, errors.New("Failed to list services")
	}
	if m.Return != nil {
		return &[]models.Service{*m.Return}, nil
	}
	return &[]models.Service{DefaultService}, nil
}

func (m *MockServices) ListByEnvID(envID, podID string) (*[]models.Service, error) {
	if m.Fail {
		return nil, fmt.Errorf("Failed to list services by environment ID %s", envID)
	}
	if m.Return != nil {
		return &[]models.Service{*m.Return}, nil
	}
	return &[]models.Service{DefaultService}, nil
}

func (m *MockServices) RetrieveByLabel(label string) (*models.Service, error) {
	if m.Fail {
		return nil, fmt.Errorf("Failed to retrieve service by label %s", label)
	}
	if m.Return != nil {
		return m.Return, nil
	}
	return &DefaultService, nil
}

func (m *MockServices) Update(svcID string, updates map[string]string) error {
	if m.Fail {
		return fmt.Errorf("Failed to update service by ID %s", svcID)
	}
	return nil
}

type MockSSL struct {
	Fail bool
}

func (m *MockSSL) Verify(chainPath, privateKeyPath, hostname string, selfSigned bool) error {
	if m.Fail {
		return errors.New("Failed to verify SSL")
	}
	return nil
}

func (m *MockSSL) Resolve(chainPath string) ([]byte, error) {
	if m.Fail {
		return nil, fmt.Errorf("Failed to resolve chain at %s", chainPath)
	}
	return []byte{}, nil
}*/
