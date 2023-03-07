package common

import (
	"fmt"
	"sync"
)

const (
	ResourceTypeStringUser                   = "user"
	ResourceTypeStringUserGroup              = "user_group"
	ResourceTypeStringUserGroupMember        = "user_group_member"
	ResourceTypeStringDepartment             = "department"
	ResourceTypeStringUserDepartments        = "user_departments"
	ResourceTypeStringGlobalProjectRole      = "global_project_role"
	ResourceTypeStringGlobalProjectField     = "global_project_field"
	ResourceTypeStringIssueType              = "issue_type"
	ResourceTypeStringProject                = "project"
	ResourceTypeStringProjectIssueType       = "project_issue_type"
	ResourceTypeStringProjectRole            = "project_role"
	ResourceTypeStringProjectRoleMember      = "project_role_member"
	ResourceTypeStringGlobalPermission       = "global_permission"
	ResourceTypeStringProjectPermission      = "project_permission"
	ResourceTypeStringProjectFieldValue      = "project_field_value"
	ResourceTypeStringTaskStatus             = "task_status"
	ResourceTypeStringTaskField              = "task_field"
	ResourceTypeStringTaskFieldOption        = "task_field_option"
	ResourceTypeStringIssueTypeField         = "issue_type_field"
	ResourceTypeStringIssueTypeLayout        = "issue_type_layout"
	ResourceTypeStringProjectIssueTypeField  = "project_issue_type_field"
	ResourceTypeStringProjectIssueTypeLayout = "project_issue_type_layout"
	ResourceTypeStringProjectSprintField     = "project_sprint_field"
	ResourceTypeStringPriority               = "priority"
	ResourceTypeStringConfig                 = "config"
	ResourceTypeStringTaskLinkType           = "task_link_type"
	ResourceTypeStringWorkflow               = "workflow"
	ResourceTypeStringSprint                 = "sprint"
	ResourceTypeStringSprintFieldValue       = "sprint_field_value"
	ResourceTypeStringTask                   = "task"
	ResourceTypeStringTaskFieldValue         = "task_field_value"
	ResourceTypeStringTaskWatcher            = "task_watcher"
	ResourceTypeStringTaskWorkLog            = "task_work_log"
	ResourceTypeStringTaskComment            = "task_comment"
	ResourceTypeStringTaskRelease            = "task_release"
	ResourceTypeStringTaskLink               = "task_link"
	ResourceTypeStringNotification           = "notification"
	ResourceTypeStringTaskAttachmentTmp      = "task_attachment_tmp"
	ResourceTypeStringTaskAttachment         = "task_attachment"
	ResourceTypeStringChangeItem             = "change_item"
)

var (
	MapOutputFile = []string{
		ResourceTypeStringUser,
		ResourceTypeStringUserGroup,
		ResourceTypeStringUserGroupMember,
		ResourceTypeStringDepartment,
		ResourceTypeStringUserDepartments,
		ResourceTypeStringGlobalProjectRole,
		ResourceTypeStringGlobalProjectField,
		ResourceTypeStringIssueType,
		ResourceTypeStringProject,
		ResourceTypeStringProjectIssueType,
		ResourceTypeStringProjectRole,
		ResourceTypeStringProjectRoleMember,
		ResourceTypeStringGlobalPermission,
		ResourceTypeStringProjectPermission,
		ResourceTypeStringProjectFieldValue,
		ResourceTypeStringTaskStatus,
		ResourceTypeStringTaskField,
		ResourceTypeStringTaskFieldOption,
		ResourceTypeStringIssueTypeField,
		ResourceTypeStringIssueTypeLayout,
		ResourceTypeStringProjectIssueTypeField,
		ResourceTypeStringProjectIssueTypeLayout,
		ResourceTypeStringProjectSprintField,
		ResourceTypeStringPriority,
		ResourceTypeStringConfig,
		ResourceTypeStringTaskLinkType,
		ResourceTypeStringWorkflow,
		ResourceTypeStringSprint,
		ResourceTypeStringSprintFieldValue,
		ResourceTypeStringTask,
		ResourceTypeStringTaskFieldValue,
		ResourceTypeStringTaskWatcher,
		ResourceTypeStringTaskWorkLog,
		ResourceTypeStringTaskComment,
		ResourceTypeStringTaskRelease,
		ResourceTypeStringTaskLink,
		ResourceTypeStringNotification,
		ResourceTypeStringTaskAttachmentTmp,
		ResourceTypeStringTaskAttachment,
		ResourceTypeStringChangeItem,
	}

	// MapResolveTime backup file size(byte) => expected resolve time（sec)
	MapResolveTime = map[int64]int64{
		1000: 3600,
		800:  3000,
		500:  2000,
		200:  1000,
		100:  500,
		50:   250,
		0:    60,
	}
)

var (
	ImportCacheMap = new(ImportCache)
)

type ImportCache struct {
	Map sync.Map
}

func (c *ImportCache) Get(cookie string) *Cache {
	res, ok := c.Map.Load(cookie)
	if !ok {
		panic(fmt.Sprintf("load sync map err, %s", cookie))
	}
	r := res.(*Cache)
	return r
}

func (c *ImportCache) Set(cookie string, input *Cache) {
	c.Map.Store(cookie, input)
}
