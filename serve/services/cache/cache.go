package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils"
)

type Cache struct {
	ImportUUID          string            `json:"import_uuid"`
	URL                 string            `json:"url"`
	Email               string            `json:"email"`
	Password            string            `json:"password"`
	ResolveStartTime    int64             `json:"start_resolve_time"`
	ResolveStatus       int               `json:"resolve_status"`
	ExpectedResolveTime int64             `json:"expected_resolve_time"`
	ResolveDoneTime     int64             `json:"resolve_done_time"`
	FileStorage         string            `json:"file_storage"`
	MultiTeam           bool              `json:"multi_team"`
	OrgName             string            `json:"org_name"`
	OrgUUID             string            `json:"org_uuid"`
	TeamUUID            string            `json:"team_uuid"`
	TeamName            string            `json:"team_name"`
	ImportUserUUID      string            `json:"import_user_uuid"`
	LocalHome           string            `json:"local_home"`
	BackupName          string            `json:"backup_name"`
	ResolveResult       *ResolveResult    `json:"resolve_result"`
	ImportResult        *ImportResult     `json:"import_result"`
	ImportScope         *ResolveResult    `json:"import_scope"`
	MapFilePath         map[string]string `json:"map_file_path"`
	MapOutputFilePath   map[string]string `json:"map_output_file_path"`
	ImportTeamUUID      string            `json:"import_team_uuid"`
	ImportTeamName      string            `json:"import_team_name"`
	UseShareDisk        bool              `json:"use_share_disk"`
	ShareDiskPath       string            `json:"share_disk_path"`

	DaysPerWeek         string              `json:"days_per_week"`
	HoursPerDay         string              `json:"hours_per_day"`
	ProjectIssueTypeMap map[string][]string `json:"project_issue_type_map"`
}

type ResolveResult struct {
	JiraVersion     string `json:"jira_version"`
	ProjectCount    int64  `json:"project_count"`
	IssueCount      int64  `json:"issue_count"`
	MemberCount     int64  `json:"member_count"`
	AttachmentCount int64  `json:"attachment_count"`
	AttachmentSize  int64  `json:"attachment_size"`
	JiraServerID    string `json:"jira_server_id"`

	DiskSetPopUps bool `json:"disk_set_pop_ups"`
}

type ImportResult struct {
	TeamName      string `json:"team_name"`
	StartTime     int64  `json:"start_time"`
	BackupTime    int64  `json:"backup_time"`
	Status        int    `json:"status"`
	ExpectedTime  int64  `json:"expected_time"`
	DoneTime      int64  `json:"done_time"`
	LastPauseTime int64  `json:"last_pause_time"`
	SpentTime     int64  `json:"spent_time"`
	BackupName    string `json:"backup_name"`
}

const cacheFile = "import.json"

var (
	attachmentSizeExpectTimeMap = map[int64]int64{
		100000000000: 624000,
		50000000000:  324000,
		5000000000:   32400,
		450000000:    15400,
		0:            60,
	}

	issueCountExpectTimeMap = map[int64]int64{
		1000000: 624000,
		500000:  324000,
		50000:   32400,
		4500:    15400,
		2000:    7400,
		1000:    3300,
		100:     300,
		0:       60,
	}

	projectCountExpectTimeMap = map[int64]int64{
		9000: 324000,
		900:  32400,
		450:  15400,
		200:  7400,
		100:  3300,
		10:   300,
		0:    60,
	}
)

func InitCacheFile() error {
	filePath := fmt.Sprintf("%s/%s", common.Path, cacheFile)
	if utils.CheckPathExist(filePath) {
		return nil
	}
	if err := os.MkdirAll(common.Path, 0755); err != nil {
		return err
	}
	_, err := os.Create(filePath)
	if err != nil {
		return err
	}
	cache := new(Cache)
	if err := SetCacheInfo(cache); err != nil {
		return err
	}
	return nil
}

func GetCacheInfo() (*Cache, error) {
	filePath := fmt.Sprintf("%s/%s", common.Path, cacheFile)
	if !utils.CheckPathExist(filePath) {
		log.Fatalf("cache file missing: %s", filePath)
		return nil, common.Errors(common.CacheFileNotFoundError, nil)
	}
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("open file error: %s, %s", filePath, err)
		return nil, common.Errors(common.ServerError, nil)
	}
	res := new(Cache)
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		log.Fatalf("parse json file error: %s, %s", filePath, err)
		return nil, common.Errors(common.ServerError, nil)
	}
	return res, nil
}

func SetCacheInfo(cache *Cache) error {
	filePath := fmt.Sprintf("%s/%s", common.Path, cacheFile)
	if !utils.CheckPathExist(filePath) {
		log.Fatalf("cache file missing: %s", filePath)
		return common.Errors(common.CacheFileNotFoundError, nil)
	}
	bytes, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		log.Fatalf("parse json file error: %s, %s", filePath, err)
		return common.Errors(common.ServerError, nil)
	}
	err = ioutil.WriteFile(filePath, bytes, 0644)
	if err != nil {
		log.Fatalf("write file error: %s, %s", filePath, err)
		return common.Errors(common.ServerError, nil)
	}
	return nil
}

func SetExpectTimeCache() {
	info, err := GetCacheInfo()
	if err != nil {
		log.Println("get cache err", err)
		return
	}
	if info.ImportScope == nil {
		return
	}
	if info.ImportScope.IssueCount != common.Calculating {
		info.ImportResult.ExpectedTime = getExpectedTime(info.ImportScope.IssueCount, issueCountExpectTimeMap)
		return
	}
	if info.ImportScope.AttachmentSize != common.Calculating {
		info.ImportResult.ExpectedTime = getExpectedTime(info.ImportScope.AttachmentSize, attachmentSizeExpectTimeMap)
		return
	}
	if info.ImportScope.ProjectCount != common.Calculating {
		info.ImportResult.ExpectedTime = getExpectedTime(info.ImportScope.ProjectCount, projectCountExpectTimeMap)
		return
	}
	return
}

func getExpectedTime(count int64, config map[int64]int64) int64 {
	for size, t := range config {
		if count > size {
			return t
		}
	}
	return 100
}
