package git

import "os/exec"

// Rm removes an existing git remote from a git repo in the current working
// directory.
func (g *SGit) Rm(remote string) error {
	_, err := exec.Command("git", "remote", "remove", remote).Output()
	return err
}
