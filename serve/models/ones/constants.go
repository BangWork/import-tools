package ones

const (
	// team
	interruptImport        = "/team/%s/queue/%s/interrupt"
	confirmImportUri       = "team/%s/jira/import"
	setPasswordUri         = "team/%s/jira/user/set_password"
	importStatusUri        = "team/%s/queues/list"
	importProgressUri      = "team/%s/queue/%s/progress"
	sendImportDataUri      = "team/%s/importer/import/%s"
	issueTypeListUri       = "team/%s/issue_types"
	thirdIssueTypeBindUri  = "team/%s/third_issue_type_bind"
	listPermissionRulesUri = "team/%s/lite_context_permission_rules"
	unboundIssueTypesUri   = "team/%s/importer/unbound_issue_types"

	// organization
	importLogUri       = "organization/%s/importer/log/%s/%d"
	importHistoryUri   = "organization/%s/importer/history"
	stampsUri          = "organization/%s/stamps/data"
	getOrgConfigUri    = "organization/%s/org_config/%d"
	updateOrgConfigUri = "organization/%s/org_config/update"

	// common
	fileConfigUri   = "importer/file_config"
	checkVersionUri = "importer/check_version/%s"

	// login
	loginUri = "auth/v2/login"

	dataTypeOrgPermissionRules = "org_permission_rules"
	administerOrganization     = "administer_organization"
	superAdministrator         = "super_administrator"
	singleUser                 = "single_user"
)

const (
	orgConfigTypeEnum = 1015
)

const (
	LanguageTagEnglish = "en"
	LanguageTagChinese = "zh"

	JiraSubTaskStyle          = "jira_subtask"
	IssueTypeStandardTaskType = 0
	IssueTypeSubTaskType      = 1
)
