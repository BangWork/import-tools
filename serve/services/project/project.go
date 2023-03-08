package project

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/bangwork/import-tools/serve/utils"

	pinyin "github.com/mozillazg/go-pinyin"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/services/importer/resolve"
	"github.com/bangwork/import-tools/serve/utils/xml"
)

type Project struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	OriginalKey string `json:"original_key"`
	Name        string `json:"name"`
	Assign      string `json:"assign"`
	Category    string `json:"category"`
	IssueCount  int    `json:"issue_count"`
	Type        string `json:"type"`
}

type Response struct {
	ProjectList []*Project `json:"projects"`
}

func GetProjectList(cookie string) (*Response, error) {
	importCache := common.ImportCacheMap.Get(cookie)
	v, found := importCache.MapFilePath[common.TagProject]
	if !found {
		return nil, common.Errors(common.NotFoundError, nil)
	}
	file, err := os.Open(v)
	if err != nil {
		return nil, err
	}
	xmlReader := resolve.NewXmlScanner(file, common.TagEntityRoot)
	res := make([]*Project, 0)
	for {
		reader, err := xml.NextElementFromReader(xmlReader)
		if err != nil {
			log.Printf("NextElementFromReader error, %+v", err)
			return nil, err
		}
		if reader == nil {
			break
		}
		data := new(Project)
		data.ID = xml.GetAttributeValue(reader, "id")
		data.Name = xml.GetAttributeValue(reader, "name")
		data.Type = xml.GetAttributeValue(reader, "projecttype")
		data.Key = xml.GetAttributeValue(reader, "key")
		data.OriginalKey = xml.GetAttributeValue(reader, "originalkey")
		data.IssueCount = xml.GetAttributeValueInt(reader, "counter")
		assign, found := importCache.ProjectAssignMap[data.ID]
		if !found {
			log.Println("ProjectAssignMap not found", data.ID)
		}
		data.Assign = assign
		category, found := importCache.ProjectCategoryMap[data.ID]
		if !found {
			log.Println("ProjectCategoryMap not found", data.ID)
		}
		data.Category = category
		res = append(res, data)
	}
	pinyinArgs := pinyin.NewArgs()
	sort.Slice(res, func(i, j int) bool {
		aName := utils.TruncateString(res[i].Name, 1)
		bName := utils.TruncateString(res[j].Name, 1)
		aPinyin := pinyin.Pinyin(aName, pinyinArgs)
		bPinyin := pinyin.Pinyin(bName, pinyinArgs)
		aCompareString := aName
		bCompareString := bName
		if len(aPinyin) != 0 {
			aCompareString = aPinyin[0][0]
		}
		if len(bPinyin) != 0 {
			bCompareString = bPinyin[0][0]
		}
		if strings.ToLower(aCompareString) < strings.ToLower(bCompareString) {
			return true
		}
		return false
	})

	return &Response{
		ProjectList: res,
	}, nil
}
