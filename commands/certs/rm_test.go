package certs

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/test"
)

var certRmTests = []struct {
	hostname  string
	expectErr bool
}{
	{certName, false},
	{"bad-cert-name", true},
}

func TestCertsRm(t *testing.T) {
	setup()
	defer teardown()
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/certs/"+certName,
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "DELETE")
			fmt.Fprint(w, "")
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","label":"service_proxy"}]`, test.SvcID))
		},
	)

	for _, data := range certRmTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdRm(data.hostname, New(settings), services.New(settings))

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}
