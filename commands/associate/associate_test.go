package associate

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/catalyzeio/cli/commands/environments"
	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/test"
)

var (
	mux     *http.ServeMux
	server  *httptest.Server
	baseURL *url.URL
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	baseURL, _ = url.Parse(server.URL)
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

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
	setup()
	defer teardown()
	settings := test.GetSettings(baseURL.String())

	mux.HandleFunc("/environments",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			if r.Header.Get("X-Pod-ID") == test.Pod {
				fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","name":"%s","namespace":"%s","organizationId":"%s"}]`, test.EnvID, test.EnvName, test.Namespace, test.OrgID))
			} else {
				fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","name":"%s","namespace":"%s","organizationId":"%s"}]`, test.EnvIDAlt, test.EnvNameAlt, test.NamespaceAlt, test.OrgIDAlt))
			}
		},
	)

	for _, data := range associateTests {
		t.Logf("Data: %+v", data)

		// reset
		settings.Environments = map[string]models.AssociatedEnvV2{}

		// test
		err := CmdAssociate(data.envName, data.alias, New(settings), environments.New(settings))

		// assertions
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
		expectedEnvs := map[string]models.AssociatedEnvV2{}
		if !data.expectErr {
			name := data.alias
			if name == "" {
				name = test.EnvName
			}
			expectedEnvs[name] = models.AssociatedEnvV2{
				Name:          test.EnvName,
				EnvironmentID: test.EnvID,
				OrgID:         test.OrgID,
				Pod:           test.Pod,
			}
			if data.envName == "" {
				expectedEnvs[test.EnvNameAlt] = models.AssociatedEnvV2{
					Name:          test.EnvNameAlt,
					EnvironmentID: test.EnvIDAlt,
					OrgID:         test.OrgIDAlt,
					Pod:           test.PodAlt,
				}
			}
		}
		if !reflect.DeepEqual(expectedEnvs, settings.Environments) {
			t.Errorf("Associated environment not added to settings object correctly. Expected: %+v. Found: %+v", expectedEnvs, settings.Environments)
		}
	}
}

func TestAssociateWithPodErrors(t *testing.T) {
	setup()
	defer teardown()
	settings := test.GetSettings(baseURL.String())
	settings.Environments = map[string]models.AssociatedEnvV2{}

	mux.HandleFunc("/environments",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			http.Error(w, `{"title":"Error","description":"error","code":1}`, 400)
		},
	)

	// test
	err := CmdAssociate("", "", New(settings), environments.New(settings))

	// assert
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	expectedEnvs := map[string]models.AssociatedEnvV2{}
	if !reflect.DeepEqual(expectedEnvs, settings.Environments) {
		t.Errorf("Associated environment not added to settings object correctly. Expected: %+v. Found: %+v", expectedEnvs, settings.Environments)
	}
}
