package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/services"
	"github.com/bangwork/import-tools/serve/services/account"
	"github.com/bangwork/import-tools/serve/services/cache"
	"github.com/bangwork/import-tools/serve/services/importer"
	"github.com/bangwork/import-tools/serve/services/importer/types"
	"github.com/bangwork/import-tools/serve/services/issue_type"
	"github.com/bangwork/import-tools/serve/services/log"
	"github.com/bangwork/import-tools/serve/services/project"
	"github.com/bangwork/import-tools/serve/utils"
	"github.com/bangwork/import-tools/serve/utils/timestamp"
)

func CheckPathExist(c *gin.Context) {
	var data path
	if err := c.BindJSON(&data); err != nil {
		return
	}
	exist := utils.CheckPathExist(common.GenPrivatePath(data.Path))
	if exist {
		RenderJSON(c, nil, nil)
	} else {
		RenderJSON(c, common.Errors(common.NotFoundError, nil), nil)
	}
}

func CheckJiraPathExist(c *gin.Context) {
	var data path
	if err := c.BindJSON(&data); err != nil {
		return
	}
	exist := utils.CheckPathExist(common.GenExportPath(data.Path))
	if exist {
		RenderJSON(c, nil, nil)
	} else {
		RenderJSON(c, common.Errors(common.NotFoundError, nil), nil)
	}
}

func JiraBackUpList(c *gin.Context) {
	var data path
	if err := c.BindJSON(&data); err != nil {
		return
	}
	exist := utils.CheckPathExist(data.Path)
	if !exist {
		RenderJSON(c, common.Errors(common.NotFoundError, nil), nil)
		return
	}
	res, err := utils.ListZipFile(common.GenExportPath(data.Path))
	RenderJSON(c, err, res)
}

func StartResolve(c *gin.Context) {
	var data *account.Account
	if err := c.BindJSON(&data); err != nil {
		return
	}
	if err := data.CheckONESAccount(); err != nil {
		RenderJSON(c, err, nil)
		return
	}
	if err := data.SetCache(); err != nil {
		RenderJSON(c, err, nil)
		return
	}
	go importer.StartResolve(data)

	cache.SaveCacheKey(data.Key())
	RenderJSON(c, nil, map[string]interface{}{
		"key": data.Key(),
	})
}

func ResolveProgress(c *gin.Context) {
	key := c.Param("key")
	info, err := cache.GetCacheInfo(key)
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}
	if info == nil {
		RenderJSON(c, nil, new(ResolveProgressResponse))
		return
	}
	if len(info.BackupName) != 0 {
		exist := utils.CheckPathExist(common.GenBackupFilePath(info.LocalHome, info.BackupName))
		if !exist {
			res := map[string]string{
				"backup_name": info.BackupName,
			}
			RenderJSON(c, common.Errors(common.NotFoundError, res), nil)
			return
		}
	}
	resp := new(ResolveProgressResponse)
	resp.TeamName = info.TeamName
	resp.OrgName = info.OrgName
	resp.MultiTeam = info.MultiTeam
	resp.StartTime = info.ResolveStartTime
	resp.ExpectedTime = info.ExpectedResolveTime
	resp.SpentTime = time.Now().Unix() - info.ResolveStartTime
	resp.Status = info.ResolveStatus
	resp.BackupName = info.BackupName
	if resp.ExpectedTime <= resp.SpentTime {
		resp.ExpectedTime = resp.SpentTime + (resp.SpentTime / 10)
	}
	if resp.Status == common.ResolveStatusDone {
		resp.SpentTime = info.ResolveDoneTime - info.ResolveStartTime
		resp.ExpectedTime = resp.SpentTime
	}
	RenderJSON(c, nil, resp)
}

func StopResolve(c *gin.Context) {
	err := importer.StopResolve()
	RenderJSON(c, err, nil)
}

func ResolveResult(c *gin.Context) {
	key := c.Param("key")
	res, err := cache.GetCacheInfo(key)
	if err != nil {
		RenderJSON(c, err, res)
		return
	}
	if res == nil {
		RenderJSON(c, nil, map[string]interface{}{
			"code": 404,
			"msg":  "invalid key",
		})
		return
	}

	accountInfo := new(account.Account)
	accountInfo.SetCurrentCache(res)
	history, err := accountInfo.GetImportHistory()
	if err != nil {
		RenderJSON(c, err, res)
		return
	}
	if res.ResolveResult.JiraVersion == "" {
		res.ResolveResult.JiraVersion = "Cloud"
	}
	resp := new(account.ResolveResultResponse)
	resp.ResolveResult = res.ResolveResult
	resp.ResolveResult.DiskSetPopUps = res.FileStorage == common.FileStorageLocal
	resp.ImportHistory = history
	RenderJSON(c, err, resp)
}

