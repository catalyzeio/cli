package certs

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/commands/ssl"
	"github.com/catalyzeio/cli/test"
)

const (
	certName = "example.com"
	pubKey   = `-----BEGIN CERTIFICATE-----
MIIDFDCCAfygAwIBAgIJAJ04dO4O6PrLMA0GCSqGSIb3DQEBBQUAMBAxDjAMBgNV
BAMTBWxvY2FsMB4XDTE3MDEyNjA2MjQwMloXDTI3MDEyNDA2MjQwMlowEDEOMAwG
A1UEAxMFbG9jYWwwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC7dwM2
rMj9N0mEZP+V9sWx0MKcuc4Uymv4BbJO/dP7ryXJEMqSZc7DrmUs1XTKEguWu9dL
0BylzCvaWalqKixojWL1Wojj6i8DfgHgFum+Fjd3EUZhbNnerfCC94Of1XCRSezG
sWP7V0gGSlxoptRvhH4NTHkyemnaZEDs323VtuhG0AgoQ8EWS/XeVAWLlSsHRPWp
BXjQn0ve33SsnbhbpkRkyB1jlH7vxbEaAX9aKrZYYSmXLz3NKp8ti8AljqybWC86
ymVl5qStd6yz/CrFiGWki0F46/BdPB8ZCY4iOsuMXbWWDiRuq7llu8iWEat651DO
VeAPKdQsRZgK/y1hAgMBAAGjcTBvMB0GA1UdDgQWBBRrj840X4a+uGDsKCRMHzX1
mtXAWTBABgNVHSMEOTA3gBRrj840X4a+uGDsKCRMHzX1mtXAWaEUpBIwEDEOMAwG
A1UEAxMFbG9jYWyCCQCdOHTuDuj6yzAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEB
BQUAA4IBAQAVpa/IkKyDPE7X4RsHZLsinEfJpAahrLsSBGDIo6cgpB3txntgmoLU
pC71ZQEE5glE4ENvflyLvvg6fAwlOVL0sax0GKfYgLJhg11CmsoRYiHCPh/bwqtU
iqAzjo7yCsyzo1Q0IMbc0RHFBmikHJEL6Dsuri1Skj+KnXLBibl8FeFuppgusV+W
8q3T/6ZNM8nFhRAPAQf7n4c4y+VjYuw/WSEdByH2NuLnLivb97E5BC0nr3/AK0Kz
MSsO3RiSxj07Gepc+Ce0VNXZkVAjUiwHvZeC7ebLC/SQs8ihogi/TVELgQkksgC/
lyUFiVHqjeeIKxYNy3d7RqGxzKKDRssi
-----END CERTIFICATE-----`
	privKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAu3cDNqzI/TdJhGT/lfbFsdDCnLnOFMpr+AWyTv3T+68lyRDK
kmXOw65lLNV0yhILlrvXS9Acpcwr2lmpaiosaI1i9VqI4+ovA34B4BbpvhY3dxFG
YWzZ3q3wgveDn9VwkUnsxrFj+1dIBkpcaKbUb4R+DUx5Mnpp2mRA7N9t1bboRtAI
KEPBFkv13lQFi5UrB0T1qQV40J9L3t90rJ24W6ZEZMgdY5R+78WxGgF/Wiq2WGEp
ly89zSqfLYvAJY6sm1gvOsplZeakrXess/wqxYhlpItBeOvwXTwfGQmOIjrLjF21
lg4kbqu5ZbvIlhGreudQzlXgDynULEWYCv8tYQIDAQABAoIBAQCN1FHzGLCLmzuc
1gjkvan+iPHkP1MiOa+MG0s3JiUugum0gGayciIHvDbBv9E3XIW2CfGuYwp5icoX
zcQ2FSg6BdY7yL5OqQveuYPTtaIsdYSLKd+0r/T522FexMKpt4MN+P8RqH37V6Kf
V70oVCffIz928kezoBfb6gOQ8s2XZRn8VHF+RDuxlT2x+eintCj9J87ynUYgcKwp
Pop2LkmRARqOCApAFCoIcywW7eV91JIXLvkmxn9J2Y/y0hdBXZyRGmPMnJEZjliI
nTanzs4RZENuOI37/zSXLn4R6M4MwRk/Lmi+Wsdd9gyOoLM+2hLtHVsdKPbHFlJ0
e5BzA35RAoGBAORU5DrKLUKQQDLcnhq70naph1taJFb7hTuX6stUaekRjZyVw1cd
2neOv3Z12h/vHOEnKcJhUVRTTCp5nn2FJVk2aSyeIUQ6LKajxlTmTe39gLubEDFr
jRmtX4WGJ+noKHz4gotZbL1Yn88PlIzBWYc6+BJOlqlckerCAza0EMKtAoGBANIu
Y35BZe2Y9BJ7BMzkxR9ZKs5ddXFAiSoT0TAI44UxX0M7R5/VMMNc4z3LmtUCHqOS
RCDdcjMunj5yiqaM1CEQ9Ol+YJ79IKtt2i0den5vDvHuHde0dNiAmpLJlFsazaIR
Zc8cLDvPiaNsb4mxM3Jq4SHfUebUemGl9FnsJOAFAoGALXLsXvthWO+Hp9gcLGwY
b4A9LiTaOOol0f/iP4jU8AyLaJCy6kNJ+iRS3gyFV3fsArEd8dAXNTbDYW0F7Cw1
i/V1p+jt7Du8KYtN7hZNisK7/hvWdE/ZLTRCYDyc80U/0ehRa9Vn/KSIYtnSEtZl
sLI/ML2t5ZZEgTsPErNy5p0CgYEAzVG9pbeTL9CsFWWRYerVWfNMKr4HnTOzCqTD
RE5anGGHsvC03kFv2ljiMBq2zQC+F4IqBYTuK2uN8GkKYvrNuuOKrJHlJ0sVYAH3
EP1sDRjGm7XF91L0lg7DcUN0Jq9/U6P1NaZK2764sSmbqAGvxUT9Wo6CvqCwULXC
hxl1SFUCgYEAyl+2eRiFXW6Opi3yLWSJ1FyqgZnqV9AUSXFTu3HFkw4yLzIwuq9M
nfOBIcrGX2exIylqMoeLxl9WfKbvZTQbL4zCzHoOtsuSTZErZywIIH0Jl5YZJnaT
EZ/6B0fi6DsLHY1tkIEvqgGI0kQX6IE84iZSi/Ubh8gQGwtutoZ1Stk=
-----END RSA PRIVATE KEY-----`
	pubKeyPath  = "example.pem"
	privKeyPath = "example-key.pem"
	invalidPath = "invalid-file.pem"
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

func TestMain(m *testing.M) {
	flag.Parse()
	if err := createCertFiles(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	statusCode := m.Run()
	cleanupCertFiles()
	os.Exit(statusCode)
}

var certCreateTests = []struct {
	hostname    string
	pubKeyPath  string
	privKeyPath string
	selfSigned  bool
	resolve     bool
	expectErr   bool
}{
	{certName, pubKeyPath, privKeyPath, true, true, false},
	{certName, pubKeyPath, privKeyPath, true, false, false},
	{certName, pubKeyPath, privKeyPath, false, true, false},
	{certName, pubKeyPath, invalidPath, true, true, true},
	{certName, invalidPath, privKeyPath, true, true, true},
	{"/?%", pubKeyPath, privKeyPath, true, true, true},
}

func TestCertsCreate(t *testing.T) {
	setup()
	defer teardown()
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services/"+test.SvcID+"/certs",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			fmt.Fprint(w, `{}`)
		},
	)
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","label":"service_proxy"}]`, test.SvcID))
		},
	)

	for _, data := range certCreateTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdCreate(data.hostname, data.pubKeyPath, data.privKeyPath, data.selfSigned, data.resolve, New(settings), services.New(settings), ssl.New(settings))

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}

func TestCertsCreateFailSSL(t *testing.T) {
	setup()
	defer teardown()
	settings := test.GetSettings(baseURL.String())
	mux.HandleFunc("/environments/"+test.EnvID+"/services",
		func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, fmt.Sprintf(`[{"id":"%s","label":"service_proxy"}]`, test.SvcID))
		},
	)

	// test
	err := CmdCreate(certName, pubKeyPath, privKeyPath, false, false, New(settings), services.New(settings), ssl.New(settings))

	// assert
	if err == nil {
		t.Fatalf("Expected error but found nil")
	}
}

func createCertFiles() error {
	cert, err := os.OpenFile(pubKeyPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer cert.Close()
	cert.WriteString(pubKey)
	key, err := os.OpenFile(privKeyPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer key.Close()
	key.WriteString(privKey)
	return nil
}

func cleanupCertFiles() error {
	err := os.Remove(pubKeyPath)
	if err != nil {
		err = os.Remove(privKeyPath)
	}
	return err
}
