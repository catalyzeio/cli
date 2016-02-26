package update

import (
	"github.com/Sirupsen/logrus"
	"github.com/catalyzeio/cli/lib/updater"
)

func CmdUpdate(iu IUpdate) error {
	logrus.Println("Checking for available updates...")
	needsUpdate, err := iu.Check()
	if err != nil {
		return err
	}
	if !needsUpdate {
		logrus.Println("You are already running the latest version of the Catalyze CLI")
		return nil
	}
	err = iu.Update()
	if err != nil {
		return err
	}
	logrus.Printf("Your CLI has been updated to version %s", updater.AutoUpdater.Info.Version)
	return nil
}

func (u *SUpdate) Check() (bool, error) {
	updater.AutoUpdater.FetchInfo()
	if updater.AutoUpdater.CurrentVersion == updater.AutoUpdater.Info.Version {
		return false, nil
	}
	return true, nil
}

// Update updates the  CLI if a new update is available.
func (u *SUpdate) Update() error {
	updater.AutoUpdater.ForcedUpgrade()
	return nil
}