func ProjectList(c *gin.Context) {
	key := c.Param("key")
	list, err := project.GetProjectList(key)
	RenderJSON(c, err, list)
}

type SaveProjectListReq struct {
	Key        string   `json:"key"`
	ProjectIDs []string `json:"project_ids"`
}

func SaveProjectList(c *gin.Context) {
	req := new(SaveProjectListReq)
	if err := c.BindJSON(&req); err != nil {
		return
	}

	cacheInfo, err := cache.GetCacheInfo(req.Key)
	if err != nil {
		return
	}
	cacheInfo.ProjectIDs = req.ProjectIDs
	err = cache.SetCacheInfo(req.Key, cacheInfo)
	RenderJSON(c, err, nil)
}

func ChooseTeam(c *gin.Context) {
	var data ChooseTeamRequest
	if err := c.BindJSON(&data); err != nil {
		return
	}
	res, err := cache.GetCacheInfo(data.Key)
	if err != nil {
		RenderJSON(c, err, res)
		return
	}
	res.ImportTeamUUID = data.TeamUUID
	res.ImportTeamName = data.TeamName
	err = cache.SetCacheInfo(data.Key, res)
	if err = log.InitTeamLogDir(data.TeamUUID); err != nil {
		RenderJSON(c, err, nil)
		return
	}
	RenderJSON(c, err, nil)
}

func IssueTypeList(c *gin.Context) {
	req := new(GetIssueTypeRequest)
	if err := c.BindJSON(&req); err != nil {
		return
	}
	acc := new(account.Account)
	cacheInfo, err := cache.GetCacheInfo(req.Key)
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}
	acc.SetCurrentCache(cacheInfo)
	typeList, err := acc.GetIssueTypeList()
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}

	issueTypes := make(map[string]bool)
	for _, pid := range req.ProjectIDs {
		for _, issueTypeID := range cacheInfo.ProjectIssueTypeMap[pid] {
			issueTypes[issueTypeID] = true
		}
	}
	list, err := issue_type.GetIssueTypeList(req.Key, typeList, issueTypes)
	RenderJSON(c, err, map[string]interface{}{
		"issue_types":    list,
		"issue_type_map": cacheInfo.IssueTypeMap,
	})
}

type SaveIssueTypeListReq struct {
	Key          string                      `json:"key"`
	IssueTypeMap []types.BuiltinIssueTypeMap `json:"issue_type_map"`
}

func SaveIssueTypeList(c *gin.Context) {
	req := new(SaveIssueTypeListReq)
	if err := c.BindJSON(&req); err != nil {
		return
	}

	cacheInfo, err := cache.GetCacheInfo(req.Key)
	if err != nil {
		return
	}
	cacheInfo.IssueTypeMap = req.IssueTypeMap
	err = cache.SetCacheInfo(req.Key, cacheInfo)
	RenderJSON(c, err, nil)
}

func StartImport(c *gin.Context) {
	var data StartImportRequest
	if err := c.BindJSON(&data); err != nil {
		return
	}
	services.PauseImportSignal = false
	services.StopImportSignal = false
	services.ImportTimeSecond = 0

	go func() {
		for {
			time.Sleep(time.Second)
			if services.CheckIsStop() {
				break
			}
			services.ImportTimeSecond++
		}
	}()

	go func() {
		importer.StartImport(data.Key, data.ProjectIDs, data.IssueTypeMap, data.Password)
	}()
	RenderJSON(c, nil, nil)
}

func Reset(c *gin.Context) {
	key := c.Param("key")
	info, err := cache.GetCacheInfo(key)
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}

	info.ImportResult.Status = 0
	info.ResolveStatus = 0
	if err = cache.SetCacheInfo(key, info); err != nil {
		RenderJSON(c, err, nil)
		return
	}
	RenderJSON(c, nil, nil)
}

func PauseImport(c *gin.Context) {
	services.PauseImportSignal = true
	services.StopImportSignal = false

	key := c.Param("key")
	info, err := cache.GetCacheInfo(key)
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}

	info.ImportResult.Status = common.ImportStatusPause
	if err = cache.SetCacheInfo(key, info); err != nil {
		RenderJSON(c, err, nil)
		return
	}

	RenderJSON(c, nil, nil)
}

