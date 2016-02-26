package disassociate

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/config"
)

func CmdDisassociate(alias string, id IDisassociate) error {
	err := id.Disassociate(alias)
	if err != nil {
		return err
	}
	logrus.Printf("WARNING: Your existing git remote *has not* been removed.\n")
	logrus.Println("Association cleared.")
	return nil
}

// Disassociate removes an existing association with the environment. The
// `catalyze` remote on the local github repo will *NOT* be removed.
func (d *SDisassociate) Disassociate(alias string) error {
	// DeleteBreadcrumb removes the environment from the settings.Environments
	// array for you
	config.DeleteBreadcrumb(alias, d.Settings)
	return nil
}
