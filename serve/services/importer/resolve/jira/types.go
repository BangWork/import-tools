package jira

type changeGroup struct {
	TaskID     string
	UserID     string
	CreateTime int64
}

type fieldScreenLayoutItem struct {
	ID               string
	FieldIdentifier  string
	FieldScreenTabID string
	Sequence         string
}

type fieldLayoutItem struct {
	ID              string
	FieldLayoutID   string
	FieldIdentifier string
	IsHidden        string
	IsRequired      string
}

type fieldScreenSchemeItem struct {
	ID                  string
	Operation           string
	FieldScreenID       string
	FieldScreenSchemeID string
}

type issueTypeScreenSchemeEntity struct {
	ID                      string
	IssueTypeScreenSchemeID string
	FieldScreenSchemeID     string
	IssueTypeID             string
}

type fieldLayoutSchemeEntity struct {
	ID                  string
	FieldLayoutSchemeID string
	IssueTypeID         string
	FieldLayoutID       string
}

type fieldLayout struct {
	ID   string
	Name string
	Type string
}

type projectRoleStruct struct {
	ProjectID string
	RoleID    string
}

type projectMemberStruct struct {
	ProjectID string
	RoleID    string
	UserID    string
}

type flowMap struct {
	ProjectIDs []string
	IssueType  string
}

type roleData struct {
	roleType          string
	roleTypeParameter string
}

type permissionRecord struct {
	id            string
	scheme        string
	permissionKey string
	spType        string
	parameter     string
}
