package invites

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/prompts"
	"github.com/catalyzeio/cli/models"
)

func CmdSend(email, envName, roleName string, ii IInvites, ip prompts.IPrompts) error {
	err := ip.YesNo("", fmt.Sprintf("Are you sure you want to invite %s to your %s organization? (y/n) ", email, envName))
	if err != nil {
		return err
	}
	roles, err := ii.ListRoles()
	if err != nil {
		return err
	}
	role := 5
	for _, r := range *roles {
		if strings.ToLower(r.Name) == strings.ToLower(roleName) {
			role = r.ID
			break
		}
	}
	err = ii.Send(email, role)
	if err != nil {
		return err
	}
	logrus.Printf("%s has been invited!", email)
	return nil
}

// Send invites a user by email to the associated environment. They do
// not need a Dashboard account prior to inviting them, but they must have a
// Dashboard account in order to accept the invitation.
func (i *SInvites) Send(email string, role int) error {
	inv := models.PostInvite{
		Email:        email,
		Role:         role,
		LinkTemplate: fmt.Sprintf("%s/accept-invite?code={inviteCode}", i.Settings.AccountsHost),
	}
	b, err := json.Marshal(inv)
	if err != nil {
		return err
	}
	headers := i.Settings.HTTPManager.GetHeaders(i.Settings.SessionToken, i.Settings.Version, i.Settings.Pod, i.Settings.UsersID)
	resp, statusCode, err := i.Settings.HTTPManager.Post(b, fmt.Sprintf("%s%s/orgs/%s/invites", i.Settings.AuthHost, i.Settings.AuthHostVersion, i.Settings.OrgID), headers)
	if err != nil {
		return err
	}
	return i.Settings.HTTPManager.ConvertResp(resp, statusCode, nil)
}
