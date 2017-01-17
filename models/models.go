package models

import "github.com/jault3/mow.cli"

// AssociatedEnvV1 holds information about an associated environment
type AssociatedEnvV1 struct {
	EnvironmentID string `json:"environmentId"`
	ServiceID     string `json:"serviceId"`
	Directory     string `json:"dir"`
	Name          string `json:"name"`
	Pod           string `json:"pod"`
	OrgID         string `json:"organizationId"`
}

// AssociatedEnvV2 holds information about an associated environment
type AssociatedEnvV2 struct {
	EnvironmentID string `json:"environmentId"`
	Name          string `json:"name"`
	Pod           string `json:"pod"`
	OrgID         string `json:"organizationId"`
}

type Cert struct {
	Name    string `json:"name"`
	PubKey  string `json:"sslCertFile"`
	PrivKey string `json:"sslPKFile"`

	Service    string `json:"service,omitempty"`
	PubKeyID   int    `json:"sslCertFileId,omitempty"`
	PrivKeyID  int    `json:"sslPKFileId,omitempty"`
	Restricted bool   `json:"restricted,omitempty"`
}

type Command struct {
	Name      string
	ShortHelp string
	LongHelp  string
	CmdFunc   func(settings *Settings) func(cmd *cli.Cmd)
}

// ConsoleCredentials hold the keys necessary for connecting to a console service
type ConsoleCredentials struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

type CPUUsage struct {
	JobID       string  `json:"job"`
	CorePercent float64 `json:"core_percent"`
	TS          int     `json:"ts"`
}

// DeployKey is an ssh key belonging to an environment's code service
type DeployKey struct {
	Name string `json:"name"`
	Key  string `json:"value"`
	Type string `json:"type"`
}

// EncryptionStore holds the values for encryption on backup/import jobs
type EncryptionStore struct {
	Key             string `json:"key"`
	KeyLogs         string `json:"keyLogs"`
	KeyInternalLogs string `json:"keyInternalLogs"`
	IV              string `json:"iv"`
}

