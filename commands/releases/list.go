package releases

import (
	"fmt"
	"sort"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/services"
	"github.com/catalyzeio/cli/models"
	"github.com/olekukonko/tablewriter"
)

// SortedReleases is a wrapper for Release array in order to sort them by CreatedAt
type SortedReleases []models.Release

func (rls SortedReleases) Len() int {
	return len(rls)
}

func (rls SortedReleases) Swap(i, j int) {
	rls[i], rls[j] = rls[j], rls[i]
}

func (rls SortedReleases) Less(i, j int) bool {
	return rls[i].CreatedAt > rls[j].CreatedAt
}

func CmdList(svcName string, ir IReleases, is services.IServices) error {
	service, err := is.RetrieveByLabel(svcName)
	if err != nil {
		return err
	}
	if service == nil {
		return fmt.Errorf("Could not find a service with the label \"%s\". You can list services with the \"catalyze services list\" command.", svcName)
	}

	rls, err := ir.List(service.ID)
	if err != nil {
		return err
	}

	if rls == nil || len(*rls) == 0 {
		logrus.Println("No releases found")
		return nil
	}

	sort.Sort(SortedReleases(*rls))
	const dateForm = "2006-01-02T15:04:05"
	data := [][]string{{"Release Name", "Created At", "Notes"}}
	for _, r := range *rls {
		name := r.Name
		if r.Name == service.ReleaseVersion {
			name = fmt.Sprintf("*%s", r.Name)
		}
		t, _ := time.Parse(dateForm, r.CreatedAt)
		data = append(data, []string{name, t.Local().Format(time.Stamp), r.Notes})
	}

	table := tablewriter.NewWriter(logrus.StandardLogger().Out)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.AppendBulk(data)
	table.Render()

	logrus.Println("\n* denotes the current release")
	return nil
}

func (r *SReleases) List(svcID string) (*[]models.Release, error) {
	headers := r.Settings.HTTPManager.GetHeaders(r.Settings.SessionToken, r.Settings.Version, r.Settings.Pod, r.Settings.UsersID)
	resp, statusCode, err := r.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/releases", r.Settings.PaasHost, r.Settings.PaasHostVersion, r.Settings.EnvironmentID, svcID), headers)
	if err != nil {
		return nil, err
	}
	var rls []models.Release
	err = r.Settings.HTTPManager.ConvertResp(resp, statusCode, &rls)
	if err != nil {
		return nil, err
	}
	return &rls, nil
}

func (r *SReleases) Retrieve(releaseName, svcID string) (*models.Release, error) {
	headers := r.Settings.HTTPManager.GetHeaders(r.Settings.SessionToken, r.Settings.Version, r.Settings.Pod, r.Settings.UsersID)
	resp, statusCode, err := r.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/environments/%s/services/%s/releases/%s", r.Settings.PaasHost, r.Settings.PaasHostVersion, r.Settings.EnvironmentID, svcID, releaseName), headers)
	if err != nil {
		return nil, err
	}
	var rls models.Release
	err = r.Settings.HTTPManager.ConvertResp(resp, statusCode, &rls)
	if err != nil {
		return nil, err
	}
	return &rls, nil
}
