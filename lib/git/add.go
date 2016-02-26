package git

import "os/exec"

// Add a git remote to a git repo in the current working directory. If the
// current working directory is not yet a git repo, this command will fail.
func (g *SGit) Add(remote, gitURL string) error {
	_, err := exec.Command("git", "remote", "add", remote, gitURL).Output()
	return err
}
