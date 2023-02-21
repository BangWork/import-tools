package common

const (
	DefaultHTTPPort = 5000

	AuthToken     = "ones-auth-token"
	UserID        = "ones-user-id"
	JiraExportDir = "export"

	TagProject   = "Project"
	TagIssueType = "IssueType"

	TagEntityRoot = "entity-engine-xml"

	TagObjectFile = "ObjectFile"

	IssueTypeDetailTypeCustom  = 0
	IssueTypeDetailTypeDemand  = 1
	IssueTypeDetailTypeTask    = 2
	IssueTypeDetailTypeDefect  = 3
	IssueTypeDetailTypeSubTask = 4

	ResolveStatusInProgress = 1
	ResolveStatusDone       = 2
	ResolveStatusFail       = 3

	ImportStatusInProgress = 1
	ImportStatusDone       = 2
	ImportStatusPause      = 3
	ImportStatusCancel     = 4
	ImportStatusFail       = 5

	ImportStatusLabelInProgress  = "in_progress"
	ImportStatusLabelDone        = "done"
	ImportStatusLabelFail        = "fail"
	ImportStatusLabelInterrupted = "interrupted"

	FileStorageLocal = "local"

	Path      = "/var/tmp/ones-files/cache"
	XmlDir    = "xml"
	OutputDir = "output"

	ImportTypeImportTools = 2
	Calculating           = -1

	ShareDiskPathPrivate = "private"
)

var (
	CachePath string
)

func GetCachePath() string {
	return CachePath
}

func SetCachePath(path string) {
	CachePath = path
}
