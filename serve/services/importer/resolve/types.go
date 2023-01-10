package resolve

// 「中间层数据结构」
// 该数据代表第三方系统数据在ONES中的呈现，所有ID均为第三方数据ID

type Base struct {
	ResourceID string `json:"resource_id"` // 第三方id
}

type ThirdConfig struct {
	Workdays            []string          `json:"workdays"`
	WorkHours           int               `json:"work_hours"`
	FixedFieldsMapping  map[string]string `json:"fixed_fields_mapping"`  // 第三方字段映射系统字段
	CustomFieldsMapping map[string]string `json:"custom_fields_mapping"` // 属性映射
	CustomFields        []string          `json:"custom_fields"`         // 属性
	TabFields           map[string]string `json:"tab_fields"`            // tab 属性
	LayoutFieldIDMap    map[string]string `json:"layout_field_id_map"`   // 第三方视图字段id 替换为 ONES字段id
}

type ThirdTask struct {
	Base
	Summary             string `json:"summary"`                // 标题
	Desc                string `json:"desc"`                   // 描述
	OwnerID             string `json:"owner_id"`               // 创建者ID
	AssignID            string `json:"assign_id"`              // 负责人ID
	Deadline            *int64 `json:"deadline"`               // 最后完成日期
	PriorityID          string `json:"priority_id"`            // 优先级ID
	StatusID            string `json:"status_id"`              // 状态ID
	ProjectID           string `json:"project_id"`             // 项目ID
	IssueTypeID         string `json:"issue_type_id"`          // 工作项类型ID
	SprintID            string `json:"sprint_id"`              // 迭代ID
	ParentID            string `json:"parent_id"`              // 父工作项ID
	CreatedTime         int64  `json:"created_time"`           // 创建时间
	UpdatedTime         int64  `json:"updated_time"`           // 更新时间
	IssueTypeDetailType int64  `json:"issue_type_detail_type"` // 使用ONES内置工作项类型
	StatusDetailType    int    `json:"status_detail_type"`     // 使用ONES内置状态
}

type ThirdTaskFieldValue struct {
	Base
	TaskID    string      `json:"task_id"`
	FieldID   string      `json:"field_id"`
	FieldType int         `json:"field_type"`
	Value     interface{} `json:"value"`
}

type ThirdTaskStatus struct {
	Base
	Name     string `json:"name"`
	Category string `json:"category"`
}

type ThirdWorkflow struct {
	ProjectID         string                   `json:"project_id"`
	IssueTypeStatuses []*ThirdTaskStatusConfig `json:"issue_type_statuses"`
	Transitions       []*ThirdTransition       `json:"transitions"`
}

type ThirdTaskStatusConfig struct {
	IssueTypeID  string   `json:"issue_type_id"`
	StatusIDs    []string `json:"status_ids"`
	InitStatusID string   `json:"init_status_id"`
	WorkflowName string   `json:"workflow_name"`
}

type ThirdTransition struct {
	ProjectID     string             `json:"project_id"`
	IssueTypeID   string             `json:"issue_type_id"`
	StartStatusID string             `json:"start_status_id"`
	Name          string             `json:"name"`
	EndStatusID   string             `json:"end_status_id"`
	UserDomains   []*ThirdUserDomain `json:"user_domains"`
}

type ThirdUserDomain struct {
	UserDomainType  string `json:"user_domain_type"`
	UserDomainParam string `json:"user_domain_param"`
}

type ThirdTaskField struct {
	Base
	Name string `json:"name"`
	Type int    `json:"type"`
}

type ThirdTaskFieldOption struct {
	Base
	ResourceTaskFieldID string `json:"resource_task_field_id"`
	Name                string `json:"name"`
	Desc                string `json:"desc"`
	Color               string `json:"color"`
	BackgroundColor     string `json:"background_color"`
}

