package jira

const (
	entitiesFile        = "entities.xml"
	activeObjectsFile   = "activeobjects.xml"
	entityRootTag       = "entity-engine-xml"
	activeObjectRootTag = "backup"

	jiraSubtaskType = "jira_subtask"
)

const (
	roleTypeUser  = "atlassian-user-role-actor"
	roleTypeGroup = "atlassian-group-role-actor"
)

const (
	configContextIssueType = "issuetype"
)

const timeLayout = "2006-01-02 15:04:05"

const (
	jiraStatusID      = "jira.status.id"
	taskStatusAllUUID = "ALL_UUID"
)

const (
	associationTypeIssueComponent  = "IssueComponent"
	associationTypeIssueVersion    = "IssueVersion"
	associationTypeIssueFixVersion = "IssueFixVersion"
	associationTypeWatchIssue      = "WatchIssue"
)

const (
	customFieldID               = "ID"
	customFieldReleaseStartDate = "Release Start Date"
	customFieldReleaseDate      = "Release Date"
	customFieldLabels           = "Labels"
	customFieldFixVersion       = "Fix Version/s"
	customFieldVersion          = "Affects Version/s"
	customFieldResolution       = "Resolution"
	customFieldComponent        = "Component/s"
	customFieldEnvironment      = "Environment"
)

const (
	applicationRole = "applicationRole"
	singleUser      = "user"
	projectRole     = "projectrole"
	userGroup       = "group"
	assignee        = "assignee"
	reporter        = "reporter"
	lead            = "lead"
)

const (
	MANAGE_SPRINTS_PERMISSION = "MANAGE_SPRINTS_PERMISSION"
	ADMINISTER_PROJECTS       = "ADMINISTER_PROJECTS"
	BROWSE_PROJECTS           = "BROWSE_PROJECTS"
	CREATE_ISSUES             = "CREATE_ISSUES"
	DELETE_ISSUES             = "DELETE_ISSUES"
	ASSIGNABLE_USER           = "ASSIGNABLE_USER"
	MANAGE_WATCHERS           = "MANAGE_WATCHERS"
	EDIT_ISSUES               = "EDIT_ISSUES"
	ASSIGN_ISSUES             = "ASSIGN_ISSUES"
	LINK_ISSUES               = "LINK_ISSUES"
	CLOSE_ISSUES              = "CLOSE_ISSUES"
	RESOLVE_ISSUES            = "RESOLVE_ISSUES"
	TRANSITION_ISSUES         = "TRANSITION_ISSUES"
	DELETE_ALL_WORKLOGS       = "DELETE_ALL_WORKLOGS"
	EDIT_ALL_WORKLOGS         = "EDIT_ALL_WORKLOGS"
	DELETE_OWN_WORKLOGS       = "DELETE_OWN_WORKLOGS"
	EDIT_OWN_WORKLOGS         = "EDIT_OWN_WORKLOGS"
)