// Environment environment
type Environment struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	Pod       string `json:"pod,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	OrgID     string `json:"organizationId"`
}

// Error is a wrapper around an array of errors from the API
type Error struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Code        int    `json:"code"`
}

// Hits contain arrays of log data
type Hits struct {
	Total    int64      `json:"total"`
	MaxScore float64    `json:"max_score"`
	Hits     *[]LogHits `json:"hits"`
}

type HTTPManager interface {
	GetHeaders(sessionToken, version, pod, userID string) map[string][]string
	ConvertResp(b []byte, statusCode int, s interface{}) error
	Get(body []byte, url string, headers map[string][]string) ([]byte, int, error)
	Post(body []byte, url string, headers map[string][]string) ([]byte, int, error)
	PostFile(filepath string, url string, headers map[string][]string) ([]byte, int, error)
	PutFile(filepath string, url string, headers map[string][]string) ([]byte, int, error)
	Put(body []byte, url string, headers map[string][]string) ([]byte, int, error)
	Delete(body []byte, url string, headers map[string][]string) ([]byte, int, error)
}

// Invite represents an invitation to an organization
type Invite struct {
	ID       string `json:"id"`
	OrgID    string `json:"orgID"`
	SenderID string `json:"senderID"`
	RoleID   int    `json:"roleID"`
	Email    string `json:"email"`
	Consumed bool   `json:"consumed"`
	Revoked  bool   `json:"revoked"`
}

// LogHits contain ordering data for logs
type LogHits struct {
	Index  string              `json:"_index"`
	Type   string              `json:"_type"`
	ID     string              `json:"_id"`
	Score  float64             `json:"_score"`
	Fields map[string][]string `json:"fields"`
}

// Login is used for making an authentication request
type Login struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

// Logs hold the log values from a successful LogQuery
type Logs struct {
	Hits *Hits `json:"hits"`
}

type Maintenance struct {
	UpstreamID string `json:"upstream"`
	CreatedAt  string `json:"createdAt"`
}

type MemoryUsage struct {
	JobID string  `json:"job"`
	Total float64 `json:"total"`
	AVG   float64 `json:"ave"`
	Max   float64 `json:"max"`
	Min   float64 `json:"min"`
	TS    int     `json:"ts"`
}

// Metrics holds all metrics data for an entire environment or a single service
type Metrics struct {
	ServiceName  string       `json:"serviceName"`
	ServiceType  string       `json:"serviceType"`
	ServiceID    string       `json:"serviceId"`
	ServiceLabel string       `json:"serviceLabel"`
	Size         ServiceSize  `json:"size"`
	Data         *MetricsData `json:"metrics"`
}

// MetricsData is a container for each type of metrics: network, memory, etc.
type MetricsData struct {
	CPUUsage     *[]CPUUsage     `json:"cpu.usage"`
	MemoryUsage  *[]MemoryUsage  `json:"memory.usage"`
	NetworkUsage *[]NetworkUsage `json:"network.usage"`
}

type NetworkUsage struct {
	JobID     string  `json:"job"`
	RXDropped float64 `json:"rx_dropped"`
	RXErrors  float64 `json:"rx_errors"`
	RXKB      float64 `json:"rx_kb"`
	RXPackets float64 `json:"rx_packets"`
	TXDropped float64 `json:"tx_dropped"`
	TXErrors  float64 `json:"tx_errors"`
	TXKB      float64 `json:"tx_kb"`
	TXPackets float64 `json:"tx_packets"`
	TS        int     `json:"ts"`
}

type Org struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// OrgUser users who have access to an org
type OrgUser struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	RoleID int    `json:"roleID"`
}

// Payload is the payload of a job
type Payload struct {
	Environment map[string]string `json:"environment"`
}

// Pod is a pod returned from the pod router
type Pod struct {
	Name    string `json:"name"`
	PHISafe bool   `json:"phiSafe"`
}

// Job job
type Job struct {
	ID               string           `json:"id"`
	Type             string           `json:"type"`
	Status           string           `json:"status,omitempty"`
	Backup           *EncryptionStore `json:"backup,omitempty"`
	Restore          *EncryptionStore `json:"restore,omitempty"`
	CreatedAt        string           `json:"created_at"`
	MetricsData      *[]MetricsData   `json:"metrics"`
	Spec             *Spec            `json:"spec"`
	Target           string           `json:"target,omitempty"`
	IsSnapshotBackup *bool            `json:"isSnapshotBackup,omitempty"`
}

// PodWrapper pod wrapper
type PodWrapper struct {
	Pods *[]Pod `json:"pods"`
}

type PostInvite struct {
	Email        string `json:"email"`
	Role         int    `json:"role"`
	LinkTemplate string `json:"linkTemplate"`
}

type Release struct {
	Name      string `json:"release,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Notes     string `json:"metadata,omitempty"`
}

