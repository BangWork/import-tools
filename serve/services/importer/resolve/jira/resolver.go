package jira

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/gin-contrib/i18n"
	"github.com/juju/errors"

	"github.com/bangwork/import-tools/serve/common"
	commonModel "github.com/bangwork/import-tools/serve/models/common"
	fieldModel "github.com/bangwork/import-tools/serve/models/field"
	"github.com/bangwork/import-tools/serve/models/issuetype"
	"github.com/bangwork/import-tools/serve/models/layout/block"
	"github.com/bangwork/import-tools/serve/models/objectlinktype"
	"github.com/bangwork/import-tools/serve/models/permission"
	"github.com/bangwork/import-tools/serve/models/project"
	roleModel "github.com/bangwork/import-tools/serve/models/role"
	statusModel "github.com/bangwork/import-tools/serve/models/status"
	userModel "github.com/bangwork/import-tools/serve/models/user"
	utilsModel "github.com/bangwork/import-tools/serve/models/utils"
	"github.com/bangwork/import-tools/serve/services"
	"github.com/bangwork/import-tools/serve/services/cache"
	"github.com/bangwork/import-tools/serve/services/importer/constants"
	"github.com/bangwork/import-tools/serve/services/importer/resolve"
	"github.com/bangwork/import-tools/serve/services/importer/types"
	"github.com/bangwork/import-tools/serve/utils"
	utils2 "github.com/bangwork/import-tools/serve/utils"
	"github.com/bangwork/import-tools/serve/utils/timestamp"
	"github.com/bangwork/import-tools/serve/utils/unique"
)

type JiraResolver struct {
	totalAttachmentSize     int64
	serverID                string
	importTask              *types.ImportTask
	daysPerWeek             string
	hoursPerDay             string
	jiraSubTaskLinkID       string
	jiraSprintCustomFieldID string
	resolveResult           cache.ResolveResult

	beganMap map[string]bool

	resolve.DefaultResolver

	jiraTaskIDMap                           map[string]struct{}
	jiraTaskKeyIDMap                        map[string]string // jira issue key => jira issue id
	jiraUserEmailMap                        map[string]string // jira email => jira user id
	jiraUserNameIDMap                       map[string]string // jira user name => jira user id
	jiraUIDNameMap                          map[string]string // jira user id => jira user name
	jiraProjectIDNameMap                    map[string]string // jira project id => jira project name
	jiraProjectIDAssignIDMap                map[string]string // jira project id => jira project assign id
	jiraUserGroupNameIDMap                  map[string]string // jira user group name => jira user group id
	jiraUserGroupIDNameMap                  map[string]string // jira user group id => jira user group name
	jiraGlobalProjectRoleNameIDMap          map[string]string // jira role name => jira user group id
	jiraUserGroupIDMembersMap               map[string][]string
	jiraFieldConfigSchemeProjectIDsMap      map[string][]string // jira scheme id => jira project ids
	jiraIssueParentIDMap                    map[string]string   // jira sub issue id => jira parent issue id
	jiraCustomFieldTypeMap                  map[string]int      // jira issue field id => field type
	jiraIssueSprintMap                      map[string]string   // jira issue id => jira sprint id
	jiraSprintIDProjectIDMap                map[string]string   // jira sprint id => jira project id
	jiraWorkflowSchemeProjectAssoc          map[string][]string // scheme id => project ids
	jiraNotificationSchemeProjectAssoc      map[string][]string // scheme id => project ids
	jiraSprintColumnNames                   []string            // jira sprint field
	jiraCustomFieldCascadingSelectMap       map[string]string
	jiraCustomFieldCascadingSelectOptionMap map[string]map[string]string
	jiraCustomFieldVersionMap               map[string]string
	jiraReleaseVersionMap                   map[string]string
	jiraCustomFieldProjectMap               map[string]string
	jiraCustomFieldMultiGroupMap            map[string]string
	jiraIssueComponent                      map[string][]string
	jiraComponentNameMap                    map[string]string
	jiraIssueFixVersion                     map[string][]string
	jiraIssueAffectsVersion                 map[string][]string
	jiraEventTypeIDNamesMap                 map[int][]string    // jira event type id => event type names
	jiraPermissionProjectAssociation        map[string][]string // jira scheme id => jira project ids
	jiraProjectTypeOptionMap                map[string]bool     // jira project type option value
	jiraProjectTypeIDMap                    map[string]string   // jira project id => project type value（jira project type without id，use value as id）

	jiraGlobalProjectField                []*resolve.ThirdGlobalProjectField    //  project field config
	jiraProjectFieldValue                 []*resolve.ThirdProjectFieldValue     //  project field config
	jiraFieldLayoutIDWithLayout           map[string]*fieldLayout               // jira field layout id => jira field layout
	jiraProjectWithFieldLayoutScheme      map[string]string                     // jira project id => jira fieldLayoutScheme id
	jiraFieldLayoutSchemeIDWithEntity     map[string][]*fieldLayoutSchemeEntity // fieldLayoutSchemeID => FieldLayoutSchemeEntity
	jiraIssueTypeScreenSchemeIDWithEntity map[string][]*issueTypeScreenSchemeEntity
	jiraProjectWithIssueTypeMap           map[string][]string // jira project id => issue type ids
	jiraProjectWithIssueTypeScreenScheme  map[string]string   // project id => issue type screen scheme id
	jiraFieldScreenSchemeIDWithName       map[string]string   // id => name
	jiraFieldScreenSchemeIDWithItem       map[string][]*fieldScreenSchemeItem
	jiraFieldScreenIDWithTabIds           map[string][]string
	jiraProjectIssueTypeNameMap           map[string]bool // jira project:jira issue type id => bool
	jiraFieldLayoutIDWithItem             map[string][]*fieldLayoutItem
	jiraFieldScreenTabIDWithItem          map[string][]*fieldScreenLayoutItem
	jiraProjectIssueTypeFieldScreen       map[string]string
	jiraFieldConfigSchemeMap              map[string]string
	jiraChangeGroupMap                    map[string]*changeGroup
	jiraCategoryProjectAssociation        map[string]string
	jiraProjectIDOriginalKeyMap           map[string]string // jira project id => jira project original key
	jiraProjectIssueTypeMap               map[string]struct{}

	jiraProjectIssueTypeLayoutSlice     []*resolve.ThirdProjectIssueTypeLayout
	jiraProjectIssueTypeWithFieldsSlice []*resolve.ThirdProjectIssueTypeWithFields
	jiraIssueTypeLayoutSlice            []*resolve.ThirdIssueTypeLayout
	jiraProjectRolesSlice               []*projectRoleStruct
	jiraProjectMembersSlice             []*projectMemberStruct
	jiraWorkflowsSlice                  []*resolve.ThirdWorkflow
	jiraGlobalPermissionSlice           []*resolve.ThirdGlobalPermission
	jiraProjectPermissionSlice          []*resolve.ThirdProjectPermission
	jiraProjectIssueTypeSlice           []*resolve.ThirdProjectIssueType
	jiraNotifications                   []*resolve.ThirdNotification
	jiraReleaseVersionSlice             []*etree.Element
	jiraTaskReleaseSlice                []*resolve.ThirdTaskRelease

	projectPermissions         []string
	projectPermissionProjectID []string

	taskReleaseIndex            int
	versionCommentIndex         int
	versionIndex                int
	projectIndex                int
	projectIssueTypeIndex       int
	projectFieldValueIndex      int
	globalProjectFieldIndex     int
	projectIssueTypeFieldIndex  int
	projectIssueTypeLayoutIndex int
	issueTypeLayoutIndex        int
	projectIssueTypeInnerIndex  int
	projectPermissionIndex      int
	globalPermissionIndex       int
	customFieldIndex            int
	customFieldOptionIndex      int
	workdayIndex                int
	projectRoleIndex            int
	projectRoleMemberIndex      int
	sprintIndex                 int
	issueLinkTypeIndex          int
	customFieldValueIndex       int
	workflowIndex               int
	notificationIndex           int
	notificationInnerIndex      int
	workLogIndex                int
	issueLinkIndex              int
	issueTypeFieldIndex         int

	root                       *etree.Element
	sprintRoot                 *etree.Element
	jiraCustomField            []*resolve.ThirdTaskField
	jiraCustomFieldOptions     []*resolve.ThirdTaskFieldOption
	jiraIssueLinkType          []*etree.Element
	jiraCustomFieldValues      []*resolve.ThirdTaskFieldValue
	jiraIssueLinks             []*etree.Element
	jiraSprints                []*etree.Element
	jiraNodeAssociations       []*etree.Element
	jiraUserAssociations       []*etree.Element
	jiraProjectElements        []*etree.Element
	jiraWorkLogs               []*resolve.ThirdTaskWorkLog
	jiraComments               []*etree.Element
	tagFilesMap                map[string]*resolve.XmlScanner
	mapTagFilePath             map[string]string
	IssueParentFilePath        string
	IssueChildFilePath         string
	IssueParentReader          *resolve.XmlScanner
	IssueChildReader           *resolve.XmlScanner
	mapIssueKeyWithAttachments map[string]map[string]string
	mapIssueIDWithKey          map[string]string
	mapIssueIDWithOriginalKey  map[string]string
	zipFile                    *zip.ReadCloser

	issueTypeDetailTypeMap   map[string]int        // issueTypeID: builtinType
	jiraProjectRoleIDDataMap map[string][]roleData // jira 项目角色 id => jira 项目角色配的用户组和单个用户

	nowTimeString string
}

type JiraResolverFactory struct{}

var (
	thirdIssueTypeFields = []*resolve.ThirdIssueTypeField{
		{
			FieldID:             customFieldReleaseStartDate,
			IssueTypeDetailType: issuetype.DetailTypePublish,
		},
	}
	customFields = []string{
		"Fix Version/s",
		"Affects Version/s",
		"Labels",
		"Resolution",
		"ID",
		"Component/s",
		fieldModel.PublishVersionFieldUUID,
		"Environment",
	}
	tabFields = map[string]string{
		"issuelinks":   block.TabLabelRelatedContent,
		"attachment":   block.TabLabelFiles,
		"timetracking": block.TabLabelAssessManhour,
		"worklog":      block.TabLabelAssessManhour,
	}
	customFieldsMapping = map[string]string{
		//"fixVersions": "Fix Version/s",
		"versions":                  "Affects Version/s",
		"labels":                    "Labels",
		"resolution":                "Resolution",
		"id":                        "ID",
		"components":                "Component/s",
		customFieldReleaseStartDate: customFieldReleaseStartDate,
		"environment":               "Environment",
	}

	// jira视图字段替换为ONES系统字段
	issueTypeLayoutJiraOnesFieldMap = map[string]string{
		"fixVersions": fieldModel.PublishVersionFieldUUID,
	}

	// jira工作项
	fixedFieldsMapping = map[string]string{
		"assignee":   "field004",
		"archivedby": "",
		//"environment":  "",
		"security":     "",
		"archiveddate": "",
		"issuetype":    "field007",
		"summary":      "field001",
		"description":  "field016",
		"priority":     "field012",
		"duedate":      "field013",
		"comment":      "",
		"reporter":     "",
		"sprint":       "field011",
	}

	eventTypeMap = map[string][]string{
		"Issue Created": {"create_task"},
		"Issue Updated": {
			"update_task_status",
			"update_task_priority",
			"update_task_title",
			"update_task_description",
			"set_task_deadline",
			"update_task_sprint",
			"update_task_related_task",
			"related_plan_case",
			"related_test_case",
			"task_related_testcase_plan",
			"upload_task_file",
			"update_issue_type",
			"update_std_to_sub_issue_type",
		},
		"Issue Assigned":        {"update_task_assign"},
		"Issue Commented":       {"update_task_message"},
		"Work Logged On Issue":  {"update_task_access_manhour", "update_task_record_manhour"},
		"Issue Worklog Updated": {"update_task_record_manhour"},
		"Issue Worklog Deleted": {"update_task_record_manhour"},
		"Generic Event":         {"update_task_other_property"},
	}

	subscriberMap = map[string]map[string]string{
		"Current_Assignee": {"type": "role", "value": "task_assign"},
		"Current_Reporter": {"type": "role", "value": "task_assign"},
		"Project_Lead":     {"type": "role", "value": "project_assign"},
		"Single_User":      {"type": "single_user_uuids", "value": ""},
		"Group_Dropdown":   {"type": "group_uuids", "value": ""},
		"Project_Role":     {"type": "role_uuids", "value": ""},
		"All_Watchers":     {"type": "role", "value": "task_watchers"},
	}

	notificationTypeMap = map[string]int{
		"single_user_uuids": constants.NotificationValueTypeUserID,
		"group_uuids":       constants.NotificationValueTypeUserGroupID,
		"role_uuids":        constants.NotificationValueTypeGlobalProjectRoleID,
		"role":              constants.NotificationValueTypeRole,
	}

	userDomainMap = map[string]map[string]string{
		// 应用程序访问权
		applicationRole: {constants.ProjectContext: "everyone", constants.IssueTypeContext: "everyone"},
		// 项目角色
		projectRole: {constants.ProjectContext: "role", constants.IssueTypeContext: "role"},
		// 用户组
		userGroup: {constants.ProjectContext: "group", constants.IssueTypeContext: "group"},
		// 单个用户
		singleUser: {constants.ProjectContext: "single_user", constants.IssueTypeContext: "single_user"},
		// 当前经办人
		assignee: {constants.ProjectContext: "project_assign", constants.IssueTypeContext: "task_assign"},
		// 报告人
		reporter: {constants.ProjectContext: "project_assign", constants.IssueTypeContext: "task_assign"},
		// 项目主管
		lead: {constants.ProjectContext: "project_assign", constants.IssueTypeContext: "project_assign"},
	}

	globalPermissionMap = map[string][]int{
		"ADMINISTER":   {permission.AdministerTeam, permission.AdministerDo},
		"SYSTEM_ADMIN": {permission.AdministerTeam, permission.AdministerDo},
	}
	needConvertPermissionMap = map[string]bool{
		ManageSprintsPermission: true,
		AdministerProjects:      true,
		BrowseProjects:          true,
		CreateIssues:            true,
		DeleteIssues:            true,
		AssignableUser:          true,
		ManageWatchers:          true,

		EditIssues:   true,
		AssignIssues: true,
		LinkIssues:   true,

		TransitionIssues: true,
		CloseIssues:      true,
		ResolveIssues:    true,

		EditAllWorkLogs:   true,
		DeleteAllWorkLogs: true,

		EditOwnWorkLogs:   true,
		DeleteOwnWorkLogs: true,
	}

	projectPermissionMap = map[string][]string{
		ManageSprintsPermission: {"manage_sprints", "be_assigned_to_sprint"},
		AdministerProjects:      {"manage_project", "view_project_reports", "manage_pipelines"},
		BrowseProjects:          {"browse_project", "view_tasks", "export_tasks"},
		CreateIssues:            {"create_tasks"},
		DeleteIssues:            {"delete_tasks"},
		AssignableUser:          {"be_assigned"},
		ManageWatchers:          {"update_task_watchers"},
		EditIssues:              {"update_tasks"},
		TransitionIssues:        {"transit_tasks"},
		EditAllWorkLogs:         {"manage_task_assess_manhour", "manage_task_record_manhours"},
		EditOwnWorkLogs:         {"manage_task_own_record_manhours", "manage_task_own_assess_manhour"},
	}
	needMergePermissionMap = map[string][]string{
		EditIssues:        {EditIssues, AssignIssues, LinkIssues}, // 取数组里的第一个作为标识，这里是EDIT_ISSUES
		AssignIssues:      {EditIssues, AssignIssues, LinkIssues},
		LinkIssues:        {EditIssues, AssignIssues, LinkIssues},
		TransitionIssues:  {TransitionIssues, CloseIssues, ResolveIssues},
		CloseIssues:       {TransitionIssues, CloseIssues, ResolveIssues},
		ResolveIssues:     {TransitionIssues, CloseIssues, ResolveIssues},
		EditAllWorkLogs:   {EditAllWorkLogs, DeleteAllWorkLogs},
		DeleteAllWorkLogs: {EditAllWorkLogs, DeleteAllWorkLogs},
		EditOwnWorkLogs:   {EditOwnWorkLogs, DeleteOwnWorkLogs},
		DeleteOwnWorkLogs: {EditOwnWorkLogs, DeleteOwnWorkLogs},
	}

	permissionContextMap = map[string]string{
		"manage_sprints":                  constants.ProjectContext,
		"be_assigned_to_sprint":           constants.ProjectContext,
		"manage_project":                  constants.ProjectContext,
		"browse_project":                  constants.ProjectContext,
		"view_project_reports":            constants.ProjectContext,
		"manage_pipelines":                constants.ProjectContext,
		"view_tasks":                      constants.IssueTypeContext,
		"export_tasks":                    constants.IssueTypeContext,
		"create_tasks":                    constants.IssueTypeContext,
		"delete_tasks":                    constants.IssueTypeContext,
		"update_tasks":                    constants.IssueTypeContext,
		"be_assigned":                     constants.IssueTypeContext,
		"transit_tasks":                   constants.IssueTypeContext,
		"update_task_watchers":            constants.IssueTypeContext,
		"manage_task_assess_manhour":      constants.IssueTypeContext,
		"manage_task_record_manhours":     constants.IssueTypeContext,
		"manage_task_own_record_manhours": constants.IssueTypeContext,
		"manage_task_own_assess_manhour":  constants.IssueTypeContext,
	}
)

