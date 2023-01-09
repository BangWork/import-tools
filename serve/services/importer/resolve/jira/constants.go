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
	ManageSprintsPermission = "MANAGE_SPRINTS_PERMISSION"
	AdministerProjects  = "ADMINISTER_PROJECTS"
	BrowseProjects      = "BROWSE_PROJECTS"
	CreateIssues      = "CREATE_ISSUES"
	DeleteIssues      = "DELETE_ISSUES"
	AssignableUser    = "ASSIGNABLE_USER"
	ManageWatchers    = "MANAGE_WATCHERS"
	EditIssues        = "EDIT_ISSUES"
	AssignIssues      = "ASSIGN_ISSUES"
	LinkIssues        = "LINK_ISSUES"
	CloseIssues       = "CLOSE_ISSUES"
	ResolveIssues     = "RESOLVE_ISSUES"
	TransitionIssues  = "TRANSITION_ISSUES"
	DeleteAllWorkLogs = "DELETE_ALL_WORKLOGS"
	EditAllWorkLogs   = "EDIT_ALL_WORKLOGS"
	DeleteOwnWorkLogs = "DELETE_OWN_WORKLOGS"
	EditOwnWorkLogs   = "EDIT_OWN_WORKLOGS"
)
