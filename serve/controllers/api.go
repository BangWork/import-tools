package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bangwork/import-tools/serve/utils/timestamp"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/services"
	"github.com/bangwork/import-tools/serve/services/account"
	"github.com/bangwork/import-tools/serve/services/cache"
	"github.com/bangwork/import-tools/serve/services/importer"
	"github.com/bangwork/import-tools/serve/services/issue_type"
	"github.com/bangwork/import-tools/serve/services/log"
	"github.com/bangwork/import-tools/serve/services/project"
	"github.com/bangwork/import-tools/serve/utils"
	"github.com/gin-gonic/gin"
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
	RenderJSON(c, nil, nil)
}

func ResolveProgress(c *gin.Context) {
	info, err := cache.GetCacheInfo()
	if err != nil {
		RenderJSON(c, err, nil)
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
	res, err := cache.GetCacheInfo()
	if err != nil {
		RenderJSON(c, err, res)
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
	list, err := project.GetProjectList()
	RenderJSON(c, err, list)
}

func ChooseTeam(c *gin.Context) {
	var data ChooseTeamRequest
	if err := c.BindJSON(&data); err != nil {
		return
	}
	res, err := cache.GetCacheInfo()
	if err != nil {
		RenderJSON(c, err, res)
		return
	}
	res.ImportTeamUUID = data.TeamUUID
	res.ImportTeamName = data.TeamName
	err = cache.SetCacheInfo(res)
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
	cacheInfo, err := cache.GetCacheInfo()
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
	list, err := issue_type.GetIssueTypeList(typeList, issueTypes)
	RenderJSON(c, err, list)
}

func SetShardDisk(c *gin.Context) {
	var data SetShareDiskRequest
	if err := c.BindJSON(&data); err != nil {
		return
	}
	if data.UseShareDisk {
		exist := utils.CheckPathExist(data.Path)
		if !exist {
			RenderJSON(c, common.Errors(common.NotFoundError, nil), nil)
			return
		}
	}
	res, err := cache.GetCacheInfo()
	if err != nil {
		RenderJSON(c, err, res)
		return
	}
	res.ShareDiskPath = data.Path
	res.UseShareDisk = data.UseShareDisk
	err = cache.SetCacheInfo(res)
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
		importer.StartImport(data.ProjectIDs, data.IssueTypeMap, data.Password)
	}()
	RenderJSON(c, nil, nil)
}

func Reset(c *gin.Context) {
	info, err := cache.GetCacheInfo()
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}

	info.ImportResult.Status = 0
	info.ResolveStatus = 0
	if err = cache.SetCacheInfo(info); err != nil {
		RenderJSON(c, err, nil)
		return
	}
	RenderJSON(c, nil, nil)
}

func PauseImport(c *gin.Context) {
	services.PauseImportSignal = true
	services.StopImportSignal = false

	info, err := cache.GetCacheInfo()
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}

	info.ImportResult.Status = common.ImportStatusPause
	if err = cache.SetCacheInfo(info); err != nil {
		RenderJSON(c, err, nil)
		return
	}

	RenderJSON(c, nil, nil)
}

func ContinueImport(c *gin.Context) {
	services.PauseImportSignal = false
	services.StopImportSignal = false
	info, err := cache.GetCacheInfo()
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}

	info.ImportResult.Status = common.ImportStatusInProgress
	if err = cache.SetCacheInfo(info); err != nil {
		RenderJSON(c, err, nil)
		return
	}
	RenderJSON(c, nil, nil)
}

func StopImport(c *gin.Context) {
	services.StopImportSignal = true
	services.PauseImportSignal = false
	info, err := cache.GetCacheInfo()
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}
	info.ImportResult.Status = common.ImportStatusCancel
	if err = cache.SetCacheInfo(info); err != nil {
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
	info, err := cache.GetCacheInfo()
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
	info, err := cache.GetCacheInfo()
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}
	res := log.GetLogTextByImportUUID(info.ImportTeamUUID, info.ImportUUID)
	RenderJSON(c, nil, res)
}

func GetImportLog(c *gin.Context) {
	startLine := c.Param("start_line")
	info, err := cache.GetCacheInfo()
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
	info, err := cache.GetCacheInfo()
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
	info, err := cache.GetCacheInfo()
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
	info, err := cache.GetCacheInfo()
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}
	RenderJSON(c, nil, info.ImportScope)
}
