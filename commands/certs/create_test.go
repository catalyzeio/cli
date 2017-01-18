package certs

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/catalyzeio/cli/models"
	"github.com/catalyzeio/cli/test"
)

const (
	certName = "example.com"
	pubKey   = `-----BEGIN CERTIFICATE-----
MIICATCCAWoCCQCsoDP5n7FfzzANBgkqhkiG9w0BAQUFADBFMQswCQYDVQQGEwJB
VTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50ZXJuZXQgV2lkZ2l0
cyBQdHkgTHRkMB4XDTE1MDYwNDE5MTkzNVoXDTE2MDYwMzE5MTkzNVowRTELMAkG
A1UEBhMCQVUxEzARBgNVBAgTClNvbWUtU3RhdGUxITAfBgNVBAoTGEludGVybmV0
IFdpZGdpdHMgUHR5IEx0ZDCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkCgYEA3+Gz
NFJhBdbUcFUxzlm70DJHXa9+nOAHZ9S6c66T1FXBRF94GfTSq8Qg9U+EOZf5cuhN
6wkLD1LLHMdb/UEjyCVVOqscfeR/nPCT5B9sv881PM8jL8C7grAUezcKiNx7Fng8
Dj9sczwziBR9P5ke5TI1g62LhHc0KGgMa8oNY7UCAwEAATANBgkqhkiG9w0BAQUF
AAOBgQBgTk8C+e13xGEw8qI2xhNfudt+8ffzIjNNWptb8rhGWblyY7EVBuU24LqE
oIOS7EH2aRhgvZjPUEQCNl+foQBRnRkYBeBhfUTl8QAUQNIyRUAHlQcPct9+VYcz
7OeuMetZkluMG3w62ooiufaGC/8orztDEySO4cj1HWssE2h/zw==
-----END CERTIFICATE-----`
	privKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDf4bM0UmEF1tRwVTHOWbvQMkddr36c4Adn1LpzrpPUVcFEX3gZ
9NKrxCD1T4Q5l/ly6E3rCQsPUsscx1v9QSPIJVU6qxx95H+c8JPkH2y/zzU8zyMv
wLuCsBR7NwqI3HsWeDwOP2xzPDOIFH0/mR7lMjWDrYuEdzQoaAxryg1jtQIDAQAB
AoGAeXoVqoYobuqqSmlvpO+7oLQnVQYsRSKp4gTjRnGrdMMzIs5KdIsK5Hh/CZwj
urxjdZ3m6Wj2v1HFM9BYcYouxx5ZYbUWx4tXeQhoVjvu8GxU6uwkDl+kQMjqcvfV
dXEoIm7ejzcvialYlHnsO8HFiB3ayhoQOK3kGcY6dGISWwECQQD/7R8/EIPAP0lU
P97w+I7j2kG79PTvCzoXygqVrmjeW6RJ6FvzT30iCnr5PVmPzHReL+q3i6tMHpGi
eeo0T0atAkEA3/I22OTH2QrKmSaW3EoPNDq78hJzsbSoVaHz+6mMn2ZungzBhJ7i
dOkUzkzuZtftYIcCQ2MtGDeSNIXuohOaKQJADwbVNta5ZahRnejCJlPxz98YzPht
CTwXhR4P0QoUjjnDQ7Oo8nhQWJdU8R1xDMhsbLtThMNmo2mIE4ok/j1JYQJBAJKg
pqSwduF3HVvVVmV54CaUZkaDKlkqLiWTWopmYvpjOP4m3/YTibZ+fe7tlBKmQng3
LZYts3Ltv77ACpT4PLECQQDDql4xPUb6WfsSjyqqfwnzkFLWADTcQQG5MmUX6iNJ
FBlcbW65DK1xPIitnX+jf803WaMPAP5YBoH6jC6VgcVH
-----END RSA PRIVATE KEY-----`
)

var (
	pubKeyPath  = "example.pem"
	privKeyPath = "example-key.pem"
	invalidPath = "invalid-file.pem"
)

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
	settings := &models.Settings{
		Environments: map[string]models.AssociatedEnvV2{
			test.Alias: models.AssociatedEnvV2{
				Name:          test.EnvName,
				EnvironmentID: test.EnvID,
				Pod:           test.Pod,
				OrgID:         test.OrgID,
			},
		},
	}
	for _, data := range certCreateTests {
		t.Logf("Data: %+v", data)

		// test
		err := CmdCreate(data.hostname, data.pubKeyPath, data.privKeyPath, data.selfSigned, data.resolve, New(settings), &test.MockServices{}, &test.MockSSL{})

		// assert
		if err != nil != data.expectErr {
			t.Errorf("Unexpected error: %s", err)
			continue
		}
	}
}

func TestCertsCreateFailSSL(t *testing.T) {
	settings := &models.Settings{
		Environments: map[string]models.AssociatedEnvV2{
			test.Alias: models.AssociatedEnvV2{
				Name:          test.EnvName,
				EnvironmentID: test.EnvID,
				Pod:           test.Pod,
				OrgID:         test.OrgID,
			},
		},
	}

	// test
	err := CmdCreate(certName, pubKeyPath, privKeyPath, false, false, New(settings), &test.MockServices{}, &test.MockSSL{Fail: true})

	// assert
	if err != nil {
		// with resolve = false, no SSL code should be called
		t.Fatalf("Unexpected error: %s", err)
	}

	// test
	err = CmdCreate(certName, pubKeyPath, privKeyPath, false, true, New(settings), &test.MockServices{}, &test.MockSSL{Fail: true})

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
