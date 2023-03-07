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

	// organization
	importLogUri     = "organization/%s/importer/log/%s/%d"
	importHistoryUri = "organization/%s/importer/history"
	fileConfigUri    = "organization/%s/file_config"
	stampsUri        = "organization/%s/stamps/data"

	// login
	loginUri = "auth/v2/login"

	dataTypeOrgPermissionRules = "org_permission_rules"
	administerOrganization     = "administer_organization"
	superAdministrator         = "super_administrator"
	singleUser                 = "single_user"
)

