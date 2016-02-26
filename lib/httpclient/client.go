package httpclient

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/updater"
	"github.com/catalyzeio/cli/models"
)

func getClient() *http.Client {
	var tr = &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		},
	}

	return &http.Client{
		Transport: tr,
	}
}

// GetHeaders builds a map of headers for a new request.
func GetHeaders(sessionToken, version, pod string) map[string][]string {
	b := make([]byte, 32)
	rand.Read(b)
	nonce := base64.StdEncoding.EncodeToString(b)
	timestamp := time.Now().Unix()
	return map[string][]string{
		"Accept":              {"application/json"},
		"Content-Type":        {"application/json"},
		"Authorization":       {fmt.Sprintf("Bearer %s", sessionToken)},
		"X-CLI-Version":       {version},
		"X-Pod-ID":            {pod},
		"X-Request-Nonce":     {nonce},
		"X-Request-Timestamp": {fmt.Sprintf("%d", timestamp)},
		"User-Agent":          {fmt.Sprintf("catalyze-cli-%s", version)},
	}
}

// ConvertResp takes in a resp from one of the httpclient methods and
// checks if it is a successful request. If not, it is parsed as an error object
// and returned as an error. Otherwise it will be marshalled into the requested
// interface. ALWAYS PASS A POINTER INTO THIS METHOD. If you don't pass a struct
// pointer your original object will be nil or an empty struct.
func ConvertResp(b []byte, statusCode int, s interface{}) error {
	logrus.Debugf("%d resp: %s", statusCode, string(b))
	if IsError(statusCode) {
		return ConvertError(b, statusCode)
	}
	if b == nil || len(b) == 0 || s == nil {
		return nil
	}
	return json.Unmarshal(b, s)
}

// IsError checks if an HTTP response code is outside of the "OK" range.
func IsError(statusCode int) bool {
	return statusCode < 200 || statusCode >= 300
}

// ConvertError attempts to convert a response into a usable error object.
func ConvertError(b []byte, statusCode int) error {
	msg := fmt.Sprintf("(%d)", statusCode)
	if b != nil || len(b) > 0 {
		var errs models.Error
		err := json.Unmarshal(b, &errs)
		if err == nil && errs.Title != "" && errs.Description != "" {
			msg = fmt.Sprintf("(%d) %s: %s", errs.Code, errs.Title, errs.Description)
		} else {
			msg = fmt.Sprintf("(%d) %s", statusCode, string(b))
		}
	}
	return errors.New(msg)
}

// Get performs a GET request
func Get(body []byte, url string, headers map[string][]string) ([]byte, int, error) {
	reader := bytes.NewReader(body)
	return MakeRequest("GET", url, reader, headers)
}

// Post performs a POST request
func Post(body []byte, url string, headers map[string][]string) ([]byte, int, error) {
	reader := bytes.NewReader(body)
	return MakeRequest("POST", url, reader, headers)
}

// PostFile uploads a file with a POST
func PostFile(filepath string, url string, headers map[string][]string) ([]byte, int, error) {
	return uploadFile("POST", filepath, url, headers)
}

// PutFile uploads a file with a PUT
func PutFile(filepath string, url string, headers map[string][]string) ([]byte, int, error) {
	return uploadFile("PUT", filepath, url, headers)
}

func uploadFile(method, filepath, url string, headers map[string][]string) ([]byte, int, error) {
	logrus.Debugf("%s %s", method, url)
	logrus.Debugf("%+v", headers)
	logrus.Debugf("%s", filepath)
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return nil, 0, err
	}
	info, _ := file.Stat()
	client := getClient()
	req, _ := http.NewRequest(method, url, file)
	req.ContentLength = info.Size()

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	return respBody, resp.StatusCode, nil
}

// Put performs a PUT request
func Put(body []byte, url string, headers map[string][]string) ([]byte, int, error) {
	reader := bytes.NewReader(body)
	return MakeRequest("PUT", url, reader, headers)
}

// Delete performs a DELETE request
func Delete(body []byte, url string, headers map[string][]string) ([]byte, int, error) {
	reader := bytes.NewReader(body)
	return MakeRequest("DELETE", url, reader, headers)
}

// MakeRequest is a generic HTTP runner that performs a request and returns
// the result body as a byte array. It's up to the caller to transform them
// into an object.
func MakeRequest(method string, url string, body io.Reader, headers map[string][]string) ([]byte, int, error) {
	logrus.Debugf("%s %s", method, url)
	logrus.Debugf("%+v", headers)
	logrus.Debugf("%s", body)
	client := getClient()
	req, _ := http.NewRequest(method, url, body)
	req.Header = headers

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 412 {
		updater.AutoUpdater.ForcedUpgrade()
		return nil, 0, fmt.Errorf("A required update has been applied. Please re-run this command.")
	}
	return respBody, resp.StatusCode, nil
}