func ContinueImport(c *gin.Context) {
	services.PauseImportSignal = false
	services.StopImportSignal = false

	key := c.Param("key")
	info, err := cache.GetCacheInfo(key)
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}

	info.ImportResult.Status = common.ImportStatusInProgress
	if err = cache.SetCacheInfo(key, info); err != nil {
		RenderJSON(c, err, nil)
		return
	}
	RenderJSON(c, nil, nil)
}

func StopImport(c *gin.Context) {
	services.StopImportSignal = true
	services.PauseImportSignal = false
	key := c.Param("key")
	info, err := cache.GetCacheInfo(key)
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}
	info.ImportResult.Status = common.ImportStatusCancel
	if err = cache.SetCacheInfo(key, info); err != nil {
		RenderJSON(c, err, nil)
		return
	}
	accountInfo := new(account.Account)
	if err = accountInfo.Login(); err != nil {
		RenderJSON(c, err, nil)
		return
	}
	RenderJSON(c, accountInfo.InterruptImport(), nil)
}

func ImportProgress(c *gin.Context) {
	key := c.Param("key")
	info, err := cache.GetCacheInfo(key)
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}
	resp := new(cache.ImportResult)
	if info.LocalHome == "" {
		RenderJSON(c, err, resp)
		return
	}
	resp.BackupTime, _ = utils.GetFileModTime(common.GenBackupFilePath(info.LocalHome, info.BackupName))
	resp.TeamName = info.ImportTeamName
	resp.BackupName = "Jira " + info.ResolveResult.JiraVersion
	if info.ImportResult == nil {
		RenderJSON(c, nil, resp)
		return
	}
	resp.StartTime = info.ImportResult.StartTime
	resp.ExpectedTime = info.ImportResult.ExpectedTime
	resp.SpentTime = services.ImportTimeSecond
	resp.Status = info.ImportResult.Status
	if resp.ExpectedTime <= resp.SpentTime {
		add := resp.SpentTime / 10
		if add == 0 {
			add = 1
		}
		resp.ExpectedTime = resp.SpentTime + add
	}
	if resp.Status == common.ImportStatusDone {
		resp.ExpectedTime = resp.SpentTime
	}
	if resp.Status == common.ImportStatusFail {
		RenderJSON(c, common.Errors(common.ServerError, nil), nil)
		return
	}
	RenderJSON(c, nil, resp)
}

func GetAllImportLog(c *gin.Context) {
	key := c.Param("key")
	info, err := cache.GetCacheInfo(key)
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}
	res := log.GetLogTextByImportUUID(info.ImportTeamUUID, info.ImportUUID)
	RenderJSON(c, nil, res)
}

func GetImportLog(c *gin.Context) {
	startLine := c.Param("start_line")
	key := c.Param("key")
	info, err := cache.GetCacheInfo(key)
	if err != nil {
		RenderJSON(c, err, startLine)
		return
	}
	start, err := strconv.Atoi(startLine)
	if err != nil {
		RenderJSON(c, err, startLine)
		return
	}
	res := log.GetLogText(info.ImportTeamUUID, info.ImportUUID, start)
	RenderJSON(c, nil, res)
}

func DownloadLogFile(c *gin.Context) {
	key := c.Param("key")
	info, err := cache.GetCacheInfo(key)
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}
	if info.ImportTeamUUID == "" {
		RenderJSON(c, err, nil)
		return
	}
	allLog, err := log.DownloadAllLog(info.ImportTeamUUID)
	if err != nil {
		return
	}
	fileContentDisposition := fmt.Sprintf("attachment;filename=%s-import-log-%s.txt", timestamp.GetDateString(), info.ImportTeamUUID)
	c.Header("Content-Type", "application/txt")
	c.Header("Content-Disposition", fileContentDisposition)
	c.Data(http.StatusOK, "", allLog)
}

func DownloadCurrentLogFile(c *gin.Context) {
	key := c.Param("key")
	info, err := cache.GetCacheInfo(key)
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}
	if info.ImportTeamUUID == "" {
		RenderJSON(c, err, nil)
		return
	}
	allLog, err := log.DownloadAllLog(info.ImportTeamUUID)
	if err != nil {
		return
	}
	fileContentDisposition := fmt.Sprintf("attachment;filename=%s-import-log-%s.txt", timestamp.GetDateString(), info.ImportTeamUUID)
	c.Header("Content-Type", "application/txt")
	c.Header("Content-Disposition", fileContentDisposition)
	c.Data(http.StatusOK, "", allLog)
}

func GetScope(c *gin.Context) {
	key := c.Param("key")
	info, err := cache.GetCacheInfo(key)
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}
	RenderJSON(c, nil, info.ImportScope)
}
