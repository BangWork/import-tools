package cache

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/services/importer/types"
	"github.com/bangwork/import-tools/serve/utils"
)

var (
	SharedDiskPath  string
	CurrentCacheKey string
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

	DaysPerWeek         string              `json:"days_per_week"`
	HoursPerDay         string              `json:"hours_per_day"`
	ProjectIssueTypeMap map[string][]string `json:"project_issue_type_map"`

	IssueTypeMap []types.BuiltinIssueTypeMap `json:"issue_type_map"`
	ProjectIDs   []string                    `json:"project_ids"`
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
	filePath := fmt.Sprintf("%s/%s", common.GetCachePath(), cacheFile)
	if utils.CheckPathExist(filePath) {
		return nil
	}
	if err := os.MkdirAll(common.GetCachePath(), 0755); err != nil {
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write([]byte("{}"))
	return nil
}

func GenCacheKey(addr string) string {
	buf := &bytes.Buffer{}
	buf.WriteString(addr)
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func GetCacheInfo(key string) (*Cache, error) {
	if key == "" {
		key = CurrentCacheKey
	}
	if key == "" {
		return nil, nil
	}
	filePath := fmt.Sprintf("%s/%s", common.GetCachePath(), cacheFile)
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("open file error: %s, %s", filePath, err)
		return nil, common.Errors(common.ServerError, nil)
	}

	d := map[string]string{}
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, common.Errors(common.ServerError, nil)
	}
	s, ok := d[key]
	if !ok {
		return new(Cache), nil
	}

	c := new(Cache)
	err = json.Unmarshal([]byte(s), &c)
	if err != nil {
		return nil, common.Errors(common.ServerError, nil)
	}
	return c, nil
}

var cacheLock sync.Mutex

func SetCacheInfo(key string, cache *Cache) error {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	if key == "" {
		key = CurrentCacheKey
	}

	filePath := fmt.Sprintf("%s/%s", common.GetCachePath(), cacheFile)
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("open file error: %s, %s", filePath, err)
		return common.Errors(common.ServerError, nil)
	}

	m := map[string]string{}
	if err := json.Unmarshal(b, &m); err != nil {
		return common.Errors(common.ServerError, nil)
	}

	cb, err := json.Marshal(cache)
	if err != nil {
		return common.Errors(common.ServerError, nil)
	}
	m[key] = string(cb)

	b, err = json.Marshal(m)
	if err != nil {
		return common.Errors(common.ServerError, nil)
	}

	err = ioutil.WriteFile(filePath, b, 0644)
	if err != nil {
		log.Printf("write file error: %s, %s", filePath, err)
		return common.Errors(common.ServerError, nil)
	}
	return nil
}

func SetExpectTimeCache(key string) {
	info, err := GetCacheInfo(key)
	if err != nil {
		log.Println("get cache err", err)
		return
	}
	if info == nil {
		log.Println("err: cache not found")
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