// ReportedError is the standard error model sent back from the API
type ReportedError struct {
	Code    int    `json:"id"`
	Message string `json:"message"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Service service
type Service struct {
	ID             string            `json:"id,omitempty"`
	Identifier     string            `json:"identifier,omitempty"`
	DNS            string            `json:"internal_domain,omitempty"`
	Type           string            `json:"type,omitempty"`
	Label          string            `json:"label"`
	Size           ServiceSize       `json:"size"`
	Name           string            `json:"name"`
	EnvVars        map[string]string `json:"environmentVariables,omitempty"`
	Source         string            `json:"source,omitempty"`
	LBIP           string            `json:"load_balancer_ip,omitempty"`
	Scale          int               `json:"scale,omitempty"`
	WorkerScale    int               `json:"worker_scale,omitempty"`
	ReleaseVersion string            `json:"release_version,omitempty"`
	Redeployable   bool              `json:"redeployable,omitempty"`
}

// ServiceFile is a file associated with a service
type ServiceFile struct {
	ID             int    `json:"id"`
	Contents       string `json:"contents"`
	GID            int    `json:"gid"`
	Mode           string `json:"mode"`
	Name           string `json:"name"`
	UID            int    `json:"uid"`
	EnableDownload bool   `json:"enable_download"`
}

// ServiceSize holds size information for a service
type ServiceSize struct {
	RAM      int    `json:"ram"`
	Storage  int    `json:"storage"`
	Behavior string `json:"behavior,omitempty"`
	Type     string `json:"type,omitempty"`
	CPU      int    `json:"cpu"`
}

type Settings SettingsV2

// SettingsV1 holds various settings for the current context. All items with
// `json:"-"` are never persisted to disk but used in memory for the current
// command.
type SettingsV1 struct {
	AccountsHost    string      `json:"-"`
	AuthHost        string      `json:"-"`
	PaasHost        string      `json:"-"`
	AuthHostVersion string      `json:"-"`
	PaasHostVersion string      `json:"-"`
	Version         string      `json:"-"`
	HTTPManager     HTTPManager `json:"-"`

	Username        string                     `json:"-"`
	Password        string                     `json:"-"`
	EnvironmentID   string                     `json:"-"` // the id of the environment used for the current command
	ServiceID       string                     `json:"-"` // the id of the service used for the current command
	Pod             string                     `json:"-"` // the pod used for the current command
	EnvironmentName string                     `json:"-"` // the name of the environment used for the current command
	OrgID           string                     `json:"-"` // the org ID of the chosen environment
	PrivateKeyPath  string                     `json:"private_key_path"`
	SessionToken    string                     `json:"token"`
	UsersID         string                     `json:"user_id"`
	Environments    map[string]AssociatedEnvV1 `json:"environments"`
	Default         string                     `json:"default"`
	Pods            *[]Pod                     `json:"pods"`
	PodCheck        int64                      `json:"pod_check"`
}

// SettingsV2 holds various settings for the current context. All items with
// `json:"-"` are never persisted to disk but used in memory for the current
// command.
type SettingsV2 struct {
	AccountsHost    string      `json:"-"`
	AuthHost        string      `json:"-"`
	PaasHost        string      `json:"-"`
	AuthHostVersion string      `json:"-"`
	PaasHostVersion string      `json:"-"`
	Version         string      `json:"-"`
	HTTPManager     HTTPManager `json:"-"`

	Username        string                     `json:"-"`
	Password        string                     `json:"-"`
	EnvironmentID   string                     `json:"-"` // the id of the environment used for the current command
	Pod             string                     `json:"-"` // the pod used for the current command
	EnvironmentName string                     `json:"-"` // the name of the environment used for the current command
	OrgID           string                     `json:"-"` // the org ID of the chosen environment
	PrivateKeyPath  string                     `json:"private_key_path"`
	SessionToken    string                     `json:"token"`
	UsersID         string                     `json:"user_id"`
	Environments    map[string]AssociatedEnvV2 `json:"environments"`
	Pods            *[]Pod                     `json:"pods"`
	PodCheck        int64                      `json:"pod_check"`
	Format          string                     `json:"format"`
}

type Site struct {
	ID              int                    `json:"id,omitempty"`
	Name            string                 `json:"name"`
	Cert            string                 `json:"cert"`
	SiteFileID      int                    `json:"siteFileId,omitempty"`
	UpstreamService string                 `json:"upstreamService"`
	Restricted      bool                   `json:"restricted,omitempty"`
	SiteValues      map[string]interface{} `json:"site_values"`
}

// Spec is a job specification
type Spec struct {
	Payload *Payload `json:"payload"`
}

// TempURL holds a URL for uploading or downloading files from a temporary URL
type TempURL struct {
	URL string `json:"url"`
}

// User is an authenticated User
type User struct {
	Username     string `json:"name"`
	Email        string `json:"email"`
	SessionToken string `json:"sessionToken"`
	UsersID      string `json:"id"`
}

// UserKey is a public key belonging to a user
type UserKey struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type Workers struct {
	Limit   int            `json:"worker_limit,omitempty"`
	Workers map[string]int `json:"workers"`
}
