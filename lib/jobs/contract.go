package jobs

import "github.com/catalyzeio/cli/models"

// IJobs
type IJobs interface {
	Delete(jobID, svcID string) error
	Retrieve(jobID, svcID string, includeSpec bool) (*models.Job, error)
	RetrieveByStatus(status string) (*[]models.Job, error)
	RetrieveByType(jobType string, page, pageSize int) (*[]models.Job, error)
	PollForStatus(statuses []string, jobID, svcID string) (string, error)
	PollTillFinished(jobID, svcID string) (string, error)
	List(svcID string, page, pageSize int) (*[]models.Job, error)
}

// SJobs is a concrete implementation of IJobs
type SJobs struct {
	Settings *models.Settings
}

// New returns an instance of IJobs
func New(settings *models.Settings) IJobs {
	return &SJobs{
		Settings: settings,
	}
}