func (p *JiraResolverFactory) CreateResolver(importTask *types.ImportTask) (resolve.ResourceResolver, error) {
	resolver := &JiraResolver{
		importTask: importTask,
	}
	resolver.initAttributes()
	resolver.initIssueTypeMap()
	for tag, file := range importTask.MapFilePath {
		fi, err := os.Open(file)
		if err != nil {
			log.Printf("open file fail: %s", err)
		}
		scanner := resolve.NewXmlScanner(fi, entityRootTag)
		resolver.tagFilesMap[tag] = scanner
	}
	resolver.mapTagFilePath = importTask.MapFilePath
	if err := resolver.loadElements(); err != nil {
		return nil, err
	}
	if err := resolver.processAttachment(); err != nil {
		return nil, err
	}
	return resolver, nil
}

func (p *JiraResolverFactory) InitImportFile(importTask *types.ImportTask) (resolve.ResourceResolver, error) {
	resolver := &JiraResolver{
		importTask: importTask,
	}
	err := resolver.InitImportFile()
	if err != nil {
		return nil, err
	}
	log.Println("start calculate attachments size")
	attachmentSize, err := utils.GetDirSize(importTask.AttachmentsPath)
	if err != nil {
		return nil, err
	}
	log.Println("end calculate attachments size")
	resolver.resolveResult.AttachmentSize = attachmentSize
	for tag, file := range resolver.mapTagFilePath {
		fi, err := os.Open(file)
		if err != nil {
			log.Printf("open file fail: %s", err)
		}
		scanner := resolve.NewXmlScanner(fi, entityRootTag)
		resolver.tagFilesMap[tag] = scanner
	}

	return resolver, nil
}

func (p *JiraResolver) PrepareResolve() error {
	if err := p.prepareProjectIssueTypeMap(); err != nil {
		return err
	}
	return nil
}

func (p *JiraResolver) prepareProjectIssueTypeMap() error {
	if err := p.resolvePrepareConfigurationContext(); err != nil {
		return err
	}
	projectIssueTypeMap, err := p.resolvePrepareProjectIssueType()
	if err != nil {
		return err
	}
	cacheInfo, err := cache.GetCacheInfo(p.importTask.Key)
	if err != nil {
		return err
	}
	cacheInfo.ProjectIssueTypeMap = projectIssueTypeMap
	return cache.SetCacheInfo(p.importTask.Key, cacheInfo)
}

func (p *JiraResolver) initIssueTypeMap() {
	p.issueTypeDetailTypeMap = make(map[string]int)
	for _, tm := range p.importTask.BuiltinIssueTypes {
		if tm.IssueTypeID != "" {
			p.issueTypeDetailTypeMap[tm.IssueTypeID] = tm.BuiltinType
		}
	}
}

func (p *JiraResolver) InitImportFile() error {
	file, err := os.Open(p.importTask.LocalFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if services.StopResolveSignal {
		return nil
	}

	entityFileReader, activeObjectsFileReader, err := p.readersFromFile(file.Name())
	if err != nil {
		return err
	}
	defer entityFileReader.Close()
	defer activeObjectsFileReader.Close()

	tagFilesMap, mapFilePathMap, handler, err := processedEntityFile(p.importTask, entityFileReader)
	if err != nil {
		return err
	}

	if services.StopResolveSignal {
		return nil
	}
	serverID, daysPerWeek, hoursPerDay := handler.ResolveResult.JiraServerID, handler.daysPerWeek, handler.hoursPerDay
	p.resolveResult = handler.ResolveResult

	objectFilePath, err := processedObjectFile(p.importTask, activeObjectsFileReader)
	if err != nil {
		return err
	}
	if services.StopResolveSignal {
		return nil
	}
	mapFilePathMap[common.TagObjectFile] = objectFilePath

	p.tagFilesMap = tagFilesMap
	p.serverID = serverID
	p.daysPerWeek = daysPerWeek
	p.hoursPerDay = hoursPerDay
	p.mapTagFilePath = mapFilePathMap
	return nil
}

func (p *JiraResolver) setCache() error {
	info, err := cache.GetCacheInfo(p.importTask.Key)
	if err != nil {
		return err
	}
	result := p.resolveResult
	info.ResolveResult = &result
	info.MapFilePath = p.mapTagFilePath
	info.HoursPerDay = p.hoursPerDay
	info.DaysPerWeek = p.daysPerWeek
	info.ResolveStatus = common.ResolveStatusDone
	info.ResolveDoneTime = time.Now().Unix()
	return cache.SetCacheInfo(p.importTask.Key, info)
}

func (p *JiraResolver) initAttributes() {
	p.jiraTaskIDMap = make(map[string]struct{})
	p.tagFilesMap = make(map[string]*resolve.XmlScanner)
	p.jiraProjectIDOriginalKeyMap = make(map[string]string)
	p.jiraProjectIssueTypeMap = make(map[string]struct{})
	p.jiraProjectIssueTypeSlice = make([]*resolve.ThirdProjectIssueType, 0)
	p.jiraNotifications = make([]*resolve.ThirdNotification, 0)
	p.jiraReleaseVersionSlice = make([]*etree.Element, 0)
	p.jiraTaskReleaseSlice = make([]*resolve.ThirdTaskRelease, 0)
	p.jiraProjectElements = make([]*etree.Element, 0)
	p.jiraProjectFieldValue = make([]*resolve.ThirdProjectFieldValue, 0)
	p.jiraGlobalProjectField = make([]*resolve.ThirdGlobalProjectField, 0)
	p.jiraProjectTypeOptionMap = make(map[string]bool)
	p.jiraProjectTypeIDMap = make(map[string]string)
	p.jiraTaskKeyIDMap = make(map[string]string)
	p.jiraCategoryProjectAssociation = make(map[string]string)
	p.jiraChangeGroupMap = make(map[string]*changeGroup)
	p.jiraFieldConfigSchemeMap = make(map[string]string)
	p.jiraProjectIssueTypeLayoutSlice = make([]*resolve.ThirdProjectIssueTypeLayout, 0)
	p.jiraProjectIssueTypeWithFieldsSlice = make([]*resolve.ThirdProjectIssueTypeWithFields, 0)
	p.jiraProjectIssueTypeFieldScreen = make(map[string]string)
	p.jiraIssueTypeLayoutSlice = make([]*resolve.ThirdIssueTypeLayout, 0)
	p.jiraFieldScreenTabIDWithItem = make(map[string][]*fieldScreenLayoutItem)
	p.jiraFieldLayoutIDWithItem = make(map[string][]*fieldLayoutItem)
	p.jiraProjectIssueTypeNameMap = make(map[string]bool)
	p.jiraIssueTypeScreenSchemeIDWithEntity = make(map[string][]*issueTypeScreenSchemeEntity)
	p.jiraFieldScreenIDWithTabIds = make(map[string][]string)
	p.jiraFieldScreenSchemeIDWithItem = make(map[string][]*fieldScreenSchemeItem)
	p.jiraFieldScreenSchemeIDWithName = make(map[string]string)
	p.jiraProjectWithIssueTypeScreenScheme = make(map[string]string)
	p.jiraProjectWithIssueTypeMap = make(map[string][]string)
	p.jiraFieldLayoutSchemeIDWithEntity = make(map[string][]*fieldLayoutSchemeEntity)
	p.jiraProjectWithFieldLayoutScheme = make(map[string]string)
	p.jiraFieldLayoutIDWithLayout = make(map[string]*fieldLayout)
	p.projectPermissions = make([]string, 0)
	p.projectPermissionProjectID = make([]string, 0)
	p.beganMap = make(map[string]bool)
	p.jiraUserEmailMap = make(map[string]string)
	p.jiraUserNameIDMap = make(map[string]string)
	p.jiraUserGroupNameIDMap = make(map[string]string)
	p.jiraUserGroupIDNameMap = make(map[string]string)
	p.jiraGlobalProjectRoleNameIDMap = make(map[string]string)
	p.jiraUIDNameMap = make(map[string]string)
	p.jiraUserGroupIDMembersMap = make(map[string][]string)
	p.jiraFieldConfigSchemeProjectIDsMap = make(map[string][]string)
	p.jiraProjectIDNameMap = make(map[string]string)
	p.jiraIssueParentIDMap = make(map[string]string)
	p.jiraCustomFieldTypeMap = make(map[string]int)
	p.jiraWorkflowsSlice = make([]*resolve.ThirdWorkflow, 0)
	p.jiraGlobalPermissionSlice = make([]*resolve.ThirdGlobalPermission, 0)
	p.jiraIssueSprintMap = make(map[string]string)
	p.jiraCustomFieldValues = make([]*resolve.ThirdTaskFieldValue, 0)
	p.jiraWorkflowSchemeProjectAssoc = make(map[string][]string)
	p.jiraNotificationSchemeProjectAssoc = make(map[string][]string)
	p.jiraSprintColumnNames = make([]string, 0)
	p.jiraSprintIDProjectIDMap = make(map[string]string)
	p.jiraProjectIDAssignIDMap = make(map[string]string)
	p.jiraCustomFieldCascadingSelectMap = make(map[string]string)
	p.jiraCustomFieldCascadingSelectOptionMap = make(map[string]map[string]string)
	p.jiraCustomFieldVersionMap = make(map[string]string)
	p.jiraReleaseVersionMap = make(map[string]string)
	p.jiraCustomFieldProjectMap = make(map[string]string)
	p.jiraCustomFieldMultiGroupMap = make(map[string]string)
	p.jiraIssueComponent = make(map[string][]string)
	p.jiraComponentNameMap = make(map[string]string)
	p.jiraIssueFixVersion = make(map[string][]string)
	p.jiraIssueAffectsVersion = make(map[string][]string)
	p.jiraEventTypeIDNamesMap = make(map[int][]string)
	p.jiraWorkLogs = make([]*resolve.ThirdTaskWorkLog, 0)
	p.jiraPermissionProjectAssociation = make(map[string][]string)
	p.mapIssueKeyWithAttachments = make(map[string]map[string]string)
	p.mapIssueIDWithKey = make(map[string]string)
	p.mapIssueIDWithOriginalKey = make(map[string]string)
	p.jiraProjectRoleIDDataMap = make(map[string][]roleData)
	p.nowTimeString = timestamp.NowTimeString()
}

func (p *JiraResolver) loadElements() error {
	customFields, err := p.getCustomField()
	if err != nil {
		return err
	}
	p.jiraCustomField = customFields

	customFieldOptions, err := p.getCustomFieldOption()
	if err != nil {
		return err
	}
	p.jiraCustomFieldOptions = customFieldOptions

	jiraIssueLinkType, err := p.getIssueLinkTypes()
	if err != nil {
		return err
	}
	p.jiraIssueLinkType = jiraIssueLinkType

	mapVersion, err := p.getVersion()
	if err != nil {
		return err
	}
	p.jiraReleaseVersionMap = mapVersion

	jiraNodeAssociations, err := p.getNodeAssociation()
	if err != nil {
		return err
	}
	p.jiraNodeAssociations = jiraNodeAssociations

	issueLinks, err := p.getIssueLinks()
	if err != nil {
		return err
	}
	p.jiraIssueLinks = issueLinks
	customFieldValues, err := p.getCustomFieldValues()
	if err != nil {
		return err
	}
	p.jiraCustomFieldValues = append(p.jiraCustomFieldValues, customFieldValues...)

	p.prepareProject()
	if err := p.preprocessIssues(); err != nil {
		return err
	}
	objectFilePath := p.mapTagFilePath[common.TagObjectFile]
	sprintDoc := etree.NewDocument()
	if err := sprintDoc.ReadFromFile(objectFilePath); err != nil {
		return err
	}
	p.sprintRoot = sprintDoc.Root()
	p.jiraSprints = p.sprintRoot.FindElements("./row")

	if err := p.prepareApplicationUser(); err != nil {
		return err
	}
	if err := p.prepareConfigurationContext(); err != nil {
		return err
	}
	p.prepareSprintColumn()
	if err := p.prepareEventTypeIDNamesMap(); err != nil {
		return err
	}
	if err := p.prepareFieldLayout(); err != nil {
		return err
	}
	if err := p.prepareFieldLayoutSchemeEntity(); err != nil {
		return err
	}
	if err := p.prepareFieldScreenScheme(); err != nil {
		return err
	}
	if err := p.prepareFieldScreenSchemeItem(); err != nil {
		return err
	}
	if err := p.prepareFieldConfigScheme(); err != nil {
		return err
	}
	if err := p.prepareIssueTypeScreenSchemeEntity(); err != nil {
		return err
	}
	if err := p.prepareFieldLayoutItem(); err != nil {
		return err
	}
	if err := p.prepareFieldScreenTab(); err != nil {
		return err
	}
	if err := p.prepareFieldScreenLayoutItem(); err != nil {
		return err
	}

	p.prepareGlobalProjectField()
	p.prepareProjectType()
	p.prepareProjectCategory()

	return nil
}

func (p *JiraResolver) prepareGlobalProjectField() {
	globalProjectFields := []*resolve.ThirdGlobalProjectField{
		{
			Base: resolve.Base{
				ResourceID: ResourceIDProjectURL,
			},
			Name: ResourceIDProjectURL,
			Type: fieldModel.FieldTypeText,
		},
		{
			Base: resolve.Base{
				ResourceID: ResourceIDProjectDescription,
			},
			Name: ResourceIDProjectDescription,
			Type: fieldModel.FieldTypeMultiLineText,
		},
	}

	p.jiraGlobalProjectField = append(p.jiraGlobalProjectField, globalProjectFields...)
}

func (p *JiraResolver) getUIDs(projectID, usersType, parameter string) []string {
	if usersType == singleUser {
		uid, ok := p.jiraUserNameIDMap[parameter]
		if !ok {
			log.Printf("uid not found, usersType: %s, parameter: %s", usersType, parameter)
			return []string{}
		}
		return []string{uid}
	}
	if usersType == userGroup {
		groupID, ok := p.jiraUserGroupNameIDMap[parameter]
		if !ok {
			log.Printf("usergroup id not found, usersType: %s, parameter: %s", usersType, parameter)
			return []string{}
		}
		return p.jiraUserGroupIDMembersMap[groupID]
	}
	if usersType == projectRole {
		roleDatas, ok := p.jiraProjectRoleIDDataMap[projectID+parameter]
		if !ok {
			log.Printf("role not found, usersType: %s, parameter: %s", usersType, parameter)
			return []string{}
		}
		var uids []string
		for _, roleData := range roleDatas {
			if roleData.roleType == roleTypeUser {
				uids = append(uids, p.getUIDs(projectID, singleUser, roleData.roleTypeParameter)...)
			} else if roleData.roleType == roleTypeGroup {
				uids = append(uids, p.getUIDs(projectID, userGroup, roleData.roleTypeParameter)...)
			}
		}
		return uids
	}
	return []string{}
}

func (p *JiraResolver) prepareProjectPermission() error {
	schemePermissionKeyPermissionRecordsMap := make(map[string]map[string][]permissionRecord)
	for {
		element, err := p.nextElement("SchemePermissions")
		if element == nil || err != nil {
			break
		}
		id := getAttributeValue(element, "id")
		scheme := getAttributeValue(element, "scheme")
		spType := getAttributeValue(element, "type")
		parameter := getAttributeValue(element, "parameter")
		permissionKey := getAttributeValue(element, "permissionKey")
		if scheme == "" {
			continue
		}

		if !needConvertPermissionMap[permissionKey] {
			continue
		}
		if schemePermissionKeyPermissionRecordsMap[scheme] == nil {
			schemePermissionKeyPermissionRecordsMap[scheme] = make(map[string][]permissionRecord)
		}
		schemePermissionKeyPermissionRecordsMap[scheme][permissionKey] = append(
			schemePermissionKeyPermissionRecordsMap[scheme][permissionKey], permissionRecord{
				id:            id,
				scheme:        scheme,
				permissionKey: permissionKey,
				spType:        spType,
				parameter:     parameter,
			})
	}
	var projectPermissionSlice []*resolve.ThirdProjectPermission
	for scheme, permissionKeyPermissionRecordsMap := range schemePermissionKeyPermissionRecordsMap {
		projectIDs, ok := p.jiraPermissionProjectAssociation[scheme]
		if !ok || len(projectIDs) == 0 {
			continue
		}
		haveMerge := make(map[string]bool)
		for permissionKey := range permissionKeyPermissionRecordsMap {
			var permissionRecords []permissionRecord
			jiraPermissions, need := needMergePermissionMap[permissionKey]
			mainPermissionKey := permissionKey
			var temp []*resolve.ThirdProjectPermission
			if need {
				mainPermissionKey = jiraPermissions[0]
				if haveMerge[mainPermissionKey] {
					// 多个jira权限转为ONES的一个权限，要避免重复合并
					continue
				}
				var toMerge = [][]permissionRecord{}
				for _, jiraPermission := range jiraPermissions {
					toMerge = append(toMerge, permissionKeyPermissionRecordsMap[jiraPermission])
				}
				temp = p.mergePermission(projectIDs, scheme, mainPermissionKey, toMerge...)
				haveMerge[mainPermissionKey] = true
			} else {
				permissionRecords = permissionKeyPermissionRecordsMap[permissionKey]
				temp = p.genProjectPermissionSlice(projectIDs, mainPermissionKey, permissionRecords)
			}
			if len(temp) == 0 {
				continue
			}
			projectPermissionSlice = append(projectPermissionSlice, temp...)
		}
	}
	// 对 projectPermissionSlice 去重
	ruleMap := make(map[string]bool)
	var uniqueProjectPermissionSlice []*resolve.ThirdProjectPermission
	for _, rule := range projectPermissionSlice {
		sign := fmt.Sprintf("%s:%s:%s:%s:%s", rule.ProjectID, rule.UserDomainType, rule.ContextType,
			rule.Permission, rule.UserDomainParam)
		if ruleMap[sign] {
			continue
		}
		ruleMap[sign] = true
		uniqueProjectPermissionSlice = append(uniqueProjectPermissionSlice, rule)
	}

	p.jiraProjectPermissionSlice = uniqueProjectPermissionSlice
	return nil
}

func (p *JiraResolver) mergePermission(projectIDs []string, scheme, permissionKey string,
	permissionRecordsToMerge ...[]permissionRecord) []*resolve.ThirdProjectPermission {
	// 这里是合并多个权限的权限记录，两个，三个
	var thirdProjectPermissions []*resolve.ThirdProjectPermission
	if len(permissionRecordsToMerge) == 0 {
		return thirdProjectPermissions
	}
	if len(permissionRecordsToMerge[0]) == 0 {
		return thirdProjectPermissions
	}
	// 规则1: 合并用户组、项目角色、单个用户的用户
	result1 := p.mergeUsergroupProjectRoleSingleUserRule(projectIDs, scheme, permissionKey, permissionRecordsToMerge...)
	thirdProjectPermissions = append(thirdProjectPermissions, result1...)
	// 规则2: 合并特殊成员域：报告人、当前经办人、项目负责人
	result2 := p.mergeSpecialUserDomain(projectIDs, scheme, permissionKey, permissionRecordsToMerge...)
	thirdProjectPermissions = append(thirdProjectPermissions, result2...)
	return thirdProjectPermissions
}

func (p *JiraResolver) mergeUsergroupProjectRoleSingleUserRule(projectIDs []string, scheme, permissionKey string,
	permissionRecordsToMerge ...[]permissionRecord) []*resolve.ThirdProjectPermission {
	var thirdProjectPermissions []*resolve.ThirdProjectPermission
	for _, projectID := range projectIDs {
		if p.ignoreProjectID(projectID) {
			continue
		}
		permissionRecords := p.mergeOneProject(projectID, scheme, permissionKey, permissionRecordsToMerge...)
		if len(permissionRecords) == 0 {
			// 规则3：为空，负责人兜底
			permissionRecords = []permissionRecord{
				{
					id:            utils2.UUID(), // 随机生成resourceID，不要和其他的重复
					scheme:        scheme,
					permissionKey: permissionKey,
					spType:        lead,
					parameter:     "",
				},
			}
		}
		thirdProjectPermissions = append(thirdProjectPermissions,
			p.genProjectPermissionSlice([]string{projectID}, permissionKey, permissionRecords)...)
	}
	return thirdProjectPermissions
}

func (p *JiraResolver) mergeOneProject(projectID string, scheme, permissionKey string,
	permissionRecordsToMerge ...[]permissionRecord) []permissionRecord {
	if len(permissionRecordsToMerge) == 0 {
		return nil
	}
	var newPermissionRecords []permissionRecord
	var toIntersect = [][]string{}
	groupOrRoleUIDsMap := make(map[string][]string)
	groupOrRolePermissionRecordMap := make(map[string]permissionRecord)
	for _, permissionRecords := range permissionRecordsToMerge {
		var tempUIDs []string
		for _, permissionRecord := range permissionRecords {
			if permissionRecord.spType != singleUser && permissionRecord.spType != projectRole &&
				permissionRecord.spType != userGroup {
				// 只处理这三种类型的
				continue
			}
			groupUIDs := p.getUIDs(projectID, permissionRecord.spType, permissionRecord.parameter)
			if permissionRecord.spType != singleUser {
				key := permissionRecord.spType + permissionRecord.parameter
				groupOrRoleUIDsMap[key] = utils2.UniqueNoNullSlice(groupUIDs...)
				groupOrRolePermissionRecordMap[key] = permissionRecord
			}
			tempUIDs = append(tempUIDs, groupUIDs...)
		}
		tempUIDs = utils2.UniqueNoNullSlice(tempUIDs...)
		toIntersect = append(toIntersect, tempUIDs)
	}
	uids := intersectUIDs(toIntersect...)
	if len(uids) == 0 {
		return newPermissionRecords
	}
	// 得到uids之后，看看是否能转为某个用户组
	uidSet := make(map[string]bool)
	for _, uid := range uids {
		uidSet[uid] = true
	}
	var allInGroupKeys []string
	// 要按人数从大到小排序
	items := make([]*utilsModel.Tuple2_String_Int, 0)
	for key, uids := range groupOrRoleUIDsMap {
		items = append(items, &utilsModel.Tuple2_String_Int{
			Ele_1: key,
			Ele_2: int64(len(uids)),
		})
	}
	itemsToSort := utilsModel.Tuple2_String_Int_Sorter(items)
	sort.Sort(sort.Reverse(itemsToSort))
	for _, item := range itemsToSort {
		allIn := true
		groupUIDs := groupOrRoleUIDsMap[item.Ele_1]
		for _, groupUID := range groupUIDs {
			if len(uidSet) == 0 || !uidSet[groupUID] {
				allIn = false
				break
			}
		}
		if allIn {
			allInGroupKeys = append(allInGroupKeys, item.Ele_1)
			// 调整uidSet
			for _, groupUID := range groupUIDs {
				delete(uidSet, groupUID)
			}
		}
	}
	for _, key := range allInGroupKeys {
		permissionRecord := groupOrRolePermissionRecordMap[key]
		permissionRecord.scheme = scheme
		permissionRecord.permissionKey = permissionKey
		newPermissionRecords = append(newPermissionRecords, permissionRecord)
	}
	for uid := range uidSet {
		// resourceID需要重新生成，不能和其他的重复
		newPermissionRecords = append(newPermissionRecords, permissionRecord{
			id:            utils2.UUID(),
			spType:        singleUser,
			parameter:     p.jiraUIDNameMap[uid],
			scheme:        scheme,
			permissionKey: permissionKey,
		})
	}
	return newPermissionRecords
}

func (p *JiraResolver) mergeSpecialUserDomain(projectIDs []string, scheme, permissionKey string,
	permissionRecordsToMerge ...[]permissionRecord) (thirdProjectPermissions []*resolve.ThirdProjectPermission) {
	defer func() {
		// 兜底
		if len(thirdProjectPermissions) > 0 {
			return
		}
		permissionRecords := []permissionRecord{
			{
				id:            utils2.UUID(), // 随机生成resourceID，不要和其他的重复
				scheme:        scheme,
				permissionKey: permissionKey,
				spType:        lead,
				parameter:     "",
			},
		}
		thirdProjectPermissions = p.genProjectPermissionSlice(projectIDs, permissionKey, permissionRecords)
	}()
	if len(permissionRecordsToMerge) == 0 {
		return
	}
	var newPermissionRecords []permissionRecord
	var toIntersect = [][]string{}
	var specialUserDomainPermissionRecordMap = make(map[string]permissionRecord)
	for _, permissionRecords := range permissionRecordsToMerge {
		var specialUserDomain []string
		for _, permissionRecord := range permissionRecords {
			if permissionRecord.spType == assignee || permissionRecord.spType == reporter ||
				permissionRecord.spType == lead {
				specialUserDomain = append(specialUserDomain, permissionRecord.spType)
				specialUserDomainPermissionRecordMap[permissionRecord.spType] = permissionRecord
			}
		}
		toIntersect = append(toIntersect, specialUserDomain)
	}
	specialUserDomains := intersectUIDs(toIntersect...)
	for _, key := range specialUserDomains {
		permissionRecord := specialUserDomainPermissionRecordMap[key]
		permissionRecord.scheme = scheme
		permissionRecord.permissionKey = permissionKey
		newPermissionRecords = append(newPermissionRecords, permissionRecord)
	}
	thirdProjectPermissions = p.genProjectPermissionSlice(projectIDs, permissionKey, newPermissionRecords)
	return
}

func (p *JiraResolver) genProjectPermissionSlice(projectIDs []string, permissionKey string,
	permissionRecords []permissionRecord) []*resolve.ThirdProjectPermission {
	var projectPermissionSlice []*resolve.ThirdProjectPermission
	permissions, ok := projectPermissionMap[permissionKey]
	if !ok {
		return projectPermissionSlice
	}
	for _, permissionRecord := range permissionRecords {
		for _, per := range permissions {
			// 这个权限应该在哪个上下文用，项目或者工作项
			permissionCtx, ok := permissionContextMap[per]
			if !ok {
				continue
			}
			// jira的用户组等映射到我们的用户组，在项目上下文是什么，在工作项上下文又是什么
			userDomainWithContext, ok := userDomainMap[permissionRecord.spType]
			if !ok {
				continue
			}
			// 在某种上下文用户域是什么
			userDomain, ok := userDomainWithContext[permissionCtx]
			if !ok {
				continue
			}
			if userDomain == "everyone" && permissionRecord.parameter != "" {
				continue
			}
			for _, pid := range projectIDs {
				_, ok := p.jiraProjectIDNameMap[pid]
				if !ok {
					continue
				}

				curParameter := ""
				domainType := userDomain
				switch domainType {
				case constants.SingleUserLabel:
					curParameter = p.jiraUserNameIDMap[permissionRecord.parameter]
					if curParameter == "" {
						log.Printf("p.jiraUserNameIDMap[permissionRecord.parameter] is empty, permissionRecord: %+v",
							permissionRecord)
						continue
					}
				case constants.GroupLabel:
					curParameter = p.jiraUserGroupNameIDMap[permissionRecord.parameter]
				case constants.RoleLabel:
					curParameter = permissionRecord.parameter
				}

				r := new(resolve.ThirdProjectPermission)
				// 所有项目的都是一样的
				r.ResourceID = permissionRecord.id
				r.ProjectID = pid
				r.UserDomainType = userDomain
				r.ContextType = permissionCtx
				r.Permission = per
				r.UserDomainParam = curParameter
				projectPermissionSlice = append(projectPermissionSlice, r)
			}
		}
	}
	return projectPermissionSlice
}

func intersectUIDs(toIntersect ...[]string) []string {
	if len(toIntersect) == 0 {
		return []string{}
	}
	if len(toIntersect) == 1 {
		return toIntersect[0]
	}
	result := toIntersect[0]
	for i := 1; i < len(toIntersect); i++ {
		if len(toIntersect[i]) == 0 {
			return nil
		}
		result = utils2.StringArrayIntersection(result, toIntersect[i])
	}
	return result
}

func (p *JiraResolver) prepareChangeGroup() error {
	for {
		element, err := p.nextElement("ChangeGroup")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}

		id := getAttributeValue(element, "id")
		taskID := getAttributeValue(element, "issue")
		userName := getAttributeValue(element, "author")
		createTime := getAttributeValue(element, "created")
		userID, found := p.jiraUserNameIDMap[userName]
		if !found {
			log.Println("jira user missing: %s", userName)
			continue
		}
		r := new(changeGroup)
		r.TaskID = taskID
		r.UserID = userID
		r.CreateTime = timestamp.StringToInt64(createTime)
		p.jiraChangeGroupMap[id] = r
	}
	return nil
}

