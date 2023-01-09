package issue_type

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/gin-contrib/i18n"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/services"
	"github.com/bangwork/import-tools/serve/services/cache"
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

func GetIssueTypeList(typeList *services.IssueTypeListResponse, issueTypes map[string]bool) (*services.IssueTypeListResponse, error) {
	list, err := cache.GetCacheInfo()
	if err != nil {
		return nil, err
	}
	v, found := list.MapFilePath[common.TagIssueType]
	if !found {
		return nil, common.Errors(common.NotFoundError, nil)
	}
	file, err := os.Open(v)
	if err != nil {
		return nil, err
	}

	mapBind := map[string]int{}
	for _, v := range typeList.JiraList {
		mapBind[v.IssueTypeID] = v.ONESDetailType
	}
	xmlReader := resolve.NewXmlScanner(file, common.TagEntityRoot)
	jiraList := make([]*services.JiraIssueType, 0)
	for {
		reader, err := xml.NextElementFromReader(xmlReader)
		if err != nil {
			log.Printf("NextElementFromReader error, %+v", err)
			return nil, err
		}
		if reader == nil {
			break
		}
		data := new(services.JiraIssueType)
		data.IssueTypeID = xml.GetAttributeValue(reader, "id")
		data.IssueTypeName = xml.GetAttributeValue(reader, "name")
		if !issueTypes[data.IssueTypeID] {
			continue
		}

		detailType, found := mapBind[data.IssueTypeID]
		if found {
			data.ONESDetailType = detailType
		}
		jiraList = append(jiraList, data)
	}
	typeList.ONESList = append(typeList.ONESList, &services.ONESIssueType{
		DetailType: common.IssueTypeDetailTypeCustom,
		Name:       i18n.MustGetMessage("issue_type_custom"),
	})

	ONESList := typeList.ONESList
	mapONESList := make(map[int]*services.ONESIssueType)
	for _, issueType := range ONESList {
		if issueType.DetailType != 0 {
			issueType.Name = i18n.MustGetMessage(fmt.Sprintf("issue_type_detail_type_%d", issueType.DetailType))
		}
		mapONESList[issueType.DetailType] = issueType
	}

	resp := new(services.IssueTypeListResponse)
	for _, detailType := range sortList {
		resp.ONESList = append(resp.ONESList, mapONESList[detailType])
	}
	sort.SliceStable(jiraList, func(i, j int) bool {
		return jiraList[i].ONESDetailType > jiraList[j].ONESDetailType
	})
	resp.JiraList = jiraList
	return resp, nil
}
