package git

import "os"

// Exists tells you whether or not a git repo exists in the current working
// directory.
func (g *SGit) Exists() bool {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return false
	}
	return true
}
