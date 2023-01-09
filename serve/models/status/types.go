package status

const (
	TaskStatusStatusNormal = 1

	TaskStatusCategoryToDo       = 1
	TaskStatusCategoryInProgress = 2
	TaskStatusCategoryDone       = 3

	TaskStatusCategoryToDoLabel       = "to_do"
	TaskStatusCategoryInProgressLabel = "in_progress"
	TaskStatusCategoryDoneLabel       = "done"

	DetailTypeDefault     = 0
	DetailTypeUnPublished = 1005
	DetailTypePublished   = 3004
	DetailTypeClosed      = 3008
)

type BuiltInTaskStatus struct {
	Name       string
	Category   string
	DetailType int
}

type TaskStatus struct {
	UUID       string `db:"uuid"`
	TeamUUID   string `db:"team_uuid"`
	Name       string `db:"name" update:"true"`
	NamePinyin string `db:"name_pinyin" update:"true"`
	Category   int    `db:"category" update:"true"`
	BuiltIn    bool   `db:"built_in"`
	CreateTime int64  `db:"create_time"`
	Status     int    `db:"status" update:"true"`
	DetailType int    `db:"detail_type"`
	Position   int    `db:"position" auto:"false"`
}

func (o *TaskStatus) TableName() string {
	return "task_status"
}

type IssueTypeTaskStatusConfig struct {
	IssueTypeUUID string `db:"issue_type_uuid"`
	TaskStatusConfig
}

type TaskStatusConfig struct {
	IssueTypeScopeUUID string `db:"issue_type_scope_uuid"`
	StatusUUID         string `db:"status_uuid"`
	Default            bool   `db:"is_default"`
	Position           int64  `db:"position"`
}

type TaskStatusWithPosition struct {
	StatusUUID string `db:"task_status_uuid"`
	Position   int64  `db:"position"`
}

type TaskStatusKey struct {
	ProjectUUID   string
	IssueTypeUUID string
}

type TaskStatusConfigSorter []*TaskStatusConfig

func (s TaskStatusConfigSorter) Len() int      { return len(s) }
func (s TaskStatusConfigSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TaskStatusConfigSorter) Less(i, j int) bool {
	return s[i].Position < s[j].Position
}
