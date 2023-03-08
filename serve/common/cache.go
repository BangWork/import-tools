package common

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/bangwork/import-tools/serve/services/importer/types"
)

var (
	SharedDiskPath  string
	CurrentCacheKey string
)

type Cache struct {
	URL                 string
	UserName            string
	Password            string
	ImportBatchTaskUUID string
	ImportTimeSecond    int64
	StopResolveSignal   bool
	PauseImportSignal   bool
	StopImportSignal    bool
	ImportUUID          string `json:"import_uuid"`
	ResolveStartTime    int64  `json:"start_resolve_time"`
	ResolveStatus       int    `json:"resolve_status"`
	ExpectedResolveTime int64  `json:"expected_resolve_time"`
	ResolveDoneTime     int64  `json:"resolve_done_time"`
	FileStorage         string `json:"file_storage"`
	MultiTeam           bool   `json:"multi_team"`
	OrgName             string `json:"org_name"`
	OrgUUID             string `json:"org_uuid"`
	TeamUUID            string `json:"team_uuid"`
	TeamName            string `json:"team_name"`
	ImportUserUUID      string `json:"import_user_uuid"`
	LocalHome           string `json:"local_home"`
	BackupName          string
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
	ProjectAssignMap    map[string]string
	ProjectCategoryMap  map[string]string

	IssueTypeMap []types.BuiltinIssueTypeMap `json:"issue_type_map"`
	ProjectIDs   []string                    `json:"project_ids"`
}

func (c *Cache) CheckIsStop() (stop bool) {
	for c.PauseImportSignal {
		time.Sleep(5 * time.Second)
	}
	return c.StopImportSignal
}

type ResolveResult struct {
	JiraVersion     string `json:"jira_version"`
	ProjectCount    int64  `json:"project_count"`
	IssueCount      int64  `json:"issue_count"`
	MemberCount     int64  `json:"member_count"`
	AttachmentCount int64  `json:"attachment_count"`
	AttachmentSize  int64  `json:"attachment_size"`
	JiraServerID    string `json:"jira_server_id"`
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
	filePath := fmt.Sprintf("%s/%s", GetCachePath(), cacheFile)
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("open file error: %s, %s", filePath, err)
		return nil, Errors(ServerError, nil)
	}

	d := map[string]string{}
	if err := json.Unmarshal(b, &d); err != nil {
		return nil, Errors(ServerError, nil)
	}
	s, ok := d[key]
	if !ok {
		return new(Cache), nil
	}

	c := new(Cache)
	err = json.Unmarshal([]byte(s), &c)
	if err != nil {
		return nil, Errors(ServerError, nil)
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

	filePath := fmt.Sprintf("%s/%s", GetCachePath(), cacheFile)
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("open file error: %s, %s", filePath, err)
		return Errors(ServerError, nil)
	}

	m := map[string]string{}
	if err := json.Unmarshal(b, &m); err != nil {
		return Errors(ServerError, nil)
	}

	cb, err := json.Marshal(cache)
	if err != nil {
		return Errors(ServerError, nil)
	}
	m[key] = string(cb)

	b, err = json.Marshal(m)
	if err != nil {
		return Errors(ServerError, nil)
	}

	err = ioutil.WriteFile(filePath, b, 0644)
	if err != nil {
		log.Printf("write file error: %s, %s", filePath, err)
		return Errors(ServerError, nil)
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
	if info.ImportScope.IssueCount != Calculating {
		info.ImportResult.ExpectedTime = getExpectedTime(info.ImportScope.IssueCount, issueCountExpectTimeMap)
		return
	}
	if info.ImportScope.AttachmentSize != Calculating {
		info.ImportResult.ExpectedTime = getExpectedTime(info.ImportScope.AttachmentSize, attachmentSizeExpectTimeMap)
		return
	}
	if info.ImportScope.ProjectCount != Calculating {
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
