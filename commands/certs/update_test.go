package certs

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/commands/ssl"
	"github.com/catalyzeio/cli/test"
)

const (
	certsUpdateCommandName    = "certs"
	certsUpdateSubcommandName = "update"
	certsUpdateStandardOutput = `Updated 'example.com'
To make your updated cert go live, you must redeploy your service proxy with the "catalyze redeploy service_proxy" command
`
)

var certUpdateTests = []struct {
	hostname    string
	pubKeyPath  string
	privKeyPath string
	selfSigned  bool
	resolve     bool
	expectErr   bool
}{
	{certName, pubKeyPath, privKeyPath, true, false, false},
	{certName, invalidPath, privKeyPath, true, false, true}, // invalid cert path
	{certName, pubKeyPath, invalidPath, true, false, true},  // invalid key path
	{certName, pubKeyPath, privKeyPath, false, false, true}, // cert not signed by CA
	{certName, pubKeyPath, privKeyPath, true, true, false},
	{"bad-cert-name", pubKeyPath, privKeyPath, true, false, true},
}

func TestCertsUpdate(t *testing.T) {
	setup()
	defer teardown()
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/certs/"+certName,
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "PUT")
			fmt.Fprint(w, fmt.Sprintf(`{"name":"%s"}`, certName))
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","label":"service_proxy"}]`, test.SvcID))
		},
	)

	for _, data := range certUpdateTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdUpdate(data.hostname, data.pubKeyPath, data.privKeyPath, data.selfSigned, data.resolve, New(settings), services.New(settings), ssl.New(settings))

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}
