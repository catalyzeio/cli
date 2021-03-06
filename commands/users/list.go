package users

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/commands/invites"
	"github.com/catalyzeio/cli/models"
	"github.com/olekukonko/tablewriter"
	"github.com/pmylund/sortutil"
)

func CmdList(myUsersID string, iu IUsers, ii invites.IInvites) error {
	orgUsers, err := iu.List()
	if err != nil {
		return err
	}
	if orgUsers == nil || len(*orgUsers) == 0 {
		logrus.Println("No users found")
		return nil
	}
	roles, err := ii.ListRoles()
	if err != nil {
		return err
	}
	rolesMap := map[int]string{}
	for _, r := range *roles {
		rolesMap[r.ID] = r.Name
	}

	sortutil.DescByField(*orgUsers, "RoleID")

	data := [][]string{{"EMAIL", "ROLE"}}
	for _, user := range *orgUsers {
		if user.ID == myUsersID {
			data = append(data, []string{user.Email, fmt.Sprintf("%s (you)", rolesMap[user.RoleID])})
		} else {
			data = append(data, []string{user.Email, rolesMap[user.RoleID]})
		}
	}

	table := tablewriter.NewWriter(logrus.StandardLogger().Out)
	table.SetBorder(false)
	table.SetRowLine(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.AppendBulk(data)
	table.Render()
	return nil
}

func (u *SUsers) List() (*[]models.OrgUser, error) {
	headers := u.Settings.HTTPManager.GetHeaders(u.Settings.SessionToken, u.Settings.Version, u.Settings.Pod, u.Settings.UsersID)
	resp, statusCode, err := u.Settings.HTTPManager.Get(nil, fmt.Sprintf("%s%s/orgs/%s/users", u.Settings.AuthHost, u.Settings.AuthHostVersion, u.Settings.OrgID), headers)
	if err != nil {
		return nil, err
	}
	var users []models.OrgUser
	err = u.Settings.HTTPManager.ConvertResp(resp, statusCode, &users)
	if err != nil {
		return nil, err
	}
	return &users, nil
}
