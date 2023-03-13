package issue_type

import (
	"log"
	"os"
	"sort"

	"github.com/juju/errors"

	"github.com/bangwork/import-tools/serve/models/ones"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/services"
	"github.com/bangwork/import-tools/serve/services/importer/resolve"
	"github.com/bangwork/import-tools/serve/utils/xml"
)

const (
	DetailTypeDemand  = 1 // demand
	DetailTypeTask    = 2 // task
	DetailTypeDefect  = 3 // defect
	DetailTypeSubTask = 4 // sub_task
)

var (
	SupportDetailType = map[int]bool{
		DetailTypeDemand:  true,
		DetailTypeTask:    true,
		DetailTypeDefect:  true,
		DetailTypeSubTask: true,
	}
	sortList = []int{
		common.IssueTypeDetailTypeCustom,
		common.IssueTypeDetailTypeDefect,
		common.IssueTypeDetailTypeTask,
		common.IssueTypeDetailTypeDemand,
		common.IssueTypeDetailTypeSubTask,
	}
)

func GetUnBoundIssueTypes(cookie, url, teamUUID string, h map[string]string) ([]*ones.UnBoundIssueTypes, error) {
	list, err := ones.GetUnBoundIssueTypes(url, teamUUID, h)
	if err != nil {
		return nil, errors.Trace(err)
	}
	//c, err := ones.DecryptCookieValueByCookie(cookie)
	//if err != nil {
	//	return nil, errors.Trace(err)
	//}
	//name := "ONES custom issue type"
	//switch c.Language {
	//case ones.LanguageTagChinese:
	//	name = "ONES 自定义工作项类型"
	//}
	//list = append(list, &ones.UnBoundIssueTypes{
	//	DetailType: common.IssueTypeDetailTypeCustom,
	//	Name:       name,
	//})
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].DetailType < list[j].DetailType
	})
	return list, nil
}

func GetIssueTypeList(url, teamUUID, cookie string, h map[string]string, projectIDs []string) (*services.IssueTypeListResponse, error) {
	importCache := common.ImportCacheMap.Get(cookie)
	issueTypes := make(map[string]bool)
	for _, pid := range projectIDs {
		for _, issueTypeID := range importCache.ProjectIssueTypeMap[pid] {
			issueTypes[issueTypeID] = true
		}
	}
	mapBind, err := ones.MapThirdIssueTypeBind(url, teamUUID, h)
	if err != nil {
		return nil, errors.Trace(err)
	}
	issueTypeFile, found := importCache.MapFilePath[common.TagIssueType]
	if !found {
		return nil, errors.Trace(common.Errors(common.NotFoundError, nil))
	}
	file, err := os.Open(issueTypeFile)
	if err != nil {
		return nil, errors.Trace(err)
	}

	xmlReader := resolve.NewXmlScanner(file, common.TagEntityRoot)
	migratedList := make([]*services.MigratedList, 0)
	readyMigrateList := make([]*services.ReadyMigrateList, 0)
	for {
		reader, err := xml.NextElementFromReader(xmlReader)
		if err != nil {
			log.Printf("NextElementFromReader error, %+v", err)
			return nil, errors.Trace(err)
		}
		if reader == nil {
			break
		}
		data := new(services.MigratedList)
		data.IssueTypeID = xml.GetAttributeValue(reader, "id")
		data.IssueTypeName = xml.GetAttributeValue(reader, "name")
		style := xml.GetAttributeValue(reader, "style")
		if !issueTypes[data.IssueTypeID] {
			continue
		}
		issueTypeTaskType := ones.IssueTypeStandardTaskType
		if style == ones.JiraSubTaskStyle {
			issueTypeTaskType = ones.IssueTypeSubTaskType
		}

		bind, found := mapBind[data.IssueTypeID]
		if found {
			data.ONESIssueTypeName = bind.ONESIssueTypeName
			data.Action = "map"
			if bind.ONESDetailType == -1 {
				data.Action = "create"
			}
			migratedList = append(migratedList, data)
		} else {
			readyMigrateList = append(readyMigrateList, &services.ReadyMigrateList{
				IssueTypeID:   data.IssueTypeID,
				IssueTypeName: data.IssueTypeName,
				Type:          issueTypeTaskType,
			})
		}
	}

	resp := new(services.IssueTypeListResponse)
	resp.MigratedList = migratedList
	resp.ReadyMigrateList = readyMigrateList
	return resp, nil
}