func (p *JiraResolver) prepareProject() {
	for {
		element, err := p.nextElement("Project")
		if err != nil {
			log.Printf("prepare project fail: %s", err)
			break
		}
		if element == nil {
			break
		}
		id := getAttributeValue(element, "id")
		if p.ignoreProjectID(id) {
			continue
		}
		p.jiraProjectElements = append(p.jiraProjectElements, element)
		originalKey := getAttributeValue(element, "originalkey")
		p.jiraProjectIDOriginalKeyMap[id] = originalKey
	}
	return
}

func (p *JiraResolver) nextElement(tag string) (*etree.Element, error) {
	line, err := p.nextLineFromTagFile(tag)
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}
	document := etree.NewDocument()
	if err := document.ReadFromString(line); err != nil {
		return nil, errors.Errorf("parse element fail: %s", line)
	}
	return document.Root(), nil
}

func (p *JiraResolver) nextElementFromReader(reader *resolve.XmlScanner) (*etree.Element, error) {
	line, err := p.nextLineFromReader(reader)
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}

	document := etree.NewDocument()
	if err := document.ReadFromString(line); err != nil {
		return nil, errors.Errorf("parse element fail: %s", line)
	}
	return document.Root(), nil
}

func (p *JiraResolver) nextElements(tag string, number int) ([]*etree.Element, error) {
	reader, ok := p.tagFilesMap[tag]
	if !ok {
		log.Printf("(%s) tag file is not exist", tag)
		return nil, nil
	}

	lines, err := p.nextLinesFromReader(reader, number)
	if err != nil {
		return nil, err
	}

	if lines == "" {
		return nil, nil
	}

	document := etree.NewDocument()
	if err := document.ReadFromString(lines); err != nil {
		return nil, err
	}

	elements := document.FindElements("./" + tag)
	return elements, nil
}

func (p *JiraResolver) nextLineFromReader(reader *resolve.XmlScanner) (string, error) {
	e := reader.NextElement()
	if e == nil {
		return "", nil
	}
	return e.Encode(), nil
}

func (p *JiraResolver) nextLinesFromReader(reader *resolve.XmlScanner, number int) (string, error) {
	var lines strings.Builder
	es := reader.NextElements(number)
	if len(es) == 0 {
		return "", nil
	}
	for _, e := range es {
		lines.WriteString(e.Encode())
	}

	return lines.String(), nil
}

func (p *JiraResolver) nextLineFromTagFile(tag string) (string, error) {
	reader, ok := p.tagFilesMap[tag]
	if !ok {
		log.Printf("(%s) tag file is not exist", tag)
		return "", nil
	}
	line, err := p.nextLineFromReader(reader)
	if err != nil {
		return line, err
	}
	return line, nil
}

func (p *JiraResolver) prepareEventTypeIDNamesMap() error {
	for {
		element, err := p.nextElement("EventType")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}
		idStr := getAttributeValue(element, "id")
		name := getAttributeValue(element, "name")
		typeNames, found := eventTypeMap[name]
		if !found {
			continue
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("id error: %s", idStr)
			continue
		}
		p.jiraEventTypeIDNamesMap[id] = typeNames
	}
	return nil
}

func (p *JiraResolver) prepareSprintColumn() {
	columns := p.sprintRoot.FindElements("./column")
	var name string
	for _, column := range columns {
		name = getAttributeValue(column, "name")
		p.jiraSprintColumnNames = append(p.jiraSprintColumnNames, name)
	}
}

func (p *JiraResolver) resolvePrepareConfigurationContext() error {
	p.jiraFieldConfigSchemeProjectIDsMap = make(map[string][]string)
	for {
		element, err := p.nextElement("ConfigurationContext")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}
		projectID := getAttributeValue(element, "project")
		key := getAttributeValue(element, "key")
		fieldConfigScheme := getAttributeValue(element, "fieldconfigscheme")
		if projectID != "" && key == configContextIssueType {
			p.jiraFieldConfigSchemeProjectIDsMap[fieldConfigScheme] = append(
				p.jiraFieldConfigSchemeProjectIDsMap[fieldConfigScheme],
				projectID,
			)
		}
	}
	return nil
}

func (p *JiraResolver) prepareConfigurationContext() error {
	for {
		element, err := p.nextElement("ConfigurationContext")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}
		projectID := getAttributeValue(element, "project")
		if p.ignoreProjectID(projectID) {
			continue
		}
		key := getAttributeValue(element, "key")
		fieldConfigScheme := getAttributeValue(element, "fieldconfigscheme")
		if projectID != "" && key == configContextIssueType {
			p.jiraFieldConfigSchemeProjectIDsMap[fieldConfigScheme] = append(
				p.jiraFieldConfigSchemeProjectIDsMap[fieldConfigScheme],
				projectID,
			)
		}
	}
	return nil
}

func (p *JiraResolver) afterAction(action string, err error) error {
	if err != nil {
		return err
	}
	switch action {
	case "NextUserGroupMember":
		if err := p.prepareProjectRoleActor(); err != nil {
			return err
		}
	case "NextProjectIssueType":
		if err := p.prepareWorkflow(); err != nil {
			return err
		}
		if err := p.prepareIssueTypeLayout(); err != nil {
			return err
		}
		if err := p.prepareProjectIssueTypeWithFields(); err != nil {
			return err
		}
		if err := p.prepareProjectIssueTypeLayout(); err != nil {
			return err
		}
	case "NextProject":
		if err := p.prepareProjectPermission(); err != nil {
			return err
		}
		if err := p.prepareProjectIssueType(); err != nil {
			return err
		}
	case "NextUserGroup":
		if err := p.prepareGlobalPermission(); err != nil {
			return err
		}
	case "NextUser":
		if err := p.prepareChangeGroup(); err != nil {
			return err
		}
	case "NextTaskLink":
		if err := p.prepareNotification(); err != nil {
			return err
		}
	}
	return nil
}

