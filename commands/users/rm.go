package users

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/httpclient"
)

func CmdRm(email string, iu IUsers) error {
	orgUsers, err := iu.List()
	if err != nil {
		return err
	}
	usersID := ""
	for _, u := range *orgUsers {
		if u.Email == email {
			usersID = u.ID
			break
		}
	}
	if usersID == "" {
		return fmt.Errorf("A user with email %s was not found", email)
	}

	err = iu.Rm(usersID)
	if err != nil {
		return err
	}
	logrus.Println("Removed.")
	return nil
}

func (u *SUsers) Rm(usersID string) error {
	headers := httpclient.GetHeaders(u.Settings.SessionToken, u.Settings.Version, u.Settings.Pod)
	resp, statusCode, err := httpclient.Delete(nil, fmt.Sprintf("%s%s/orgs/%s/users/%s", u.Settings.AuthHost, u.Settings.AuthHostVersion, u.Settings.OrgID, usersID), headers)
	if err != nil {
		return err
	}
	return httpclient.ConvertResp(resp, statusCode, nil)
}
