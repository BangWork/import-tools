package constants

const (
	ImportTypeJira = 1

	UnAssignedIssueTypeID = "0"

	NotificationValueTypeUserID              = 1
	NotificationValueTypeUserGroupID         = 2
	NotificationValueTypeGlobalProjectRoleID = 3
	NotificationValueTypeRole                = 4

	DefaultHash = 0

	CustomFieldPrefix = "customfield_"

	SingleUserLabel = "single_user"
	GroupLabel      = "group"
	RoleLabel       = "role"
)

const (
	ProjectContext   = "project"
	IssueTypeContext = "issue_type"
)
const (
	ChangeItemFieldAttachment           = "Attachment"
	ChangeItemFieldDescription          = "description"
	ChangeItemFieldPriority             = "priority"
	ChangeItemFieldStatus               = "status"
	ChangeItemFieldIssueType            = "issuetype"
	ChangeItemFieldAssignee             = "assignee"
	ChangeItemFieldSprint               = "Sprint"
	ChangeItemFieldLink                 = "Link"
	ChangeItemFieldEpicChild            = "Epic Child"
	ChangeItemFieldTimeOriginalEstimate = "timeoriginalestimate"
	ChangeItemFieldTimeEstimate         = "timeestimate"
	ChangeItemFieldTimeSpent            = "timespent"
)

const (
	ThirdTaskWorkLogTypeEstimate  = 1
	ThirdTaskWorkLogTypeRemaining = 2
	ThirdTaskWorkLogTypeLog       = 3
)

const (
	TaskWorkLogDescMaxLen = 1024
)