func (p *JiraResolver) prepareGlobalPermission() error {
	for {
		element, err := p.nextElement("GlobalPermissionEntry")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}

		permissionValue := getAttributeValue(element, "permission")
		if list, found := globalPermissionMap[permissionValue]; found {
			for _, v := range list {
				r := new(resolve.ThirdGlobalPermission)
				r.ResourceID = getAttributeValue(element, "id")
				r.Permission = v
				r.GroupID = p.jiraUserGroupNameIDMap[getAttributeValue(element, "group_id")]
				p.jiraGlobalPermissionSlice = append(p.jiraGlobalPermissionSlice, r)
			}
		}
	}
	return nil
}

func (p *JiraResolver) prepareProjectRoleActor() error {
	p.jiraProjectRolesSlice = make([]*projectRoleStruct, 0)
	p.jiraProjectMembersSlice = make([]*projectMemberStruct, 0)
	for {
		element, err := p.nextElement("ProjectRoleActor")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}
		pid := getAttributeValue(element, "pid")
		if p.ignoreProjectID(pid) {
			continue
		}
		// 要记录项目角色 -> [用户组，单个用户]的映射
		roleId := getAttributeValue(element, "projectroleid")
		roleType := getAttributeValue(element, "roletype")
		roleTypeParameter := getAttributeValue(element, "roletypeparameter")
		if roleId == "" {
			continue
		}
		projectRoleKey := pid + roleId
		p.jiraProjectRoleIDDataMap[projectRoleKey] = append(p.jiraProjectRoleIDDataMap[projectRoleKey],
			roleData{
				roleType:          roleType,
				roleTypeParameter: roleTypeParameter,
			})
		if roleType == roleTypeGroup {
			if _, found := p.jiraUserGroupNameIDMap[roleTypeParameter]; !found {
				log.Printf("non-exists user group1: %s", roleTypeParameter)
				continue
			}
			userGroupID := p.jiraUserGroupNameIDMap[roleTypeParameter]
			if _, found := p.jiraUserGroupIDMembersMap[userGroupID]; !found {
				log.Printf("non-exists user group2: %s", userGroupID)
				continue
			}
			for _, userID := range p.jiraUserGroupIDMembersMap[userGroupID] {
				p.jiraProjectMembersSlice = append(p.jiraProjectMembersSlice, &projectMemberStruct{
					ProjectID: pid,
					RoleID:    roleId,
					UserID:    userID,
				})
			}
		}
		if roleType == roleTypeUser {
			if _, found := p.jiraUserNameIDMap[roleTypeParameter]; !found {
				log.Printf("non-exists user name: %s", roleTypeParameter)
				continue
			}
			p.jiraProjectMembersSlice = append(p.jiraProjectMembersSlice, &projectMemberStruct{
				ProjectID: pid,
				RoleID:    roleId,
				UserID:    p.jiraUserNameIDMap[roleTypeParameter],
			})
		}
		p.jiraProjectRolesSlice = append(p.jiraProjectRolesSlice, &projectRoleStruct{RoleID: roleId, ProjectID: pid})
	}
	return nil
}

func (p *JiraResolver) ignoreProjectID(pid string) bool {
	if _, found := p.importTask.SelectedProjectIDs[pid]; found {
		return false
	}
	return true
}

func (p *JiraResolver) ignoreTaskID(pid string) bool {
	if _, found := p.jiraTaskIDMap[pid]; found {
		return false
	}
	return true
}

func (p *JiraResolver) prepareWorkflow() error {
	flowMap, err := p.mapWorkFlow()
	if err != nil {
		return err
	}
	projectIDResultMap := make(map[string][]*resolve.ThirdWorkflow)
	projectIDTransitionsMap := make(map[string][]*resolve.ThirdTransition)

	for {
		element, err := p.nextElement("Workflow")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}
		flowName := getAttributeValue(element, "name")
		descNode := element.FindElement("./descriptor")
		if descNode == nil {
			log.Printf("[importer] parse worker fail: %s", descNode)
			continue
		}
		docDesc := etree.NewDocument()
		if err := docDesc.ReadFromString(descNode.Text()); err != nil {
			log.Printf("parse xml fail %s:%s", getAttributeValue(element, "name"), descNode.Text())
			continue
		}

		steps := docDesc.FindElements("./workflow/steps/step")
		initialAction := docDesc.FindElement("./workflow/initial-actions/action")
		commonActions := docDesc.FindElements("./workflow/common-actions/action")
		globalActions := docDesc.FindElements("./workflow/global-actions/action")

		var statuses []string
		stepIDStatusMap := make(map[string]string)
		for _, step := range steps {
			id := getAttributeValue(step, "id")
			for _, meta := range step.FindElements("./meta") {
				name := getAttributeValue(meta, "name")
				if name == jiraStatusID {
					stepIDStatusMap[id] = meta.Text()
					statuses = append(statuses, meta.Text())
					break
				}
			}
		}
		var MapCommonAction = make(map[string]*etree.Element)
		for _, a := range commonActions {
			id := getAttributeValue(a, "id")
			MapCommonAction[id] = a
		}

		transitionsWithoutProject := make([]*resolve.ThirdTransition, 0)

		// 工作流 Transition condition
		for _, step := range steps {
			actions := step.FindElements("./actions/action")
			id := getAttributeValue(step, "id")
			for _, action := range actions {
				name := getAttributeValue(action, "name")
				unConditionalResultStep := getAttributeValue(action.FindElement("./results/unconditional-result"), "step")
				conditions := action.FindElements("./restrict-to/conditions/condition")
				userDomains := buildUserDomainFromCondition(conditions)
				transitionsWithoutProject = append(transitionsWithoutProject, &resolve.ThirdTransition{
					StartStatusID: stepIDStatusMap[id],
					Name:          name,
					EndStatusID:   stepIDStatusMap[unConditionalResultStep],
					UserDomains:   userDomains,
				})
			}

			curStepCommonActions := step.FindElements("./actions/common-action")
			for _, a := range curStepCommonActions {
				commonStepId := getAttributeValue(a, "id")
				action, ok := MapCommonAction[commonStepId]
				if !ok {
					continue
				}
				conditions := action.FindElements("./restrict-to/conditions/condition")
				unConditionalResultStep := getAttributeValue(action.FindElement("./results/unconditional-result"), "step")
				name := getAttributeValue(action, "name")
				userDomains := buildUserDomainFromCondition(conditions)
				if _, found := stepIDStatusMap[id]; !found {
					continue
				}
				transitionsWithoutProject = append(transitionsWithoutProject, &resolve.ThirdTransition{
					StartStatusID: stepIDStatusMap[id],
					Name:          name,
					EndStatusID:   stepIDStatusMap[unConditionalResultStep],
					UserDomains:   userDomains,
				})
			}
		}
		for _, action := range globalActions {
			conditions := action.FindElements("./restrict-to/conditions/condition")
			unConditionalResultStep := getAttributeValue(action.FindElement("./results/unconditional-result"), "step")
			name := getAttributeValue(action, "name")
			userDomains := buildUserDomainFromCondition(conditions)
			transitionsWithoutProject = append(transitionsWithoutProject, &resolve.ThirdTransition{
				StartStatusID: taskStatusAllUUID,
				Name:          name,
				EndStatusID:   stepIDStatusMap[unConditionalResultStep],
				UserDomains:   userDomains,
			})
		}

		fMaps := flowMap[flowName]
		for _, f := range fMaps {
			for _, projectID := range f.ProjectIDs {
				ret := new(resolve.ThirdWorkflow)
				ret.ProjectID = projectID
				itStatus := &resolve.ThirdTaskStatusConfig{
					IssueTypeID:  f.IssueType,
					StatusIDs:    statuses,
					WorkflowName: flowName,
				}
				if initialAction != nil {
					stepIndex := getAttributeValue(initialAction.FindElement("./results/unconditional-result"), "step")
					itStatus.InitStatusID = stepIDStatusMap[stepIndex]
				}

				ret.IssueTypeStatuses = append(ret.IssueTypeStatuses, itStatus)
				for _, t := range transitionsWithoutProject {
					projectIDTransitionsMap[projectID] = append(projectIDTransitionsMap[projectID], &resolve.ThirdTransition{
						ProjectID:     projectID,
						IssueTypeID:   f.IssueType,
						StartStatusID: t.StartStatusID,
						Name:          t.Name,
						EndStatusID:   t.EndStatusID,
						UserDomains:   t.UserDomains,
					})
				}
				projectIDResultMap[projectID] = append(projectIDResultMap[projectID], ret)
			}
		}
	}

	for projectID, results := range projectIDResultMap {
		result := results[0]
		if len(results) > 0 {
			for i := 1; i < len(results); i++ {
				result.IssueTypeStatuses = append(result.IssueTypeStatuses, results[i].IssueTypeStatuses...)
			}
		}
		if p.ignoreProjectID(projectID) {
			continue
		}
		result.Transitions = projectIDTransitionsMap[projectID]
		p.jiraWorkflowsSlice = append(p.jiraWorkflowsSlice, result)
	}
	return nil
}

func (p *JiraResolver) prepareProjectType() {
	curProjectField := new(resolve.ThirdGlobalProjectField)
	curProjectField.ResourceID = ResourceIDProjectType
	curProjectField.Name = curProjectField.ResourceID
	curProjectField.Type = fieldModel.FieldTypeOption
	curProjectField.Options = make([]*resolve.ThirdGlobalProjectFieldOption, 0)

	projectTypeOptionMap := make(map[string]bool)
	projectTypeIDMap := make(map[string]string)
	for {
		if p.projectIndex >= len(p.jiraProjectElements) {
			p.projectIndex = 0
			break
		}
		element := p.jiraProjectElements[p.projectIndex]
		p.projectIndex++
		projectType := getAttributeValue(element, "projecttype")
		projectID := getAttributeValue(element, "id")
		projectURL := getAttributeValue(element, "url")
		projectDescription := getAttributeValue(element, "description")
		projectFieldValues := []*resolve.ThirdProjectFieldValue{
			{
				Base: resolve.Base{
					ResourceID: projectID,
				},
				ProjectID:      projectID,
				ProjectFieldID: ResourceIDProjectDescription,
				Value:          projectDescription,
			},
			{
				Base: resolve.Base{
					ResourceID: projectID,
				},
				ProjectID:      projectID,
				ProjectFieldID: ResourceIDProjectURL,
				Value:          projectURL,
			},
		}
		p.jiraProjectFieldValue = append(p.jiraProjectFieldValue, projectFieldValues...)
		if projectType == "" {
			continue
		}
		if !projectTypeOptionMap[projectType] {
			curOption := new(resolve.ThirdGlobalProjectFieldOption)
			curOption.ResourceID = projectType
			curOption.Value = projectType
			curOption.Color = "#fff"
			curProjectField.Options = append(curProjectField.Options, curOption)
			projectTypeOptionMap[projectType] = true
		}

		curProjectFieldValue := new(resolve.ThirdProjectFieldValue)
		curProjectFieldValue.ResourceID = projectID // project type 没有id， 与项目一一对应，所以使用项目id作为属性id
		curProjectFieldValue.ProjectID = projectID
		curProjectFieldValue.ProjectFieldID = ResourceIDProjectType
		curProjectFieldValue.Type = fieldModel.FieldTypeOption
		curProjectFieldValue.Value = projectType

		p.jiraProjectFieldValue = append(p.jiraProjectFieldValue, curProjectFieldValue)
	}
	p.jiraGlobalProjectField = append(p.jiraGlobalProjectField, curProjectField)
	p.jiraProjectTypeIDMap = projectTypeIDMap
}

func (p *JiraResolver) prepareProjectCategory() {
	curProjectField := new(resolve.ThirdGlobalProjectField)
	curProjectField.ResourceID = ResourceIDProjectCategory
	curProjectField.Name = curProjectField.ResourceID
	curProjectField.Type = fieldModel.FieldTypeOption
	curProjectField.Options = make([]*resolve.ThirdGlobalProjectFieldOption, 0)

	for {
		element, err := p.nextElement("ProjectCategory")
		if element == nil {
			break
		}
		if err != nil {
			log.Printf("get project element fail")
			break
		}
		name := getAttributeValue(element, "name")
		desc := getAttributeValue(element, "description")
		optionID := getAttributeValue(element, "id")
		if name == "" || optionID == "" {
			continue
		}
		curOption := new(resolve.ThirdGlobalProjectFieldOption)
		curOption.ResourceID = optionID
		curOption.Value = name
		curOption.Color = "#fff"
		curOption.Desc = desc
		curProjectField.Options = append(curProjectField.Options, curOption)
	}
	p.jiraGlobalProjectField = append(p.jiraGlobalProjectField, curProjectField)

	for projectID, categoryID := range p.jiraCategoryProjectAssociation {
		curProjectFieldValue := new(resolve.ThirdProjectFieldValue)
		curProjectFieldValue.ResourceID = fmt.Sprintf("%s-%s", projectID, categoryID)
		curProjectFieldValue.ProjectID = projectID
		curProjectFieldValue.Type = fieldModel.FieldTypeOption
		curProjectFieldValue.Value = categoryID
		curProjectFieldValue.ProjectFieldID = ResourceIDProjectCategory

		p.jiraProjectFieldValue = append(p.jiraProjectFieldValue, curProjectFieldValue)
	}
}

func (p *JiraResolver) mapWorkFlow() (map[string][]*flowMap, error) {
	flowMaps := make(map[string][]*flowMap)
	entityMap := make(map[string][]*etree.Element, 0)

	for {
		element, err := p.nextElement("WorkflowSchemeEntity")
		if err != nil {
			return nil, err
		}
		if element == nil {
			break
		}
		scheme := getAttributeValue(element, "scheme")
		if _, found := entityMap[scheme]; !found {
			entityMap[scheme] = make([]*etree.Element, 0)
		}
		entityMap[scheme] = append(entityMap[scheme], element)
	}

	for scheme, element := range entityMap {
		projectIDs := p.jiraWorkflowSchemeProjectAssoc[scheme]
		projectIDs = utils2.UniqueNoNullSlice(projectIDs...)
		for _, entity := range element {
			fm := new(flowMap)
			fm.ProjectIDs = projectIDs
			fm.IssueType = getAttributeValue(entity, "issuetype")
			workflowName := getAttributeValue(entity, "workflow")
			flowMaps[workflowName] = append(flowMaps[workflowName], fm)

		}
	}
	return flowMaps, nil
}

func buildUserDomainFromCondition(conditions []*etree.Element) []*resolve.ThirdUserDomain {
	userDomains := make([]*resolve.ThirdUserDomain, 0)
	for _, condition := range conditions {
		var isAssign, isReporter bool
		var projectRoleID, groupKey string
		args := condition.FindElements("./arg")
		for _, arg := range args {
			name := getAttributeValue(arg, "name")
			if name == "jira.projectrole.id" {
				projectRoleID = name
			} else if name == "group" {
				groupKey = name
			} else if name == "class.name" && strings.Contains(name, "AllowOnlyReporter") {
				isReporter = true
			} else if name == "class.name" && strings.Contains(name, "AllowOnlyAssignee") {
				isAssign = true
			}
		}
		if len(projectRoleID) > 0 {
			userDomains = append(userDomains, &resolve.ThirdUserDomain{
				UserDomainType:  "role",
				UserDomainParam: projectRoleID,
			})
		} else if len(groupKey) > 0 {
			userDomains = append(userDomains, &resolve.ThirdUserDomain{
				UserDomainType:  "group",
				UserDomainParam: groupKey,
			})
		} else if isAssign {
			userDomains = append(userDomains, &resolve.ThirdUserDomain{
				UserDomainType: "task_assign",
			})
		} else if isReporter {
			userDomains = append(userDomains, &resolve.ThirdUserDomain{
				UserDomainType: "task_owner",
			})
		}
	}
	return userDomains
}

func (p *JiraResolver) prepareApplicationUser() error {
	for {
		element, err := p.nextElement("ApplicationUser")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}
		id := getAttributeValue(element, "id")
		userName := getAttributeValue(element, "userKey")
		if id == "" || userName == "" {
			log.Printf("prepare application user map fail: %s, %s", id, userName)
			continue
		}
		p.jiraUserNameIDMap[userName] = id
		p.jiraUIDNameMap[id] = userName
	}
	return nil
}

func (p *JiraResolver) resolvePrepareProjectIssueType() (map[string][]string, error) {
	optionConfigProjectIssueTypeMap := make(map[string]struct{})
	p.jiraProjectIssueTypeMap = make(map[string]struct{})
	for {
		element, err := p.nextElement("OptionConfiguration")
		if element == nil {
			break
		}
		if err != nil {
			log.Printf("prepare project issue type: %s", err)
			return nil, err
		}
		fieldID := getAttributeValue(element, "fieldid")
		scheme := getAttributeValue(element, "fieldconfig")
		projectIDs, ok := p.jiraFieldConfigSchemeProjectIDsMap[scheme]
		if !ok || len(projectIDs) == 0 || fieldID != configContextIssueType {
			continue
		}

		issueTypeID := getAttributeValue(element, "optionid")
		for _, projectID := range projectIDs {
			mapKey := fmt.Sprintf("%s:%s", projectID, issueTypeID)
			optionConfigProjectIssueTypeMap[mapKey] = struct{}{}
		}
	}

	for thirdProjectIssueType, _ := range optionConfigProjectIssueTypeMap {
		p.jiraProjectIssueTypeMap[thirdProjectIssueType] = struct{}{}
	}

	projectIssueTypeMap := make(map[string][]string)
	for key, _ := range p.jiraProjectIssueTypeMap {
		keys := strings.Split(key, ":")
		if len(keys) != 2 {
			log.Printf("project issue type error: %s", key)
			continue
		}
		projectID := keys[0]
		issueTypeID := keys[1]
		projectIssueTypeMap[projectID] = append(projectIssueTypeMap[projectID], issueTypeID)
	}

	p.jiraProjectIssueTypeMap = nil
	return projectIssueTypeMap, nil
}

