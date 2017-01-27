package console

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"golang.org/x/net/websocket"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/models"
	"github.com/docker/docker/pkg/term"
)

func CmdConsole(svcName, command string, ic IConsole, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.\n", svcName)
	}
	return ic.Open(command, service)
}

// Open opens a secure console to a code or database service. For code
// services, a command is required. This command is executed as root in the
// context of the application root directory. For database services, no command
// is needed - instead, the appropriate command for the database type is run.
// For example, for a postgres database, psql is run.
func (c *SConsole) Open(command string, service *models.Service) error {
	stdin, stdout, _ := term.StdStreams()
	fdIn, isTermIn := term.GetFdInfo(stdin)
	if !isTermIn {
		return errors.New("StdIn is not a terminal")
	}
	var size *term.Winsize
	var err error
	if runtime.GOOS != "windows" {
		size, err = term.GetWinsize(fdIn)
	} else {
		fdOut, _ := term.GetFdInfo(stdout)
		size, err = term.GetWinsize(fdOut)
	}

	if err != nil {
		return err
	}
	if size.Width != 80 {
		logrus.Warnln("Your terminal width is not 80 characters. Please resize your terminal to be exactly 80 characters wide to avoid line wrapping issues.")
	} else {
		logrus.Warnln("Keep your terminal width at 80 characters. Resizing your terminal will introduce line wrapping issues.")
	}

	logrus.Printf("Opening console to %s (%s)", service.Name, service.ID)
	job, err := c.Request(command, service)
	if err != nil {
		return err
	}
	// all because logrus treats print, println, and printf the same
	logrus.StandardLogger().Out.Write([]byte(fmt.Sprintf("Waiting for the console (job ID = %s) to be ready. This might take a minute.", job.ID)))

	validStatuses := []string{"running", "finished", "failed"}
	status, err := c.Jobs.PollForStatus(validStatuses, job.ID, service.ID)
	if err != nil {
		return err
	}
	found := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("\nCould not open a console connection. Entered state '%s'", status)
	}
	job.Status = status
	defer c.Destroy(job.ID, service)
	creds, err := c.RetrieveTokens(job.ID, service)
	if err != nil {
		return err
	}

	creds.URL = strings.Replace(creds.URL, "http", "ws", 1)
	logrus.Println("\nConnecting...")

	// BEGIN websocket impl
	config, _ := websocket.NewConfig(creds.URL, "ws://localhost:9443/")
	config.TlsConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	config.Header["X-Console-Token"] = []string{creds.Token}
	ws, err := websocket.DialConfig(config)
	if err != nil {
		return err
	}
	defer ws.Close()
	logrus.Println("Connection opened")

	oldState, err := term.SetRawTerminal(fdIn)
	if err != nil {
		return err
	}
	defer term.RestoreTerminal(fdIn, oldState)

	signal.Notify(make(chan os.Signal, 1), os.Interrupt)

	done := make(chan struct{}, 2)
	go readWS(ws, stdout, done)
	go readStdin(stdin, ws, done)

	<-done
	return nil
}

func (c *SConsole) Request(command string, service *models.Service) (*models.Job, error) {
	console := map[string]string{}
	if command != "" {
		console["command"] = command
	}
	b, err := json.Marshal(console)
	if err != nil {
		return nil, err
	}
	headers := c.Settings.HTTPManager.GetHeaders(c.Settings.SessionToken, c.Settings.Version, c.Settings.Pod, c.Settings.UsersID)
	resp, statusCode, err := c.Settings.HTTPManager.Post(b, fmt.Sprintf("%s%s/environments/%s/services/%s/console", c.Settings.PaasHost, c.Settings.PaasHostVersion, c.Settings.EnvironmentID, service.ID), headers)
	if err != nil {
		return nil, err
	}
	var job models.Job
	err = c.Settings.HTTPManager.ConvertResp(resp, statusCode, &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (c *SConsole) RetrieveTokens(jobID string, service *models.Service) (*models.ConsoleCredentials, error) {
	headers := c.Settings.HTTPManager.GetHeaders(c.Settings.SessionToken, c.Settings.Version, c.Settings.Pod, c.Settings.UsersID)
	resp, statusCode, err := c.Settings.HTTPManager.Post(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/jobs/%s/console-token", c.Settings.PaasHost, c.Settings.PaasHostVersion, c.Settings.EnvironmentID, service.ID, jobID), headers)
	if err != nil {
		return nil, err
	}
	var credentials models.ConsoleCredentials
	err = c.Settings.HTTPManager.ConvertResp(resp, statusCode, &credentials)
	if err != nil {
		return nil, err
	}
	return &credentials, nil
}

func (c *SConsole) Destroy(jobID string, service *models.Service) error {
	return c.Jobs.Delete(jobID, service.ID)
}

// Reads incoming data from the websocket and forwards it to stdout.
func readWS(ws *websocket.Conn, t io.Writer, done chan struct{}) {
	_, err := io.Copy(t, ws)
	if err == io.EOF {
		logrus.Println("Connection closed")
	} else if err != nil {
		logrus.Printf("Error reading data from server: %s", err)
	}
	done <- struct{}{}
}

// Reads data from stdin and writes it to the websocket.
func readStdin(t io.ReadCloser, ws *websocket.Conn, done chan struct{}) {
	_, err := io.Copy(ws, t)
	if err == io.EOF {
		logrus.Println("Input closed")
	} else if err != nil {
		logrus.Printf("Error writing data to server: %s", err)
	}
	done <- struct{}{}
}