type ThirdPriority struct {
	Base
	Name            string `json:"name"`
	Desc            string `json:"desc"`
	Color           string `json:"color"`
	BackgroundColor string `json:"background_color"`
}

type ThirdTaskLinkType struct {
	Base
	Name    string `json:"name"`
	DescOut string `json:"desc_out"`
	DescIn  string `json:"desc_in"`
	Type    int    `json:"type"`
}

type ThirdUser struct {
	Base
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Title      string `json:"title"`
	Company    string `json:"company"`
	Status     int    `json:"status"`
	CreateTime int64  `json:"create_time" `
	ModifyTime int64  `json:"modify_time" `
}

type ThirdUserGroup struct {
	Base
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Status     int    `json:"status"`
	CreateTime int64  `json:"create_time"`
}

type ThirdIssueType struct {
	Base
	Name       string `json:"name"`
	Status     int    ` json:"status"`
	Type       int    `json:"type"`
	DetailType int    `json:"detail_type"`
	CreateTime int64  `json:"create_time"`
	Icon       int
}

type ThirdProjectIssueType struct {
	Base
	ProjectID   string `json:"project_id"`
	IssueTypeID string `json:"issue_type_id"`
}

type ThirdGlobalProjectRole struct {
	Base
	Name       string `json:"name"`
	Status     int    `json:"status"`
	CreateTime int64  `json:"create_time"`
}

type ThirdProjectRole struct {
	Base
	ResourceRoleID    string `json:"resource_role_id"`
	ResourceProjectID string `json:"resource_project_id"`
	CreateTime        int64  `json:"create_time"`
}

type ThirdProjectRoleMember struct {
	Base
	ResourceRoleID    string `json:"resource_role_id"`
	ResourceProjectID string `json:"resource_project_id"`
	ResourceUserID    string `json:"resource_user_id"`
}

type ThirdProjectCategory struct {
	Base
	Name string `json:"name"`
}

type ThirdUserGroupMember struct {
	Base
	ResourceGroupID string `json:"resource_group_id"`
	ResourceUserID  string `json:"resource_user_id"`
}

type ThirdProject struct {
	Base
	Name         string `json:"name"`
	AssignID     string `json:"assign_id"` // 负责人ID
	CreateTime   int64  `json:"create_time"`
	Type         int    `json:"type"`         // 项目类型，通用，瀑布
	Status       int8   `json:"status"`       // 项目状态， 1正常，2删除
	StatusUUID   string `json:"status_uuid"`  // 项目进行状态，to_do未开始，in_progress进行中，done已结束
	Announcement string `json:"announcement"` // 公告
}

type ThirdGlobalProjectField struct {
	Base
	Name    string                           `json:"name"`
	Type    int                              `json:"type"`    // 属性类型  1单选，2单行文本
	Options []*ThirdGlobalProjectFieldOption `json:"options"` // 下拉框选项值
}

type ThirdGlobalProjectFieldOption struct {
	Base
	Value           string `json:"value"`
	Desc            string `json:"desc"`
	Color           string `json:"color"`
	BackgroundColor string `json:"background_color"`
}

type ThirdProjectFieldValue struct {
	Base
	ProjectID      string      `json:"project_id"`
	ProjectFieldID string      `json:"project_field_id"`
	Type           int         `json:"type"` // 属性类型  1单选，2单行文本
	Value          interface{} `json:"value"`
}

type ThirdIssueTypeLayout struct {
	IssueTypeID       string   `json:"issue_type_id"`
	FieldConfigID     string   `json:"field_config_id"`
	FieldConfigName   string   `json:"field_config_name"`
	ScreenSchemeID    string   `json:"screen_scheme_id"`
	ScreenSchemeName  string   `json:"screen_scheme_name"`
	CreateIssueConfig []*Field `json:"create_issue_config"`
	ViewIssueConfig   []*Field `json:"view_issue_config"`
}

