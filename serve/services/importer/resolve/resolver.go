package resolve

import (
	"github.com/bangwork/import-tools/serve/services/importer/types"
)

type ResourceResolver interface {
	// Global
	NextGlobalProjectField() ([]byte, error)
	NextTaskStatus() ([]byte, bool, error)
	NextTaskField() ([]byte, error)
	NextTaskFieldOption() ([]byte, error)
	NextPriority() ([]byte, error)
	NextTaskLinkType() ([]byte, error)
	NextGlobalProjectRole() ([]byte, error)
	NextGlobalPermission() ([]byte, error)
	NextIssueType() ([]byte, error)
	NextIssueTypeLayout() ([]byte, error)
	NextIssueTypeField() ([]byte, error)
	NextWorkflow() ([]byte, error)

	Config() ([]byte, error) // 配置

	// User
	NextUser() ([]byte, error)
	NextUserGroup() ([]byte, error)
	NextUserGroupMember() ([]byte, error)

	// Project
	NextProject() ([]byte, error)
	NextProjectFieldValue() ([]byte, error)
	NextProjectRole() ([]byte, error)
	NextProjectPermission() ([]byte, error)
	NextProjectRoleMember() ([]byte, error)
	NextProjectIssueType() ([]byte, error)
	NextProjectIssueTypeLayout() ([]byte, error)
	NextProjectIssueTypeField() ([]byte, error)
	NextNotification() ([]byte, error)

	// Issue、Sprint
	NextSprint() ([]byte, error)
	NextTask() ([]byte, error)
	NextChangeItem() ([]byte, bool, error)
	NextTaskFieldValue() ([]byte, bool, error)
	NextTaskWatcher() ([]byte, bool, error)
	NextTaskWorkLog() ([]byte, bool, error)
	NextTaskComment() ([]byte, bool, error)
	NextTaskLink() ([]byte, error)
	NextTaskRelease() ([]byte, error)

	NextTaskAttachment() ([]byte, bool, error)

	TotalAttachmentSize() int64
	PrepareResolve() error
	Clear() error
	ServerID() string
}

type DefaultResolver struct{}

func (p *DefaultResolver) NextProject() (*ThirdProject, error) {
	return nil, nil
}
func (p *DefaultResolver) NextTask() (*ThirdTask, error) {
	return nil, nil
}
func (p *DefaultResolver) NextUser() (*ThirdUser, error) {
	return nil, nil
}
func (p *DefaultResolver) Clear() error {
	return nil
}

type ResolverFactory interface {
	CreateResolver(importTask *types.ImportTask) (ResourceResolver, error)
	InitImportFile(importTask *types.ImportTask) (ResourceResolver, error)
}
