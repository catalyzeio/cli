package certs

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/test"
)

func TestCertsList(t *testing.T) {
	setup()
	defer teardown()
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/certs",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `[{"name":"cert1"},{"name":"cert2"}]`)
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","label":"service_proxy"}]`, test.SvcID))
		},
	)

	// test
	err := CmdList(New(settings), services.New(settings))

	// assert
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}
