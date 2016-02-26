package invites

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/auth"
	"github.com/catalyzeio/cli/lib/httpclient"
	"github.com/catalyzeio/cli/lib/prompts"
)

// CmdAccept creates an environment from a JSON env spec
func CmdAccept(inviteCode string, ii IInvites, ia auth.IAuth, ip prompts.IPrompts) error {
	user, err := ia.Signin()
	if err != nil {
		return err
	}
	err = ip.YesNo(fmt.Sprintf("Are you sure you want to accept this org invitation as %s? (y/n) ", user.Email))
	if err != nil {
		return err
	}

	orgID, err := ii.Accept(inviteCode)
	if err != nil {
		return err
	}
	logrus.Printf("Successfully joined organization (%s) as %s\n", orgID, user.Email)
	return nil
}

func (i *SInvites) Accept(inviteCode string) (string, error) {
	headers := httpclient.GetHeaders(i.Settings.SessionToken, i.Settings.Version, i.Settings.Pod)
	resp, statusCode, err := httpclient.Post(nil, fmt.Sprintf("%s%s/orgs/accept-invite/%s", i.Settings.AuthHost, i.Settings.AuthHostVersion, inviteCode), headers)
	if err != nil {
		return "", err
	}
	var org map[string]string
	err = httpclient.ConvertResp(resp, statusCode, &org)
	if err != nil {
		return "", err
	}
	return org["orgID"], nil
}
