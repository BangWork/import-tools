package project

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/bangwork/import-tools/serve/utils"

	pinyin "github.com/mozillazg/go-pinyin"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/services/cache"
	"github.com/bangwork/import-tools/serve/services/importer/resolve"
	"github.com/bangwork/import-tools/serve/utils/xml"
)

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProjectRes struct {
	ProjectList []*Project  `json:"projects"`
	Cache       interface{} `json:"cache"`
}

func GetProjectList(key string) (*ProjectRes, error) {
	list, err := cache.GetCacheInfo(key)
	if err != nil {
		return nil, err
	}
	v, found := list.MapFilePath[common.TagProject]
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

	return &ProjectRes{
		ProjectList: res,
		Cache:       list.ProjectIDs[list.TeamUUID],
	}, nil
}