func (p *JiraResolver) prepareProjectIssueType() error {
	optionConfigProjectIssueTypeMap := make(map[string]struct{})
	for {
		element, err := p.nextElement("OptionConfiguration")
		if element == nil {
			break
		}
		if err != nil {
			log.Printf("prepare project issue type: %s", err)
			return err
		}
		fieldID := getAttributeValue(element, "fieldid")
		scheme := getAttributeValue(element, "fieldconfig")
		projectIDs, ok := p.jiraFieldConfigSchemeProjectIDsMap[scheme]
		if !ok || len(projectIDs) == 0 || fieldID != configContextIssueType {
			continue
		}

		issueTypeID := getAttributeValue(element, "optionid")
		for _, projectID := range projectIDs {
			if _, found := p.jiraProjectIDNameMap[projectID]; !found {
				continue
			}
			mapKey := fmt.Sprintf("%s:%s", projectID, issueTypeID)
			optionConfigProjectIssueTypeMap[mapKey] = struct{}{}
		}
	}

	for thirdProjectIssueType, _ := range optionConfigProjectIssueTypeMap {
		p.jiraProjectIssueTypeMap[thirdProjectIssueType] = struct{}{}
	}

	for key, _ := range p.jiraProjectIssueTypeMap {
		keys := strings.Split(key, ":")
		if len(keys) != 2 {
			log.Printf("project issue type error: %s", key)
			continue
		}
		r := new(resolve.ThirdProjectIssueType)
		r.ResourceID = keys[0]
		r.ProjectID = keys[0]
		r.IssueTypeID = keys[1]
		p.jiraProjectIssueTypeSlice = append(p.jiraProjectIssueTypeSlice, r)
		p.jiraProjectIssueTypeNameMap[key] = true
		p.jiraProjectWithIssueTypeMap[r.ProjectID] = append(p.jiraProjectWithIssueTypeMap[r.ProjectID], r.IssueTypeID)
	}

	p.jiraProjectIssueTypeMap = nil
	return nil
}

func (p *JiraResolver) ServerID() string {
	return p.serverID
}

func getAttributeValue(element *etree.Element, attribute string) string {
	var e = element
	var name = attribute
	var resp string
	a := e.SelectAttr(name)
	if a == nil {
		child := e.SelectElement(name)
		if child != nil {
			resp = child.Text()
		}
	} else {
		resp = a.Value
	}
	return resp
}

