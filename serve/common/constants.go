package common

const (
	CurrentVersion = "3.x.x"
)

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
	LoginCookieName      = "IMPORTTOOLS"

	DefaultAesKey = "1234567890abcdef"

	InstallAreaAsia    = "Asia"
	InstallAreaAmerica = "America"

	ProjectTypeSoftware = "software"
	ProjectTypeBusiness = "business"
)

const (
	ResourceTypeStringUser                   = "user"
	ResourceTypeStringUserGroup              = "user_group"
	ResourceTypeStringUserGroupMember        = "user_group_member"
	ResourceTypeStringDepartment             = "department"
	ResourceTypeStringUserDepartments        = "user_departments"
	ResourceTypeStringGlobalProjectRole      = "global_project_role"
	ResourceTypeStringGlobalProjectField     = "global_project_field"
	ResourceTypeStringIssueType              = "issue_type"
	ResourceTypeStringProject                = "project"
	ResourceTypeStringProjectIssueType       = "project_issue_type"
	ResourceTypeStringProjectRole            = "project_role"
	ResourceTypeStringProjectRoleMember      = "project_role_member"
	ResourceTypeStringGlobalPermission       = "global_permission"
	ResourceTypeStringProjectPermission      = "project_permission"
	ResourceTypeStringProjectFieldValue      = "project_field_value"
	ResourceTypeStringTaskStatus             = "task_status"
	ResourceTypeStringTaskField              = "task_field"
	ResourceTypeStringTaskFieldOption        = "task_field_option"
	ResourceTypeStringIssueTypeField         = "issue_type_field"
	ResourceTypeStringIssueTypeLayout        = "issue_type_layout"
	ResourceTypeStringProjectIssueTypeField  = "project_issue_type_field"
	ResourceTypeStringProjectIssueTypeLayout = "project_issue_type_layout"
	ResourceTypeStringProjectSprintField     = "project_sprint_field"
	ResourceTypeStringPriority               = "priority"
	ResourceTypeStringConfig                 = "config"
	ResourceTypeStringTaskLinkType           = "task_link_type"
	ResourceTypeStringWorkflow               = "workflow"
	ResourceTypeStringSprint                 = "sprint"
	ResourceTypeStringSprintFieldValue       = "sprint_field_value"
	ResourceTypeStringTask                   = "task"
	ResourceTypeStringTaskFieldValue         = "task_field_value"
	ResourceTypeStringTaskWatcher            = "task_watcher"
	ResourceTypeStringTaskWorkLog            = "task_work_log"
	ResourceTypeStringTaskComment            = "task_comment"
	ResourceTypeStringTaskRelease            = "task_release"
	ResourceTypeStringTaskLink               = "task_link"
	ResourceTypeStringNotification           = "notification"
	ResourceTypeStringTaskAttachmentTmp      = "task_attachment_tmp"
	ResourceTypeStringTaskAttachment         = "task_attachment"
	ResourceTypeStringChangeItem             = "change_item"
)