type ThirdChangeItem struct {
	Base
	TaskID     string `json:"task_id"`
	FieldType  string `json:"field_type"`
	Field      string `json:"field"`
	OldValue   string `json:"old_value"`
	NewValue   string `json:"new_value"`
	OldString  string `json:"old_string"`
	NewString  string `json:"new_string"`
	UserID     string `json:"user_id"`
	CreateTime int64  `json:"create_time"`
}

type ThirdIssueTypeField struct {
	Base
	IssueTypeID         string `json:"issue_type_id"`
	FieldID             string `json:"field_id"`
	IssueTypeDetailType int64  `json:"issue_type_detail_type"` // 使用ONES内置工作项类型
}

type Field struct {
	FieldIdentifier string `json:"field_identifier"`
	Required        bool   `json:"required"`
}

type ThirdProjectIssueTypeWithFields struct {
	ProjectID           string                 `json:"project_id"`
	IssueTypeWithFields []*IssueTypeWithFields `json:"issue_type_with_fields"`
}

type IssueTypeWithFields struct {
	IssueTypeID         string   `json:"issue_type_id"`
	Fields              []*Field `json:"fields"`
	IssueTypeDetailType int64    `json:"issue_type_detail_type"`
}

type ThirdProjectIssueTypeLayout struct {
	ProjectID    string             `json:"project_id"`
	ScopeConfigs []*ScopeConfigItem `json:"scope_configs"`
}

type ScopeConfigItem struct {
	IssueTypeID    string `json:"issue_type_id"`
	FieldConfigID  string `json:"field_config_id"`
	ScreenSchemeID string `json:"screen_scheme_id"`
}

type ThirdSprint struct {
	Base
	Name      string `json:"name"`
	Status    int    `json:"status"`
	ProjectID string `json:"project_id"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	AssignID  string `json:"assign_id"`
	OwnerID   string `json:"owner_id"`
}

type ThirdTaskWatcher struct {
	Base
	TaskID string `json:"task_id"`
	UserID string `json:"user_id"`
}

type ThirdProjectPermission struct {
	Base
	ProjectID       string `json:"project_id"`
	UserDomainType  string `json:"user_domain_type"`
	ContextType     string `json:"context_type"`
	Permission      string `json:"permission"`
	UserDomainParam string `json:"user_domain_param"`
}

type ThirdGlobalPermission struct {
	Base
	Permission int    `json:"permission"`
	GroupID    string `json:"group_id"`
}

type ThirdRelease struct {
	Base
	Permission int    `json:"permission"`
	GroupID    string `json:"group_id"`
}

type ThirdTaskComment struct {
	Base
	TaskID     string `json:"task_id"`
	UserID     string `json:"user_id"`
	Body       string `json:"body"`
	CreateTime int64  ` json:"create_time"`
}

type ThirdTaskLink struct {
	Base
	LinkTypeID   string ` json:"link_type_id"`
	SourceTaskID string `json:"source_task_id"`
	TargetTaskID string `json:"target_task_id"`
}

type ThirdTaskRelease struct {
	TaskID    string `json:"task_id"`
	ReleaseID string `json:"release_id"`
}

type ThirdTaskAttachment struct {
	Base
	TaskID       string `json:"task_id"`
	FilePath     string `json:"file_path"`
	FileName     string `json:"file_name"`
	ResourceUUID string `json:"resource_uuid"`
}

type ThirdTaskWorkLog struct {
	Base
	TaskID      string  `json:"task_id"`
	UserID      string  `json:"user_id"`
	StartTime   int64   `json:"start_time"`
	Hours       float64 `json:"hours"`
	Type        int     `json:"type"`
	CreateTime  int64   `json:"create_time"`
	Description string  `json:"description"`
}

type ThirdNotification struct {
	ProjectID  string `json:"project_id"`
	ValueType  int    `json:"value_type"`
	Value      string `json:"value"`
	ConfigType string `json:"config_type"`
}
