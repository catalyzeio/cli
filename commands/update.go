package commands

import (
	"fmt"

	"github.com/catalyzeio/catalyze/updater"
)

// Update updates the  CLI if a new update is available.
func Update() {
	fmt.Println("Checking for available updates...")
	updater.AutoUpdater.FetchInfo()
	if updater.AutoUpdater.CurrentVersion == updater.AutoUpdater.Info.Version {
		fmt.Println("You are already running the latest 2.x version of the Catalyze CLI")
		return
	}
	updater.AutoUpdater.ForcedUpgrade()
	fmt.Printf("Your CLI has been updated to version %s\n", updater.AutoUpdater.Info.Version)
}
