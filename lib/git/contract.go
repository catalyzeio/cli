package git

// IGit is an interface through which you can perform git operations
type IGit interface {
	Add(remote, gitURL string) error
	Exists() bool
	List() ([]string, error)
	Rm(remote string) error
}

// SGit is an implementor of IGit
type SGit struct{}

// New creates a new instance of IGit
func New() IGit {
	return &SGit{}
}
