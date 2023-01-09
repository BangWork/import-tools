package issuetype

const (
	StatusNormal = 1
)

const (
	StandardTaskType = 0
	SubTaskType      = 1
)

const (
	DetailTypeCustom  = 0
	DetailTypeSubTodo = 8
	DetailTypePublish = 9
)

type IssueType struct {
	UUID            string `db:"uuid"`
	TeamUUID        string `db:"team_uuid"`
	Name            string `db:"name"`
	NamePinyin      string `db:"name_pinyin"`
	Icon            int    `db:"icon"`
	BuiltIn         bool   `db:"built_in"`
	DefaultSelected bool   `db:"default_selected"`
	CreateTime      int64  `db:"create_time"`
	Status          int    `db:"status"`
	DefaultConfigs  string `db:"default_configs"`
	Type            int    `db:"type"`
	DetailType      int    `db:"detail_type"`
	Index           int
}

var (
	WorkFlowReadOnlyIssueTypes = map[int]struct{}{
		DetailTypeSubTodo: {},
	}
)

func (i *IssueType) IsWorkFlowEditable() bool {
	return IsWorkFlowEditable(i.DetailType)
}

func IsWorkFlowEditable(detailType int) bool {
	_, readonly := WorkFlowReadOnlyIssueTypes[detailType]
	return !readonly
}

func (i *IssueType) IsSubTaskType() bool {
	return i.Type == SubTaskType
}

type IssueTypeConfig struct {
	Scope         string `db:"scope" json:"scope"`
	IssueTypeUUID string `db:"issue_type_uuid" json:"issue_type_uuid"`
}

type IssueTypeConfigWithIssueTypePinyin struct {
	Scope               string `db:"scope" json:"scope"`
	IssueTypeUUID       string `db:"issue_type_uuid" json:"issue_type_uuid"`
	IssueTypeNamePinyin string `db:"issue_type_pinyin" json:"-"`
}

type ProjectIssueType struct {
	Scope         string `db:"scope"`
	IssueTypeUUID string `db:"issue_type_uuid"`
	Position      string `db:"position"`
}

type Scope struct {
	UUID          string `db:"uuid" json:"uuid"`
	TeamUUID      string `db:"team_uuid" json:"-"`
	Name          string `db:"name" json:"name"`
	NamePinyin    string `db:"name_pinyin" json:"name_pinyin"`
	Scope         string `db:"scope" json:"scope"`
	ScopeType     int    `db:"scope_type" json:"scope_type"`
	IssueTypeUUID string `db:"issue_type_uuid" json:"issue_type_uuid"`
	Position      int    `db:"position" json:"-"`
}
