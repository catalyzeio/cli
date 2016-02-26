package git

import (
	"os/exec"
	"strings"
)

// List returns a list of all git removes in the current working directory.
func (g *SGit) List() ([]string, error) {
	out, err := exec.Command("git", "remote").Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(string(out), "\n"), nil
}
