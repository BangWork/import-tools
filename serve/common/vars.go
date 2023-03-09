package common

import (
	"fmt"
	"sync"
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

// MapResolveTime backup file size(byte) => expected resolve time（sec)
var MapResolveTime = map[int64]int64{
	1000: 3600,
	800:  3000,
	500:  2000,
	200:  1000,
	100:  500,
	50:   250,
	0:    60,
}

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
)
