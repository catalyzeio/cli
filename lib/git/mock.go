package git

import "errors"

// MGit is a mock implementor of IGit
type MGit struct {
	ReturnError bool
}

func (g *MGit) Add(remote, gitURL string) error {
	if g.ReturnError {
		return errors.New("Mock error returned")
	}
	return nil
}

func (g *MGit) Exists() bool {
	if g.ReturnError {
		return false
	}
	return true
}

func (g *MGit) List() ([]string, error) {
	if g.ReturnError {
		return nil, errors.New("Mock error returned")
	}
	return []string{"catalyze"}, nil
}

func (g *MGit) Rm(remote string) error {
	if g.ReturnError {
		return errors.New("Mock error returned")
	}
	return nil
}
