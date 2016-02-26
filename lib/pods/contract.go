package pods

import "github.com/catalyzeio/cli/models"

// IPods
type IPods interface {
	List() (*[]models.Pod, error)
}

// SPods is a concrete implementation of IPods
type SPods struct {
	Settings *models.Settings
}

// New returns an instance of IPods
func New(settings *models.Settings) IPods {
	return &SPods{
		Settings: settings,
	}
}