func (p *JiraResolver) NextTaskRelease() ([]byte, error) {
	if p.taskReleaseIndex >= len(p.jiraTaskReleaseSlice) {
		p.jiraTaskReleaseSlice = nil
		return nil, nil
	}
	r := p.jiraTaskReleaseSlice[p.taskReleaseIndex]
	p.taskReleaseIndex++
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextGlobalProjectField() ([]byte, error) {
	if p.globalProjectFieldIndex >= len(p.jiraGlobalProjectField) {
		p.jiraGlobalProjectField = nil
		return nil, nil
	}
	r := p.jiraGlobalProjectField[p.globalProjectFieldIndex]
	p.globalProjectFieldIndex++
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextProjectFieldValue() ([]byte, error) {
	if p.projectFieldValueIndex >= len(p.jiraProjectFieldValue) {
		p.jiraProjectFieldValue = nil
		return nil, nil
	}
	r := p.jiraProjectFieldValue[p.projectFieldValueIndex]
	p.projectFieldValueIndex++
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextIssueTypeField() ([]byte, error) {
	if p.issueTypeFieldIndex < len(thirdIssueTypeFields) {
		r := thirdIssueTypeFields[p.issueTypeFieldIndex]
		p.issueTypeFieldIndex++
		return utils2.OutputJSON(r), nil
	}

	element, err := p.nextElement("FieldConfigSchemeIssueType")
	if element == nil || err != nil {
		return nil, p.afterAction("NextProject", err)
	}

	r := new(resolve.ThirdIssueTypeField)
	r.ResourceID = getAttributeValue(element, "id")
	r.IssueTypeID = getAttributeValue(element, "issuetype")
	scheme := getAttributeValue(element, "fieldconfigscheme")

	if r.IssueTypeID == "" {
		r.IssueTypeID = "all"
	}

	fieldID, ok := p.jiraFieldConfigSchemeMap[scheme]
	if !ok {
		return p.NextIssueTypeField()
	}

	r.FieldID = strings.Split(fieldID, "_")[1]

	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextProject() ([]byte, error) {
	if p.projectIndex >= len(p.jiraProjectElements) {
		p.projectIndex = 0
		return nil, p.afterAction("NextProject", nil)
	}
	element := p.jiraProjectElements[p.projectIndex]
	p.projectIndex++

	r := new(resolve.ThirdProject)
	r.ResourceID = getAttributeValue(element, "id")
	r.Name = getAttributeValue(element, "name")
	r.Status = utilsModel.ProjectStatusNormal
	r.CreateTime = time.Now().Unix() * 1000
	r.StatusUUID = project.ProjectStatusCategoryTodo
	r.Announcement = ""
	r.Type = project.ProjectTypeAgileEnum
	r.AssignID = p.jiraUserNameIDMap[getAttributeValue(element, "lead")]
	p.jiraProjectIDNameMap[r.ResourceID] = r.Name
	p.jiraProjectIDAssignIDMap[r.ResourceID] = r.AssignID

	perScheme := element.SelectAttr("permissionscheme")
	if perScheme != nil {
		permissionScheme := perScheme.Value
		projectIDs := p.jiraPermissionProjectAssociation[permissionScheme]
		if !utils2.StringArrayContains(projectIDs, r.ResourceID) {
			p.jiraPermissionProjectAssociation[permissionScheme] = append(p.jiraPermissionProjectAssociation[permissionScheme], r.ResourceID)
		}
	}
	notification := element.SelectAttr("permissionscheme")
	if notification != nil {
		notificationScheme := notification.Value
		projectIDs := p.jiraNotificationSchemeProjectAssoc[notificationScheme]
		if !utils2.StringArrayContains(projectIDs, r.ResourceID) {
			p.jiraNotificationSchemeProjectAssoc[notificationScheme] = append(p.jiraNotificationSchemeProjectAssoc[notificationScheme], r.ResourceID)
		}
	}
	workflow := element.SelectAttr("workflowscheme")
	if workflow != nil {
		workflowScheme := workflow.Value
		p.jiraWorkflowSchemeProjectAssoc[workflowScheme] = append(p.jiraWorkflowSchemeProjectAssoc[workflowScheme], r.ResourceID)
	}

	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextChangeItem() ([]byte, bool, error) {
	element, err := p.nextElement("ChangeItem")
	if element == nil || err != nil {
		return nil, false, err
	}

	groupID := getAttributeValue(element, "group")

	group, found := p.jiraChangeGroupMap[groupID]
	if !found {
		return nil, true, nil
	}
	if p.ignoreTaskID(group.TaskID) {
		return nil, true, nil
	}
	r := new(resolve.ThirdChangeItem)
	r.ResourceID = getAttributeValue(element, "id")
	r.TaskID = group.TaskID
	r.UserID = group.UserID
	r.CreateTime = group.CreateTime
	r.Field = getAttributeValue(element, "field")
	r.FieldType = getAttributeValue(element, "fieldtype")
	r.OldString = getAttributeValue(element, "oldstring")
	r.OldValue = getAttributeValue(element, "oldvalue")
	r.NewString = getAttributeValue(element, "newstring")
	r.NewValue = getAttributeValue(element, "newvalue")

	if r.Field == constants.ChangeItemFieldLink {
		if r.NewValue != "" {
			r.NewValue = p.jiraTaskKeyIDMap[r.NewValue]
		}
		if r.OldValue != "" {
			r.OldValue = p.jiraTaskKeyIDMap[r.OldValue]
		}
	}
	if r.Field == constants.ChangeItemFieldEpicChild {
		r.Field = constants.ChangeItemFieldLink
	}

	if r.Field == constants.ChangeItemFieldAssignee {
		r.NewValue = p.jiraUserNameIDMap[r.NewValue]
		r.OldValue = p.jiraUserNameIDMap[r.OldValue]
	}

	return utils2.OutputJSON(r), false, nil
}

func (p *JiraResolver) NextIssueTypeLayout() ([]byte, error) {
	if p.issueTypeLayoutIndex >= len(p.jiraIssueTypeLayoutSlice) {
		p.jiraIssueTypeLayoutSlice = nil
		return nil, nil
	}
	r := p.jiraIssueTypeLayoutSlice[p.issueTypeLayoutIndex]
	p.issueTypeLayoutIndex++
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextProjectIssueTypeLayout() ([]byte, error) {
	if p.projectIssueTypeLayoutIndex >= len(p.jiraProjectIssueTypeLayoutSlice) {
		p.jiraProjectIssueTypeLayoutSlice = nil
		return nil, nil
	}
	r := p.jiraProjectIssueTypeLayoutSlice[p.projectIssueTypeLayoutIndex]
	p.projectIssueTypeLayoutIndex++
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextProjectIssueTypeField() ([]byte, error) {
	if p.projectIssueTypeFieldIndex >= len(p.jiraProjectIssueTypeWithFieldsSlice) {
		p.jiraProjectIssueTypeWithFieldsSlice = nil
		return nil, nil
	}
	r := p.jiraProjectIssueTypeWithFieldsSlice[p.projectIssueTypeFieldIndex]
	p.projectIssueTypeFieldIndex++
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextProjectPermission() ([]byte, error) {
	if p.projectPermissionIndex >= len(p.jiraProjectPermissionSlice) {
		p.jiraProjectPermissionSlice = nil
		return nil, nil
	}
	r := p.jiraProjectPermissionSlice[p.projectPermissionIndex]
	p.projectPermissionIndex++
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextGlobalPermission() ([]byte, error) {
	if p.globalPermissionIndex >= len(p.jiraGlobalPermissionSlice) {
		p.jiraGlobalPermissionSlice = nil
		return nil, nil
	}
	r := p.jiraGlobalPermissionSlice[p.globalPermissionIndex]
	p.globalPermissionIndex++
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextSprint() ([]byte, error) {
	if p.sprintIndex >= len(p.jiraSprints) {
		p.jiraSprints = nil
		return nil, nil
	}
	element := p.jiraSprints[p.sprintIndex]
	children := element.ChildElements()

	r := new(resolve.ThirdSprint)
	r.Status = utilsModel.SprintStatusCategoryToDo
	closed := false
	started := false
	for k, child := range children {
		columnName := p.jiraSprintColumnNames[k]
		childText := child.Text()
		switch columnName {
		case "ID":
			r.ResourceID = childText
		case "NAME":
			if childText == "" {
				childText = "without name"
			}
			r.Name = childText
		case "START_DATE":
			if childText == "" {
				childText = strconv.Itoa(int(time.Now().Unix()) * 1000)
			}
			startTime, err := strconv.Atoi(childText)
			if err != nil {
				return nil, err
			}
			r.StartTime = int64(startTime / 1000)
		case "END_DATE":
			if childText == "" {
				childText = strconv.Itoa(int(time.Now().Unix()) * 1000)
			}
			endTime, err := strconv.Atoi(childText)
			if err != nil {
				return nil, err
			}
			r.EndTime = int64(endTime / 1000)
		case "CLOSED":
			if childText == "true" {
				closed = true
			}
			break
		case "STARTED":
			if childText == "true" {
				started = true
			}
		}
	}
	if closed {
		r.Status = utilsModel.SprintStatusCategoryDone
	} else if started {
		r.Status = utilsModel.SprintStatusCategoryInProgress
	}
	r.ProjectID = p.jiraSprintIDProjectIDMap[r.ResourceID]
	r.AssignID = p.jiraProjectIDAssignIDMap[r.ProjectID]
	p.sprintIndex += 1

	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) prepareNotification() error {
	notificationMap := map[string]bool{}
	for {
		element, err := p.nextElement("Notification")
		if element == nil || err != nil {
			return err
		}

		scheme := getAttributeValue(element, "scheme")
		nType := getAttributeValue(element, "type")
		param := getAttributeValue(element, "parameter")
		eventTypeIDStr := getAttributeValue(element, "eventTypeId")
		var eventTypeID int
		if eventTypeIDStr != "" {
			eventTypeID, err = strconv.Atoi(eventTypeIDStr)
			if err != nil {
				log.Printf("get event type id fail: %s.", eventTypeIDStr)
				continue
			}
		}

		subscriber, ok := subscriberMap[nType]
		if !ok {
			log.Printf("not found notification type")
			continue
		}

		valueType := subscriber["type"]
		value := subscriber["value"]
		if nType == "Single_User" || nType == "Group_Dropdown" || nType == "Project_Role" {
			if param == "" {
				continue
			}
			value = param
			if nType == "Single_User" {
				v := p.jiraUserNameIDMap[param]
				if v != "" {
					value = v
				}
			}
			if nType == "Group_Dropdown" {
				v := p.jiraUserGroupNameIDMap[param]
				if v != "" {
					value = v
				}
			}
		}
		eventTypes := p.jiraEventTypeIDNamesMap[eventTypeID]
		projectIDs := p.jiraNotificationSchemeProjectAssoc[scheme]
		projectIDs = unique.Strings(projectIDs)
		for _, eventType := range eventTypes {
			for _, projectID := range projectIDs {
				if p.ignoreProjectID(projectID) {
					continue
				}
				notify := new(resolve.ThirdNotification)
				notify.ValueType = notificationTypeMap[valueType]
				notify.Value = value
				notify.ConfigType = eventType
				notify.ProjectID = projectID
				key := getNotificationKey(notify)
				if notificationMap[key] {
					continue
				}
				notificationMap[key] = true
				p.jiraNotifications = append(p.jiraNotifications, notify)
			}
		}
	}
}

func getNotificationKey(notify *resolve.ThirdNotification) string {
	return fmt.Sprintf("%s:%d:%s:%s", notify.ProjectID, notify.ValueType, notify.Value, notify.ConfigType)
}

func (p *JiraResolver) NextNotification() ([]byte, error) {
	if p.notificationIndex >= len(p.jiraNotifications) {
		p.jiraNotifications = nil
		return nil, nil
	}
	r := p.jiraNotifications[p.notificationIndex]
	p.notificationIndex++

	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextTask() ([]byte, error) {
	if p.versionIndex < len(p.jiraReleaseVersionSlice) {
		element := p.jiraReleaseVersionSlice[p.versionIndex]
		p.versionIndex++
		if element != nil {
			r := new(resolve.ThirdTask)
			r.ResourceID = versionIDToTaskID(getAttributeValue(element, "id"))
			r.Summary = getAttributeValue(element, "name")
			r.Desc = getAttributeValue(element, "description")
			r.ProjectID = getAttributeValue(element, "project")
			r.AssignID = p.jiraProjectIDAssignIDMap[r.ProjectID]
			r.OwnerID = r.AssignID
			r.StatusDetailType = statusModel.DetailTypeUnPublished
			r.IssueTypeDetailType = issuetype.DetailTypePublish
			r.CreatedTime = time.Now().Unix()
			r.UpdatedTime = time.Now().Unix()

			released := getAttributeValue(element, "released")
			archived := getAttributeValue(element, "archived")
			deadline := getAttributeValue(element, "releasedate")

			if deadline != "" {
				dUnix := timestamp.StringToInt64(deadline)
				r.Deadline = &dUnix
			}
			if released == "true" {
				r.StatusDetailType = statusModel.DetailTypePublished
			}
			if archived == "true" {
				r.StatusDetailType = statusModel.DetailTypeClosed
			}
			p.jiraTaskIDMap[r.ResourceID] = struct{}{}
			return utils2.OutputJSON(r), nil
		}
	}

	element, err := p.nextElementFromReader(p.IssueParentReader)
	if err != nil {
		return nil, err
	}
	if element == nil {
		element, err = p.nextElementFromReader(p.IssueChildReader)
		if err != nil || element == nil {
			return nil, err
		}
	}

	var o = element

	r := new(resolve.ThirdTask)

	r.ResourceID = getAttributeValue(o, "id")
	key := getAttributeValue(o, "key")
	if key != "" {
		p.jiraTaskKeyIDMap[key] = r.ResourceID
	}

	parentID, ok := p.jiraIssueParentIDMap[r.ResourceID]
	if ok {
		r.ParentID = parentID
	}

	// convert key to field value
	p.handleIssueKey(o)

	// convert resolution to field value
	p.handleIssueResolution(o)
	r.Summary = getAttributeValue(o, "summary")
	r.Desc = getAttributeValue(o, "description")
	creator := getAttributeValue(o, "creator")

	ownerID, ok := p.jiraUserNameIDMap[creator]
	if !ok {
		log.Printf("creator %s not found.", creator)
	}
	r.OwnerID = ownerID

	assignee := getAttributeValue(o, "assignee")
	assignID, ok := p.jiraUserNameIDMap[assignee]
	if !ok && assignee != "" {
		log.Printf("assignee (%s) not found.", assignee)
	}
	if assignID == "" {
		assignID = r.OwnerID
	}
	r.AssignID = assignID
	duedate := getAttributeValue(o, "duedate")
	if len(duedate) > 0 {
		dueTime := timestamp.StringToInt64(duedate)
		r.Deadline = &dueTime
	}

	r.PriorityID = getAttributeValue(o, "priority")
	r.StatusID = getAttributeValue(o, "status")
	r.ProjectID = getAttributeValue(o, "project")
	r.IssueTypeID = getAttributeValue(o, "type")

	sprintID, _ := p.jiraIssueSprintMap[r.ResourceID]
	if len(parentID) > 0 {
		sprintID, _ = p.jiraIssueSprintMap[parentID]
	}
	r.SprintID = sprintID
	created := getAttributeValue(o, "created")
	r.CreatedTime = timestamp.StringToInt64(created)
	updated := getAttributeValue(o, "updated")
	r.UpdatedTime = timestamp.StringToInt64(updated)

	// timeoriginalestimate 属性，预估工时
	estimateStr := getAttributeValue(o, "timeoriginalestimate")
	if len(estimateStr) > 0 {
		estimateTime, err := strconv.ParseInt(estimateStr, 10, 64)
		if err != nil {
			log.Println("parse timeoriginalestimate %s failed, %+v\n", estimateStr, err)
		}

		now := time.Now().Unix()
		wl := &resolve.ThirdTaskWorkLog{
			Base:       resolve.Base{ResourceID: fmt.Sprintf("%s-%s", r.ResourceID, "timeoriginalestimate")},
			TaskID:     r.ResourceID,
			UserID:     r.OwnerID,
			StartTime:  now,
			Hours:      float64(estimateTime) / 3600,
			Type:       constants.ThirdTaskWorkLogTypeEstimate,
			CreateTime: now,
		}
		p.jiraWorkLogs = append(p.jiraWorkLogs, wl)
	}

	// timeestimate 属性，剩余工时
	remainingTimeStr := getAttributeValue(o, "timeestimate")
	if len(remainingTimeStr) > 0 {
		remainingTime, err := strconv.ParseInt(remainingTimeStr, 10, 64)
		if err != nil {
			log.Println("parse timeestimate %s failed, %+v\n", remainingTimeStr, err)
		}
		now := time.Now().Unix()
		wl := &resolve.ThirdTaskWorkLog{
			Base:       resolve.Base{ResourceID: fmt.Sprintf("%s-%s", r.ResourceID, "timeestimate")},
			TaskID:     r.ResourceID,
			UserID:     r.OwnerID,
			StartTime:  now,
			Hours:      float64(remainingTime) / 3600,
			Type:       constants.ThirdTaskWorkLogTypeRemaining,
			CreateTime: now,
		}
		p.jiraWorkLogs = append(p.jiraWorkLogs, wl)
	}
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextUser() ([]byte, error) {
	element, err := p.nextElement("User")
	if element == nil || err != nil {
		return []byte{}, p.afterAction("NextUser", err)
	}

	email := getAttributeValue(element, "emailAddress")
	if email == "" {
		log.Printf("jira user email empty")
		return p.NextUser()
	}
	if _, found := p.jiraUserEmailMap[email]; found {
		log.Printf("jira import user email exist: %s", email)
		return p.NextUser()
	}
	r := new(resolve.ThirdUser)
	r.ResourceID = getAttributeValue(element, "id")
	r.Name = getAttributeValue(element, "displayName")
	r.Email = email
	r.Phone = ""
	r.Title = ""
	r.Company = ""
	r.Status = userModel.UserStatusNormal
	r.CreateTime = timestamp.StringToInt64Micro(getAttributeValue(element, "createdDate"))
	r.ModifyTime = timestamp.StringToInt64Micro(getAttributeValue(element, "updatedDate"))

	userName := getAttributeValue(element, "userName")
	p.jiraUserNameIDMap[userName] = r.ResourceID
	p.jiraUIDNameMap[r.ResourceID] = userName
	p.jiraUserEmailMap[email] = r.ResourceID

	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextWorkflow() ([]byte, error) {
	if p.workflowIndex >= len(p.jiraWorkflowsSlice) {
		p.jiraWorkflowsSlice = nil
		return nil, nil
	}
	r := p.jiraWorkflowsSlice[p.workflowIndex]
	p.workflowIndex++

	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextIssueType() ([]byte, error) {
	element, err := p.nextElement("IssueType")
	if element == nil || err != nil {
		return nil, err
	}
	r := new(resolve.ThirdIssueType)
	r.ResourceID = getAttributeValue(element, "id")
	r.Name = getAttributeValue(element, "name")
	r.Status = issuetype.StatusNormal
	r.Type = issuetype.StandardTaskType
	r.DetailType = issuetype.DetailTypeCustom
	r.Icon = 1

	style := getAttributeValue(element, "style")
	if style == jiraSubtaskType {
		r.Type = issuetype.SubTaskType
		r.Icon = 15
	}
	// 系统工作项类型映射
	if detailType, ok := p.issueTypeDetailTypeMap[r.ResourceID]; ok {
		r.DetailType = detailType
	}

	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextProjectIssueType() ([]byte, error) {
	if p.projectIssueTypeIndex >= len(p.jiraProjectIssueTypeSlice) {
		p.jiraProjectIssueTypeSlice = nil
		return nil, p.afterAction("NextProjectIssueType", nil)
	}
	r := p.jiraProjectIssueTypeSlice[p.projectIssueTypeIndex]
	p.projectIssueTypeIndex++
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextUserGroup() ([]byte, error) {
	element, err := p.nextElement("Group")
	if element == nil || err != nil {
		return nil, p.afterAction("NextUserGroup", err)
	}

	r := new(resolve.ThirdUserGroup)
	r.ResourceID = getAttributeValue(element, "id")
	r.Name = getAttributeValue(element, "groupName")
	r.Status = utilsModel.UserGroupStatusNormal
	r.CreateTime = timestamp.StringToInt64Micro(getAttributeValue(element, "createdDate"))

	p.jiraUserGroupIDNameMap[r.ResourceID] = r.Name

	if _, ok := p.jiraUserGroupNameIDMap[r.Name]; ok {
		return p.NextUserGroup()
	}

	p.jiraUserGroupNameIDMap[r.Name] = r.ResourceID
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextGlobalProjectRole() ([]byte, error) {
	element, err := p.nextElement("ProjectRole")
	if element == nil || err != nil {
		return nil, err
	}

	r := new(resolve.ThirdGlobalProjectRole)
	r.ResourceID = getAttributeValue(element, "id")
	r.Name = getAttributeValue(element, "name")
	r.Status = roleModel.RoleStatusNormal
	r.CreateTime = time.Now().Unix()

	p.jiraGlobalProjectRoleNameIDMap[r.Name] = r.ResourceID

	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextProjectRole() ([]byte, error) {
	if p.projectRoleIndex >= len(p.jiraProjectRolesSlice) {
		p.jiraProjectRolesSlice = nil
		return nil, nil
	}

	element := p.jiraProjectRolesSlice[p.projectRoleIndex]
	r := new(resolve.ThirdProjectRole)
	r.ResourceRoleID = element.RoleID
	r.ResourceProjectID = element.ProjectID
	r.CreateTime = time.Now().Unix()
	p.projectRoleIndex += 1

	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextProjectRoleMember() ([]byte, error) {
	if p.projectRoleMemberIndex >= len(p.jiraProjectMembersSlice) {
		p.jiraProjectMembersSlice = nil
		return nil, nil
	}

	element := p.jiraProjectMembersSlice[p.projectRoleMemberIndex]
	r := new(resolve.ThirdProjectRoleMember)
	r.ResourceRoleID = element.RoleID
	r.ResourceProjectID = element.ProjectID
	//r.ResourceUserID = p.jiraUserNameIDMap[element.UserName]
	r.ResourceUserID = element.UserID
	p.projectRoleMemberIndex += 1

	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextUserGroupMember() ([]byte, error) {
	element, err := p.nextElement("Membership")
	if element == nil || err != nil {
		return nil, p.afterAction("NextUserGroupMember", err)
	}
	r := new(resolve.ThirdUserGroupMember)
	r.ResourceID = getAttributeValue(element, "id")
	r.ResourceGroupID = getAttributeValue(element, "parentId")
	r.ResourceUserID = getAttributeValue(element, "childId")

	groupName, ok := p.jiraUserGroupIDNameMap[r.ResourceGroupID]
	if !ok {
		return p.NextUserGroupMember()
	}
	groupID, ok := p.jiraUserGroupNameIDMap[groupName]
	if !ok {
		return p.NextUserGroupMember()
	}

	r.ResourceGroupID = groupID

	if _, found := p.jiraUserGroupIDMembersMap[r.ResourceGroupID]; !found {
		p.jiraUserGroupIDMembersMap[r.ResourceGroupID] = []string{}
	}
	p.jiraUserGroupIDMembersMap[r.ResourceGroupID] = append(p.jiraUserGroupIDMembersMap[r.ResourceGroupID], r.ResourceUserID)
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextTaskStatus() ([]byte, bool, error) {
	element, err := p.nextElement("Status")
	if element == nil || err != nil {
		return nil, false, err
	}

	r := new(resolve.ThirdTaskStatus)

	r.ResourceID = getAttributeValue(element, "id")
	r.Name = getAttributeValue(element, "name")
	c := getAttributeValue(element, "statuscategory")
	category, ok := p.getStatusCategory(c)
	if !ok {
		return nil, true, nil
	}
	r.Category = category
	return utils2.OutputJSON(r), false, nil
}

func (p *JiraResolver) NextTaskField() ([]byte, error) {
	if p.customFieldIndex >= len(p.jiraCustomField) {
		p.jiraCustomField = nil
		return nil, nil
	}

	r := p.jiraCustomField[p.customFieldIndex]
	p.customFieldIndex += 1
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextTaskFieldOption() ([]byte, error) {
	if p.customFieldOptionIndex >= len(p.jiraCustomFieldOptions) {
		p.jiraCustomFieldOptions = nil
		return nil, nil
	}

	r := p.jiraCustomFieldOptions[p.customFieldOptionIndex]
	p.customFieldOptionIndex += 1
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextPriority() ([]byte, error) {
	element, err := p.nextElement("Priority")
	if element == nil || err != nil {
		return nil, err
	}

	o := element
	r := new(resolve.ThirdPriority)
	r.ResourceID = getAttributeValue(o, "id")
	r.Name = getAttributeValue(o, "name")
	r.Desc = getAttributeValue(o, "description")
	r.Color = "#307fe2"
	r.BackgroundColor = "#e0ecfb"
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) Config() ([]byte, error) {
	weekDays := commonModel.TeamWorkdays

	daysPerWeek, _ := strconv.Atoi(strings.TrimSpace(p.daysPerWeek))
	if daysPerWeek > 7 {
		daysPerWeek = 7
	} else if daysPerWeek <= 0 {
		daysPerWeek = 5
	}
	workdays := weekDays[:daysPerWeek]

	workhours, _ := strconv.Atoi(strings.TrimSpace(p.hoursPerDay))

	r := &resolve.ThirdConfig{
		Workdays:            workdays,
		WorkHours:           workhours,
		FixedFieldsMapping:  fixedFieldsMapping,
		CustomFieldsMapping: customFieldsMapping,
		CustomFields:        customFields,
		TabFields:           tabFields,
		LayoutFieldIDMap:    issueTypeLayoutJiraOnesFieldMap,
	}
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextTaskLinkType() ([]byte, error) {
	var index = p.issueLinkTypeIndex
	var datas = p.jiraIssueLinkType
	if index >= len(datas) {
		p.jiraIssueLinkType = nil
		return nil, nil
	}
	p.issueLinkTypeIndex += 1

	o := datas[index]
	r := new(resolve.ThirdTaskLinkType)
	r.ResourceID = getAttributeValue(o, "id")
	r.Name = getAttributeValue(o, "linkname")
	r.DescIn = getAttributeValue(o, "inward")
	r.DescOut = getAttributeValue(o, "outward")
	if r.DescIn == r.DescOut {
		r.Type = objectlinktype.LinkModelTwoWayMany
	} else {
		r.Type = objectlinktype.LinkModelManyToMany
	}
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextTaskFieldValue() ([]byte, bool, error) {
	if p.customFieldValueIndex < len(p.jiraCustomFieldValues) {
		r := p.jiraCustomFieldValues[p.customFieldValueIndex]
		p.customFieldValueIndex += 1
		if p.ignoreTaskID(r.TaskID) {
			return nil, true, nil
		}
		return utils2.OutputJSON(r), false, nil
	}

	element, err := p.nextElement("CustomFieldValue")
	if element == nil || err != nil {
		return nil, false, err
	}
	var d = element

	id := getAttributeValue(d, "id")
	issueID := getAttributeValue(d, "issue")
	if p.ignoreTaskID(issueID) {
		return nil, true, nil
	}
	fieldID := getAttributeValue(d, "customfield")

	if fieldID == p.jiraSprintCustomFieldID {
		return nil, true, nil
	}

	var value interface{}
	stringValue := getAttributeValue(d, "stringvalue")
	numberValue := getAttributeValue(d, "numbervalue")
	dateValue := getAttributeValue(d, "datevalue")
	textValue := getAttributeValue(d, "textvalue")

	if len(stringValue) > 0 {
		value = stringValue
	} else if len(numberValue) > 0 {
		value = numberValue
	} else if len(dateValue) > 0 {
		value = timestamp.StringToInt64(dateValue)
	} else if len(textValue) > 0 {
		value = textValue
	}

	fieldType, ok := p.jiraCustomFieldTypeMap[fieldID]
	if !ok {
		log.Printf("field id %s not found in field type", fieldID)
		return nil, true, nil
	}

	if _, ok := p.jiraCustomFieldCascadingSelectMap[fieldID]; ok {
		return nil, true, nil
	} else if _, ok := p.jiraCustomFieldVersionMap[fieldID]; ok {
		return nil, true, nil
	} else if _, ok := p.jiraCustomFieldProjectMap[fieldID]; ok {
		valueStr, _ := value.(string)
		k := strings.Split(valueStr, ".")[0]
		v, ok := p.jiraProjectIDNameMap[k]
		if !ok {
			return nil, true, nil
		}

		value = v
	} else if _, ok := p.jiraCustomFieldMultiGroupMap[fieldID]; ok {
		return nil, true, nil
	}

	if fieldType == fieldModel.FieldTypeUserList || fieldType == fieldModel.FieldTypeUser {
		valueStr, _ := value.(string)
		userID, ok := p.jiraUserNameIDMap[valueStr]
		if !ok {
			log.Printf("field %s user %s not found in jiraUserNameIDMap", fieldID, value)
			return nil, true, nil
		}
		value = userID
	}

	r := &resolve.ThirdTaskFieldValue{
		Base: resolve.Base{
			ResourceID: id,
		},
		TaskID:    issueID,
		FieldID:   fieldID,
		FieldType: fieldType,
		Value:     value,
	}
	return utils2.OutputJSON(r), false, nil
}

func (p *JiraResolver) NextTaskWatcher() ([]byte, bool, error) {
	element, err := p.nextElement("UserAssociation")
	if element == nil || err != nil {
		return nil, false, err
	}

	var o = element

	associationType := getAttributeValue(o, "associationType")
	if associationType != associationTypeWatchIssue {
		return nil, true, nil
	}

	r := new(resolve.ThirdTaskWatcher)
	r.TaskID = getAttributeValue(o, "sinkNodeId")
	if p.ignoreTaskID(r.TaskID) {
		return nil, true, nil
	}
	userName := getAttributeValue(o, "sourceName")
	userID := p.jiraUserNameIDMap[userName]
	r.UserID = userID
	r.ResourceID = fmt.Sprintf("%s-%s", r.TaskID, r.UserID)
	return utils2.OutputJSON(r), false, nil
}

func (p *JiraResolver) NextTaskWorkLog() ([]byte, bool, error) {
	if p.workLogIndex < len(p.jiraWorkLogs) {
		r := p.jiraWorkLogs[p.workLogIndex]
		p.workLogIndex += 1
		return utils2.OutputJSON(r), false, nil
	}

	element, err := p.nextElement("Worklog")
	if element == nil || err != nil {
		return nil, false, err
	}

	var d = element

	id := getAttributeValue(d, "id")
	issueID := getAttributeValue(d, "issue")
	if p.ignoreTaskID(issueID) {
		return nil, true, err
	}

	authorName := getAttributeValue(d, "author")
	authorID, ok := p.jiraUserNameIDMap[authorName]
	if !ok {
		log.Printf("jira work log author %s not found", authorName)
	}

	created := getAttributeValue(d, "created")
	createdTime := timestamp.StringToInt64(created)

	var timeWorked int64
	timeWorkedStr := getAttributeValue(d, "timeworked")
	if len(timeWorkedStr) > 0 {
		timeWorked, err = strconv.ParseInt(timeWorkedStr, 10, 64)
		if err != nil {
			log.Println("parse timeWorked %s failed, %+v\n", timeWorkedStr, err)
			return nil, true, err
		}
	}

	startDate := getAttributeValue(d, "startdate")
	startTime := timestamp.StringToInt64(startDate)

	r := &resolve.ThirdTaskWorkLog{
		Base:        resolve.Base{ResourceID: id},
		TaskID:      issueID,
		UserID:      authorID,
		StartTime:   startTime,
		Hours:       float64(timeWorked) / 3600,
		Type:        constants.ThirdTaskWorkLogTypeLog,
		CreateTime:  createdTime,
		Description: getAttributeValue(d, "body"),
	}
	return utils2.OutputJSON(r), false, nil
}

func versionIDToTaskID(versionID string) string {
	return fmt.Sprintf("Version-%s", versionID)
}

func versionIDToCommentID(versionID string) string {
	return fmt.Sprintf("comment-%s", versionID)
}

func (p *JiraResolver) NextTaskComment() ([]byte, bool, error) {
	if p.versionCommentIndex < len(p.jiraReleaseVersionSlice) {
		element := p.jiraReleaseVersionSlice[p.versionCommentIndex]
		r := new(resolve.ThirdTaskComment)
		versionID := getAttributeValue(element, "id")
		projectID := getAttributeValue(element, "project")
		name := getAttributeValue(element, "name")
		r.TaskID = versionIDToTaskID(versionID)
		r.ResourceID = versionIDToCommentID(versionID)
		r.UserID = p.jiraProjectIDAssignIDMap[projectID]
		r.Body = fmt.Sprintf("%s - %s - %s %s", i18n.MustGetMessage("version_note"), p.jiraProjectIDNameMap[projectID], i18n.MustGetMessage("version"), name)
		r.CreateTime = time.Now().Unix()
		p.versionCommentIndex++
		return utils2.OutputJSON(r), false, nil
	}

	element, err := p.nextElement("Action")
	if element == nil || err != nil {
		return nil, false, err
	}

	//
	o := element
	t := getAttributeValue(o, "type")
	if t != "comment" {
		return nil, true, nil
	}

	r := new(resolve.ThirdTaskComment)
	r.ResourceID = getAttributeValue(o, "id")
	r.TaskID = getAttributeValue(o, "issue")
	userName := getAttributeValue(o, "author")
	userID, ok := p.jiraUserNameIDMap[userName]
	if !ok {
		log.Printf("jira task comment author %s not found in jiraUserNameIDMap", userName)
	}
	r.UserID = userID
	r.Body = getAttributeValue(o, "body")
	created := getAttributeValue(o, "created")
	r.CreateTime = timestamp.StringToInt64(created)
	return utils2.OutputJSON(r), false, nil
}

func (p *JiraResolver) NextTaskLink() ([]byte, error) {
	var index = p.issueLinkIndex
	var datas = p.jiraIssueLinks
	if index >= len(datas) {
		p.jiraIssueLinks = nil
		return nil, p.afterAction("NextTaskLink", nil)
	}
	p.issueLinkIndex += 1

	o := datas[index]

	r := new(resolve.ThirdTaskLink)
	r.ResourceID = getAttributeValue(o, "id")
	r.LinkTypeID = getAttributeValue(o, "linktype")
	r.SourceTaskID = getAttributeValue(o, "source")
	r.TargetTaskID = getAttributeValue(o, "destination")
	return utils2.OutputJSON(r), nil
}

func (p *JiraResolver) NextTaskAttachment() ([]byte, bool, error) {
	element, err := p.nextElement("FileAttachment")
	if element == nil || err != nil {
		return nil, false, err
	}
	var o = element
	//
	r := new(resolve.ThirdTaskAttachment)
	r.ResourceID = getAttributeValue(o, "id")
	r.TaskID = getAttributeValue(o, "issue")
	r.FileName = getAttributeValue(o, "filename")

	if p.ignoreTaskID(r.TaskID) {
		return nil, true, nil
	}

	issueKey, ok := p.mapIssueIDWithKey[r.TaskID]
	if !ok {
		log.Printf("the jira issue key is not found. issue id: (%s), FileAttachment id: (%s)", r.TaskID, r.ResourceID)
		return nil, true, nil
	}

	m, ok := p.mapIssueKeyWithAttachments[issueKey]
	if !ok {
		originalIssueKey, _ := p.mapIssueIDWithOriginalKey[r.TaskID]
		m, ok = p.mapIssueKeyWithAttachments[originalIssueKey]
		if !ok {
			log.Printf("the jira issue key is not found in mapIssueKeyWithAttachments. issue id: (%s), issue key: (%s), FileAttachment id: (%s)", r.TaskID, issueKey, r.ResourceID)
			return p.NextTaskAttachment()
		}
	}

	shortFilePath, ok := m[r.ResourceID]
	if !ok {
		log.Printf("the jira attachment file is not found. issue key: (%s), FileAttachment id: (%s)", r.TaskID, r.ResourceID)
		return nil, true, nil
	}
	filePath := fmt.Sprintf("%s%s", p.importTask.AttachmentsPath, shortFilePath)
	stat, err := os.Stat(filePath)
	if err != nil {
		log.Printf("get file stat err: %s", err)
	}
	p.totalAttachmentSize += stat.Size()
	r.FilePath = shortFilePath
	return utils2.OutputJSON(r), false, nil
}

func (p *JiraResolver) TotalAttachmentSize() int64 {
	return p.totalAttachmentSize
}

func (p *JiraResolver) Clear() error {
	return p.setCache()
}

func (p *JiraResolver) getStatusCategory(category string) (string, bool) {
	var m = map[string]string{
		"1": statusModel.TaskStatusCategoryInProgressLabel,
		"2": statusModel.TaskStatusCategoryToDoLabel,
		"3": statusModel.TaskStatusCategoryDoneLabel,
		"4": statusModel.TaskStatusCategoryInProgressLabel,
		"5": statusModel.TaskStatusCategoryInProgressLabel,
		"":  statusModel.TaskStatusCategoryInProgressLabel,
	}
	resp, ok := m[category]
	return resp, ok
}

func (p *JiraResolver) getFieldType(fieldType string) (int, bool) {
	var m = map[string]int{
		"float":            fieldModel.FieldTypeFloat,
		"textfield":        fieldModel.FieldTypeText,
		"radiobuttons":     fieldModel.FieldTypeOption,
		"textarea":         fieldModel.FieldTypeMultiLineText,
		"datetime":         fieldModel.FieldTypeTime,
		"userpicker":       fieldModel.FieldTypeUser,
		"datepicker":       fieldModel.FieldTypeDate,
		"multiselect":      fieldModel.FieldTypeMultiOption,
		"labels":           fieldModel.FieldTypeText,
		"multicheckboxes":  fieldModel.FieldTypeMultiOption,
		"select":           fieldModel.FieldTypeOption,
		"cascadingselect":  fieldModel.FieldTypeText,
		"project":          fieldModel.FieldTypeText,
		"multiuserpicker":  fieldModel.FieldTypeUserList,
		"grouppicker":      fieldModel.FieldTypeText,
		"multigrouppicker": fieldModel.FieldTypeText,
		"readonlyfield":    fieldModel.FieldTypeMultiLineText,
		"multiversion":     fieldModel.FieldTypeText,
		"version":          fieldModel.FieldTypeText,
		"jobcheckbox":      fieldModel.FieldTypeText,
		"hiddenjobswitch":  fieldModel.FieldTypeText,
		"devsummary":       fieldModel.FieldTypeText,
		"url":              fieldModel.FieldTypeText,
		"gh-lexo-rank":     fieldModel.FieldTypeText,
		"gh-epic-label":    fieldModel.FieldTypeText,
		"gh-epic-status":   fieldModel.FieldTypeOption,
		"greenhopper-releasedmultiversionhistory": fieldModel.FieldTypeText,
	}
	resp, ok := m[fieldType]
	return resp, ok
}

func (p *JiraResolver) readersFromFile(path string) (io.ReadCloser, io.ReadCloser, error) {
	zipFile, err := zip.OpenReader(path)
	if err != nil {
		return nil, nil, err
	}

	var entityFile, activeObjectFile *zip.File
	for _, f := range zipFile.File {
		if entityFile != nil && activeObjectFile != nil {
			break
		}
		if f.FileInfo().IsDir() {
			continue
		}
		fileName := f.FileInfo().Name()
		if fileName == entitiesFile {
			entityFile = f
			continue
		}
		if fileName == activeObjectsFile {
			activeObjectFile = f
		}
	}
	if entityFile == nil {
		err = fmt.Errorf("%s file not found", entitiesFile)
		return nil, nil, err
	}
	if activeObjectFile == nil {
		err = fmt.Errorf("%s file not found", activeObjectsFile)
		return nil, nil, err
	}

	entityFileReader, e := entityFile.Open()
	if e != nil && e != io.EOF {
		return nil, nil, e
	}

	activeObjectsFileReader, e := activeObjectFile.Open()
	if e != nil && e != io.EOF {
		return nil, nil, e
	}

	return entityFileReader, activeObjectsFileReader, nil
}

func (p *JiraResolver) getCustomField() ([]*resolve.ThirdTaskField, error) {
	var fields = make([]*resolve.ThirdTaskField, 0)
	var err error

	fields = append(fields, &resolve.ThirdTaskField{
		Base: resolve.Base{
			ResourceID: customFieldReleaseStartDate,
		},
		Name: i18n.MustGetMessage("field.release_start_date"),
		Type: fieldModel.FieldTypeDate,
	})
	fields = append(fields, &resolve.ThirdTaskField{
		Base: resolve.Base{
			ResourceID: customFieldID,
		},
		Name: customFieldID,
		Type: fieldModel.FieldTypeText,
	})
	fields = append(fields, &resolve.ThirdTaskField{
		Base: resolve.Base{
			ResourceID: customFieldLabels,
		},
		Name: customFieldLabels,
		Type: fieldModel.FieldTypeText,
	})
	fields = append(fields, &resolve.ThirdTaskField{
		Base: resolve.Base{
			ResourceID: customFieldFixVersion,
		},
		Name: customFieldFixVersion,
		Type: fieldModel.FieldTypeText,
	})
	fields = append(fields, &resolve.ThirdTaskField{
		Base: resolve.Base{
			ResourceID: customFieldVersion,
		},
		Name: customFieldVersion,
		Type: fieldModel.FieldTypeText,
	})
	fields = append(fields, &resolve.ThirdTaskField{
		Base: resolve.Base{
			ResourceID: customFieldResolution,
		},
		Name: customFieldResolution,
		Type: fieldModel.FieldTypeOption,
	})
	fields = append(fields, &resolve.ThirdTaskField{
		Base: resolve.Base{
			ResourceID: customFieldComponent,
		},
		Name: customFieldComponent,
		Type: fieldModel.FieldTypeText,
	})
	fields = append(fields, &resolve.ThirdTaskField{
		Base: resolve.Base{
			ResourceID: customFieldEnvironment,
		},
		Name: customFieldEnvironment,
		Type: fieldModel.FieldTypeMultiLineText,
	})

	for {
		element, e := p.nextElement("CustomField")
		if element == nil || e != nil {
			err = e
			break
		}

		var f = element

		name := getAttributeValue(f, "name")

		var k = getAttributeValue(f, "customfieldtypekey")
		parts := strings.Split(k, ":")
		if len(parts) != 2 {
			continue
		}

		t := parts[1]

		id := getAttributeValue(f, "id")

		if name == "Rank" && t == "gh-lexo-rank" {
			continue
		}

		if name == "Sprint" && t == "gh-sprint" {
			p.jiraSprintCustomFieldID = id
			continue
		}

		if t == "cascadingselect" {
			p.jiraCustomFieldCascadingSelectMap[id] = ""
		} else if t == "version" || t == "multiversion" {
			p.jiraCustomFieldVersionMap[id] = ""
		} else if t == "project" {
			p.jiraCustomFieldProjectMap[id] = ""
		} else if t == "multigrouppicker" {
			p.jiraCustomFieldMultiGroupMap[id] = ""
		}

		o := new(resolve.ThirdTaskField)
		o.ResourceID = id
		o.Name = name
		fieldType, ok := p.getFieldType(t)
		if !ok {
			continue
		}

		p.jiraCustomFieldTypeMap[id] = fieldType

		o.Type = fieldType
		fields = append(fields, o)
	}

	for {
		element, e := p.nextElement("Component")
		if element == nil || e != nil {
			err = e
			break
		}
		var o = element

		id := getAttributeValue(o, "id")
		name := getAttributeValue(o, "name")
		p.jiraComponentNameMap[id] = name
	}
	return fields, err
}

func (p *JiraResolver) getCustomFieldOption() ([]*resolve.ThirdTaskFieldOption, error) {
	var datas = make([]*resolve.ThirdTaskFieldOption, 0)
	var err error
	for {
		element, e := p.nextElement("CustomFieldOption")
		if element == nil || e != nil {
			err = e
			break
		}
		var op = element
		d := new(resolve.ThirdTaskFieldOption)
		d.ResourceID = getAttributeValue(op, "id")
		d.ResourceTaskFieldID = getAttributeValue(op, "customfield")
		d.Name = getAttributeValue(op, "value")
		d.BackgroundColor = "#307fe2"
		d.Color = "#fff"

		if _, ok := p.jiraCustomFieldCascadingSelectMap[d.ResourceTaskFieldID]; ok {
			m, ok := p.jiraCustomFieldCascadingSelectOptionMap[d.ResourceTaskFieldID]
			if !ok {
				m = make(map[string]string)
			}
			m[d.ResourceID] = d.Name
			p.jiraCustomFieldCascadingSelectOptionMap[d.ResourceTaskFieldID] = m
		}

		datas = append(datas, d)
	}

	for {
		element, e := p.nextElement("Resolution")
		if element == nil || e != nil {
			err = e
			break
		}
		var r = element
		d := new(resolve.ThirdTaskFieldOption)
		d.ResourceID = fmt.Sprintf("%s-%s", customFieldResolution, getAttributeValue(r, "id"))
		d.ResourceTaskFieldID = customFieldResolution
		d.Name = getAttributeValue(r, "name")
		d.BackgroundColor = "#307fe2"
		d.Color = "#fff"
		datas = append(datas, d)
	}

	datas = append(datas, &resolve.ThirdTaskFieldOption{
		Base: resolve.Base{
			ResourceID: fmt.Sprintf("%s-%s", customFieldResolution, "Unresolved"),
		},
		ResourceTaskFieldID: customFieldResolution,
		Name:                "Unresolved",
		Color:               "#fff",
		BackgroundColor:     "#307fe2",
	})
	return datas, err
}

func (p *JiraResolver) getIssueLinks() ([]*etree.Element, error) {
	var resp []*etree.Element
	var err error
	for {
		element, e := p.nextElement("IssueLink")
		if element == nil || e != nil {
			err = e
			break
		}

		var link = element
		parentID := getAttributeValue(link, "source")
		childID := getAttributeValue(link, "destination")
		linkTypeID := getAttributeValue(link, "linktype")
		if linkTypeID == p.jiraSubTaskLinkID {
			p.jiraIssueParentIDMap[childID] = parentID
		} else {
			resp = append(resp, link)
		}
	}
	return resp, err
}

func (p *JiraResolver) getIssueLinkTypes() ([]*etree.Element, error) {
	var resp []*etree.Element
	var err error
	for {
		element, e := p.nextElement("IssueLinkType")
		if element == nil || e != nil {
			err = e
			break
		}
		var linkType = element
		linkname := getAttributeValue(linkType, "linkname")
		if linkname == "jira_subtask_link" {
			p.jiraSubTaskLinkID = getAttributeValue(linkType, "id")
		} else {
			resp = append(resp, linkType)
		}
	}
	return resp, err
}

func (p *JiraResolver) preprocessIssues() error {
	dir := importerFileDir()
	parentPath := fmt.Sprintf("%s_*.xml", "IssueParent")
	parentFile, err := ioutil.TempFile(dir, parentPath)
	if err != nil {
		return err
	}
	defer parentFile.Close()

	childPath := fmt.Sprintf("%s_*.xml", "IssueChild")
	childFile, err := ioutil.TempFile(dir, childPath)
	if err != nil {
		return err
	}
	defer childFile.Close()

	for {
		line, e := p.nextLineFromTagFile("Issue")
		if line == "" || e != nil {
			break
		}

		doc := etree.NewDocument()
		if err := doc.ReadFromString(line); err != nil {
			return err
		}

		var o = doc.Root()
		projectID := getAttributeValue(o, "project")
		if p.ignoreProjectID(projectID) {
			continue
		}
		childID := getAttributeValue(o, "id")
		if childID == "" {
			continue
		}
		issueTypeID := getAttributeValue(o, "type")
		environment := getAttributeValue(o, "environment")
		if environment != "" {
			r := &resolve.ThirdTaskFieldValue{
				Base: resolve.Base{
					ResourceID: fmt.Sprintf("%s-%s", childID, customFieldEnvironment),
				},
				TaskID:    childID,
				FieldID:   customFieldEnvironment,
				FieldType: fieldModel.FieldTypeMultiLineText,
				Value:     environment,
			}
			p.jiraCustomFieldValues = append(p.jiraCustomFieldValues, r)
		}
		issueID := childID
		sprintID, ok := p.jiraIssueSprintMap[issueID]
		if ok {
			p.jiraSprintIDProjectIDMap[sprintID] = projectID
		}

		parentIssueID, ok := p.jiraIssueParentIDMap[childID]
		if ok {
			// child
			childFile.WriteString(line + "\n")
		} else {
			// parent
			parentFile.WriteString(line + "\n")
		}
		if parentIssueID != "" {
			p.jiraTaskIDMap[parentIssueID] = struct{}{}
		}
		p.jiraTaskIDMap[childID] = struct{}{}

		mapKey := fmt.Sprintf("%s:%s", projectID, issueTypeID)
		if _, found := p.jiraProjectIssueTypeMap[mapKey]; !found {
			p.jiraProjectIssueTypeMap[mapKey] = struct{}{}
		}
	}

	p.IssueChildFilePath = childFile.Name()
	p.IssueParentFilePath = parentFile.Name()

	if err := parentFile.Sync(); err != nil {
		return err
	}
	if err := childFile.Sync(); err != nil {
		return err
	}

	cf, err := os.Open(childFile.Name())
	if err != nil {
		return err
	}
	p.IssueChildReader = resolve.NewXmlScanner(cf, entityRootTag)

	pf, err := os.Open(parentFile.Name())
	if err != nil {
		return err
	}
	p.IssueParentReader = resolve.NewXmlScanner(pf, entityRootTag)
	return err
}

func (p *JiraResolver) getCustomFieldValues() ([]*resolve.ThirdTaskFieldValue, error) {
	resp := make([]*resolve.ThirdTaskFieldValue, 0)
	var err error

	cascadingSelectFieldValue := make(map[string]map[string][]string)
	versionFieldValue := make(map[string]map[string][]string)
	multiGroupFieldValue := make(map[string]map[string][]string)

	var tag = "CustomFieldValue"
	count := 0

	log.Println("[jira import] start getCustomFieldValues1")
	for {
		elements, err := p.nextElements(tag, 1000000)
		if err != nil {
			return nil, err
		}
		if elements == nil || len(elements) == 0 {
			break
		}
		count += len(elements)
		log.Printf(" CustomFieldValue count: %d", count)

		for _, e := range elements {
			var d = e

			issueID := getAttributeValue(d, "issue")
			if len(issueID) == 0 {
				continue
			}

			fieldID := getAttributeValue(d, "customfield")

			var value string
			stringValue := getAttributeValue(d, "stringvalue")
			numberValue := getAttributeValue(d, "numbervalue")
			dateValue := getAttributeValue(d, "datevalue")
			textValue := getAttributeValue(d, "textvalue")

			if len(stringValue) > 0 {
				value = stringValue
			} else if len(numberValue) > 0 {
				value = numberValue
			} else if len(dateValue) > 0 {
				value = dateValue
			} else if len(textValue) > 0 {
				value = textValue
			}

			if fieldID == p.jiraSprintCustomFieldID {
				p.jiraIssueSprintMap[issueID] = value
				continue
			}

			// 如果是 cascadingselect, Version, multiversion 属性类型，要先特殊处理
			if _, ok := p.jiraCustomFieldCascadingSelectMap[fieldID]; ok {
				m, ok := p.jiraCustomFieldCascadingSelectOptionMap[fieldID]
				if !ok {
					continue
				}
				v, ok := m[value]
				if !ok {
					continue
				}

				m2, ok := cascadingSelectFieldValue[issueID]
				if !ok {
					m2 = make(map[string][]string)
				}
				m2[fieldID] = append(m2[fieldID], v)
				cascadingSelectFieldValue[issueID] = m2
			} else if _, ok := p.jiraCustomFieldVersionMap[fieldID]; ok {
				value = strings.Split(value, ".")[0]
				v, ok := p.jiraReleaseVersionMap[value]
				if !ok {
					continue
				}

				m2, ok := versionFieldValue[issueID]
				if !ok {
					m2 = make(map[string][]string)
				}
				m2[fieldID] = append(m2[fieldID], v)
				versionFieldValue[issueID] = m2
			} else if _, ok := p.jiraCustomFieldMultiGroupMap[fieldID]; ok {
				m2, ok := multiGroupFieldValue[issueID]
				if !ok {
					m2 = make(map[string][]string)
				}
				m2[fieldID] = append(m2[fieldID], value)
				multiGroupFieldValue[issueID] = m2
			}
		}
	}

	log.Println("[jira import] end getCustomFieldValues1")
	// reset customfieldvalue tag file
	path, ok := p.mapTagFilePath[tag]
	if !ok {
		return nil, errors.New("%s tag file path not found")
	}

	newReader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	p.tagFilesMap[tag] = resolve.NewXmlScanner(newReader, entityRootTag)

	//
	for issueID, mapFieldIDValues := range cascadingSelectFieldValue {
		for fieldID, vs := range mapFieldIDValues {
			r := &resolve.ThirdTaskFieldValue{
				Base: resolve.Base{
					ResourceID: fmt.Sprintf("%s-%s", issueID, fieldID),
				},
				TaskID:    issueID,
				FieldID:   fieldID,
				FieldType: fieldModel.FieldTypeText,
				Value:     strings.Join(vs, ","),
			}
			resp = append(resp, r)
		}
	}
	cascadingSelectFieldValue = nil

	for issueID, mapFieldIDValues := range versionFieldValue {
		for fieldID, vs := range mapFieldIDValues {
			r := &resolve.ThirdTaskFieldValue{
				Base: resolve.Base{
					ResourceID: fmt.Sprintf("%s-%s", issueID, fieldID),
				},
				TaskID:    issueID,
				FieldID:   fieldID,
				FieldType: fieldModel.FieldTypeText,
				Value:     strings.Join(vs, ","),
			}
			resp = append(resp, r)
		}
	}
	versionFieldValue = nil

	for issueID, mapFieldIDValues := range multiGroupFieldValue {
		for fieldID, vs := range mapFieldIDValues {
			r := &resolve.ThirdTaskFieldValue{
				Base: resolve.Base{
					ResourceID: fmt.Sprintf("%s-%s", issueID, fieldID),
				},
				TaskID:    issueID,
				FieldID:   fieldID,
				FieldType: fieldModel.FieldTypeText,
				Value:     strings.Join(vs, ","),
			}
			resp = append(resp, r)
		}
	}
	multiGroupFieldValue = nil

	for issueID, values := range p.jiraIssueComponent {
		r := &resolve.ThirdTaskFieldValue{
			Base: resolve.Base{
				ResourceID: fmt.Sprintf("%s-%s", issueID, customFieldComponent),
			},
			TaskID:    issueID,
			FieldID:   customFieldComponent,
			FieldType: fieldModel.FieldTypeText,
			Value:     strings.Join(values, ","),
		}
		resp = append(resp, r)
	}
	p.jiraIssueComponent = nil

	for _, e := range p.jiraReleaseVersionSlice {
		date := getAttributeValue(e, "startdate")
		if date != "" {
			startDate := timestamp.StringToInt64(date)
			taskID := versionIDToTaskID(getAttributeValue(e, "id"))
			r := &resolve.ThirdTaskFieldValue{
				Base: resolve.Base{
					ResourceID: fmt.Sprintf("%s-%s", taskID, customFieldReleaseStartDate),
				},
				TaskID:    taskID,
				FieldID:   customFieldReleaseStartDate,
				FieldType: fieldModel.FieldTypeDate,
				Value:     startDate,
			}
			resp = append(resp, r)
		}

		date = getAttributeValue(e, "releasedate")
		if date == "" {
			date = timestamp.SecToDateTimeString()
		}
		releaseDate := timestamp.StringToInt64(date)
		taskID := versionIDToTaskID(getAttributeValue(e, "id"))
		r := &resolve.ThirdTaskFieldValue{
			Base: resolve.Base{
				ResourceID: fmt.Sprintf("%s-%s", taskID, customFieldReleaseDate),
			},
			TaskID:    taskID,
			FieldID:   fieldModel.PublishDateFieldUUID,
			FieldType: fieldModel.FieldTypeDate,
			Value:     releaseDate,
		}
		resp = append(resp, r)

	}
	p.jiraIssueFixVersion = nil
	//
	for issueID, values := range p.jiraIssueAffectsVersion {
		r := &resolve.ThirdTaskFieldValue{
			Base: resolve.Base{
				ResourceID: fmt.Sprintf("%s-%s", issueID, customFieldVersion),
			},
			TaskID:    issueID,
			FieldID:   customFieldVersion,
			FieldType: fieldModel.FieldTypeText,
			Value:     strings.Join(values, ","),
		}
		resp = append(resp, r)
	}
	p.jiraIssueAffectsVersion = nil

	issueLabelsMap := make(map[string][]string, 0)
	userDefinedIssueLabelsMap := make(map[string][]string, 0)

	for {
		element, e := p.nextElement("Label")
		if element == nil || e != nil {
			err = e
			break
		}
		var d = element

		issueID := getAttributeValue(d, "issue")
		if len(issueID) == 0 {
			continue
		}
		label := getAttributeValue(d, "label")
		fieldID := getAttributeValue(d, "fieldid")

		if len(fieldID) == 0 {
			issueLabelsMap[issueID] = append(issueLabelsMap[issueID], label)
			continue
		}
		key := fmt.Sprintf("%s-%s", issueID, fieldID)
		userDefinedIssueLabelsMap[key] = append(userDefinedIssueLabelsMap[key], label)
	}

	for issueID, labels := range issueLabelsMap {
		r := &resolve.ThirdTaskFieldValue{
			Base: resolve.Base{
				ResourceID: fmt.Sprintf("%s-%s", issueID, customFieldLabels),
			},
			TaskID:    issueID,
			FieldID:   customFieldLabels,
			FieldType: fieldModel.FieldTypeText,
			Value:     strings.Join(labels, ","),
		}
		resp = append(resp, r)
	}
	for key, labels := range userDefinedIssueLabelsMap {
		keys := strings.Split(key, "-")
		if len(keys[0]) == 0 || len(keys[1]) == 0 {
			log.Println("userDefinedIssueLabelsMap empty", keys)
			continue
		}
		r := &resolve.ThirdTaskFieldValue{
			Base: resolve.Base{
				ResourceID: fmt.Sprintf("%s-%s", keys[0], userDefinedFieldLabels),
			},
			TaskID:    keys[0],
			FieldID:   keys[1],
			FieldType: fieldModel.FieldTypeText,
			Value:     strings.Join(labels, ","),
		}
		resp = append(resp, r)
	}
	log.Println("[jira import] end getCustomFieldValues")
	return resp, err
}

func (p *JiraResolver) handleIssueResolution(o *etree.Element) {
	issueID := getAttributeValue(o, "id")
	resolution := getAttributeValue(o, "resolution")
	if len(resolution) == 0 {
		resolution = "Unresolved"
	}
	f := &resolve.ThirdTaskFieldValue{
		Base: resolve.Base{
			ResourceID: fmt.Sprintf("%s-%s", issueID, customFieldResolution),
		},
		TaskID:    issueID,
		FieldID:   customFieldResolution,
		FieldType: fieldModel.FieldTypeOption,
		Value:     fmt.Sprintf("%s-%s", customFieldResolution, resolution),
	}
	p.jiraCustomFieldValues = append(p.jiraCustomFieldValues, f)
}

func (p *JiraResolver) handleIssueKey(o *etree.Element) {
	issueID := getAttributeValue(o, "id")
	key := getAttributeValue(o, "key")
	projectID := getAttributeValue(o, "project")
	number := getAttributeValue(o, "number")
	if len(key) == 0 {
		// 兼容 8.x
		projectKey := getAttributeValue(o, "projectKey")
		key = fmt.Sprintf("%s-%s", projectKey, number)
	}

	if len(key) > 0 {
		f := &resolve.ThirdTaskFieldValue{
			Base: resolve.Base{
				ResourceID: key,
			},
			TaskID:    issueID,
			FieldID:   customFieldID,
			FieldType: fieldModel.FieldTypeText,
			Value:     key,
		}
		p.jiraCustomFieldValues = append(p.jiraCustomFieldValues, f)

		//
		p.mapIssueIDWithKey[issueID] = key
	}
	originalKey := fmt.Sprintf("%s-%s", p.jiraProjectIDOriginalKeyMap[projectID], number)
	p.mapIssueIDWithOriginalKey[issueID] = originalKey
}

func (p *JiraResolver) getVersion() (map[string]string, error) {
	var resp = make(map[string]string)
	var err error
	for {
		element, e := p.nextElement("Version")
		if element == nil || e != nil {
			err = e
			break
		}
		var o = element
		projectID := getAttributeValue(o, "project")
		if p.ignoreProjectID(projectID) {
			continue
		}
		id := getAttributeValue(o, "id")
		name := getAttributeValue(o, "name")
		resp[id] = name
		p.jiraReleaseVersionSlice = append(p.jiraReleaseVersionSlice, element)
	}
	return resp, err
}

func (p *JiraResolver) getNodeAssociation() ([]*etree.Element, error) {
	var resp []*etree.Element
	var err error
	for {
		element, e := p.nextElement("NodeAssociation")
		if element == nil || e != nil {
			err = e
			break
		}

		var o = element

		sourceNodeId := getAttributeValue(o, "sourceNodeId")
		sinkNodeId := getAttributeValue(o, "sinkNodeId")
		associationType := getAttributeValue(o, "associationType")
		sourceNodeEntity := getAttributeValue(o, "sourceNodeEntity")
		sinkNodeEntity := getAttributeValue(o, "sinkNodeEntity")
		key := fmt.Sprintf("%s:%s", sourceNodeEntity, sinkNodeEntity)
		switch key {
		case "Project:PermissionScheme":
			if _, found := p.jiraPermissionProjectAssociation[sinkNodeId]; !found {
				p.jiraPermissionProjectAssociation[sinkNodeId] = make([]string, 0)
			}
			p.jiraPermissionProjectAssociation[sinkNodeId] = append(p.jiraPermissionProjectAssociation[sinkNodeId], sourceNodeId)
		case "Project:WorkflowScheme":
			if _, found := p.jiraWorkflowSchemeProjectAssoc[sinkNodeId]; !found {
				p.jiraWorkflowSchemeProjectAssoc[sinkNodeId] = make([]string, 0)
			}
			p.jiraWorkflowSchemeProjectAssoc[sinkNodeId] = append(p.jiraWorkflowSchemeProjectAssoc[sinkNodeId], sourceNodeId)
		case "Project:NotificationScheme":
			if _, found := p.jiraNotificationSchemeProjectAssoc[sinkNodeId]; !found {
				p.jiraNotificationSchemeProjectAssoc[sinkNodeId] = make([]string, 0)
			}
			p.jiraNotificationSchemeProjectAssoc[sinkNodeId] = append(p.jiraNotificationSchemeProjectAssoc[sinkNodeId], sourceNodeId)
		case "Project:FieldLayoutScheme":
			p.jiraProjectWithFieldLayoutScheme[sourceNodeId] = sinkNodeId
		case "Project:IssueTypeScreenScheme":
			p.jiraProjectWithIssueTypeScreenScheme[sourceNodeId] = sinkNodeId
		case "Project:ProjectCategory":
			p.jiraCategoryProjectAssociation[sourceNodeId] = sinkNodeId
		}

		switch associationType {
		case associationTypeIssueComponent:
			n, ok := p.jiraComponentNameMap[sinkNodeId]
			if !ok {
				continue
			}
			p.jiraIssueComponent[sourceNodeId] = append(p.jiraIssueComponent[sourceNodeId], n)
		case associationTypeIssueVersion:
			n, ok := p.jiraReleaseVersionMap[sinkNodeId]
			if !ok {
				continue
			}
			p.jiraIssueAffectsVersion[sourceNodeId] = append(p.jiraIssueAffectsVersion[sourceNodeId], n)
		case associationTypeIssueFixVersion:
			p.jiraIssueFixVersion[sinkNodeId] = append(p.jiraIssueFixVersion[sinkNodeId], sourceNodeId)
		default:
			resp = append(resp, o)
		}
	}

	for releaseID, taskIDs := range p.jiraIssueFixVersion {
		for _, taskID := range taskIDs {
			p.jiraTaskReleaseSlice = append(p.jiraTaskReleaseSlice, &resolve.ThirdTaskRelease{
				TaskID:    taskID,
				ReleaseID: versionIDToTaskID(releaseID),
			})
		}
	}
	p.jiraIssueFixVersion = nil

	return resp, err
}

func (p *JiraResolver) processAttachment() error {
	log.Println("[jira import] start process attachment")
	files := make([]string, 0)
	err := filepath.Walk(p.importTask.AttachmentsPath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return errors.Trace(err)
	}
	for _, file := range files {
		shortName := strings.ReplaceAll(file, p.importTask.AttachmentsPath, "")
		if strings.Contains(file, ".") {
			continue
		}
		splits := strings.Split(shortName, "/")
		if len(splits) != 5 {
			continue
		}
		issueKey := splits[3]
		attachmentID := splits[4]
		m, ok := p.mapIssueKeyWithAttachments[issueKey]
		if !ok {
			m = make(map[string]string)
		}
		m[attachmentID] = shortName
		p.mapIssueKeyWithAttachments[issueKey] = m
	}
	log.Println("[jira import] end process attachment")
	return nil
}
