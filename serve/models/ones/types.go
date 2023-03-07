package ones

type Stamps struct {
	OrgPermissionRules PermissionRules `json:"org_permission_rules"`
}

type PermissionRules struct {
	Rules []PermissionRule `json:"permission_rules"`
}
type PermissionRule struct {
	Permission      string `json:"permission"`
	UserDomainType  string `json:"user_domain_type"`
	UserDomainParam string `json:"user_domain_param"`
}

type loginErrorResponse struct {
	RetryCount int `json:"retry_count"`
}

type ImportHistory struct {
	TeamUUID   string `json:"team_uuid"`
	TeamName   string `json:"team_name"`
	ImportList []struct {
		ImportTime   int    `json:"import_time"`
		JiraVersion  string `json:"jira_version"`
		JiraServerID string `json:"jira_server_id"`
	} `json:"import_list"`
}

type LoginResponse struct {
	User  User         `json:"user"`
	Org   Organization `json:"org"`
	Teams []Team       `json:"teams"`
}

type User struct {
	UUID string `json:"uuid"`
}

type Organization struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	MultiTeam bool   `json:"visibility"`
}

type Team struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
