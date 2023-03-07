package controllers

import (
	"fmt"
	logApi "log"
	"net/http"
	"strconv"
	"time"

	"github.com/bangwork/import-tools/serve/services/config"

	"github.com/bangwork/import-tools/serve/services/team"

	common2 "github.com/bangwork/import-tools/serve/models/common"

	"github.com/bangwork/import-tools/serve/services/auth"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/services"
	"github.com/bangwork/import-tools/serve/services/account"
	"github.com/bangwork/import-tools/serve/services/importer"
	"github.com/bangwork/import-tools/serve/services/importer/types"
	"github.com/bangwork/import-tools/serve/services/log"
	"github.com/bangwork/import-tools/serve/services/project"
	"github.com/bangwork/import-tools/serve/utils"
	"github.com/bangwork/import-tools/serve/utils/timestamp"
	"github.com/gin-gonic/gin"
)

func Config(c *gin.Context) {
	RenderJSON(c, nil, common2.GetConfig())
}

func Login(c *gin.Context) {
	var data *services.LoginRequest
	if err := c.BindJSON(&data); err != nil {
		return
	}
	cookie, err := auth.Login(data)
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}
	common.ImportCacheMap.Set(cookie, new(common.Cache))
	c.SetCookie(common.LoginCookieName, cookie, common.GetCookieExpireTime(), "/", "", false, true)
	RenderJSON(c, err, nil)
}

func Logout(c *gin.Context) {
	cookie, err := c.Cookie(common.LoginCookieName)
	if err != nil {
		logApi.Printf("get cookie fail:%+v", err)
	}
	auth.Logout(cookie)
	RenderJSON(c, nil, nil)
}

func UserJiraConfig(c *gin.Context) {
	jiraConfig, err := config.GetUserJiraConfig(getCookie(c))
	RenderJSON(c, err, jiraConfig)
}

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
	var req *services.StartResolveRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}
	userUUID := c.GetString("userUUID")
	cookie := getCookie(c)
	//onesUrl := c.GetString("onesUrl")
	go importer.StartResolve(cookie, userUUID, req.JiraLocalHome, req.BackupName)

	RenderJSON(c, nil, nil)
	//cache.CurrentCacheKey = cache.GenCacheKey(onesUrl)
	//RenderJSON(c, nil, map[string]interface{}{
	//	"key": cache.CurrentCacheKey,
	//})
}

func ResolveProgress(c *gin.Context) {
	importCache := common.ImportCacheMap.Get(getCookie(c))
	if len(importCache.BackupName) != 0 {
		exist := utils.CheckPathExist(common.GenBackupFilePath(importCache.LocalHome, importCache.BackupName))
		if !exist {
			res := map[string]string{
				"backup_name": importCache.BackupName,
			}
			RenderJSON(c, common.Errors(common.NotFoundError, res), nil)
			return
		}
	}
	resp := new(ResolveProgressResponse)
	resp.StartTime = importCache.ResolveStartTime
	resp.ExpectedTime = importCache.ExpectedResolveTime
	resp.SpentTime = time.Now().Unix() - importCache.ResolveStartTime
	resp.Status = importCache.ResolveStatus
	if resp.ExpectedTime <= resp.SpentTime {
		resp.ExpectedTime = resp.SpentTime + (resp.SpentTime / 10)
	}
	if resp.Status == common.ResolveStatusDone {
		resp.SpentTime = importCache.ResolveDoneTime - importCache.ResolveStartTime
		resp.ExpectedTime = resp.SpentTime
	}
	RenderJSON(c, nil, resp)
}

func StopResolve(c *gin.Context) {
	cookie := getCookie(c)
	err := importer.StopResolve(cookie)
	RenderJSON(c, err, nil)
}

func ResolveResult(c *gin.Context) {
	cookie := getCookie(c)
	importCache := common.ImportCacheMap.Get(cookie)
	if importCache.ResolveResult.JiraVersion == "" {
		importCache.ResolveResult.JiraVersion = "Cloud"
	}
	resp := new(account.ResolveResultResponse)
	resp.ResolveResult = importCache.ResolveResult
	RenderJSON(c, nil, resp)
}

func TeamList(c *gin.Context) {
	h := getONESHeader(c)
	orgUUID := getOrgUUID(c)
	url := getONESUrl(c)
	l, err := team.GetImportHistory(orgUUID, url, h)
	RenderJSON(c, err, l)
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

	cacheInfo, err := common.GetCacheInfo(req.Key)
	if err != nil {
		return
	}
	cacheInfo.ProjectIDs = req.ProjectIDs
	err = common.SetCacheInfo(req.Key, cacheInfo)
	RenderJSON(c, err, nil)
}

func ChooseTeam(c *gin.Context) {
	var data ChooseTeamRequest
	if err := c.BindJSON(&data); err != nil {
		return
	}
	res, err := common.GetCacheInfo(data.Key)
	if err != nil {
		RenderJSON(c, err, res)
		return
	}
	res.ImportTeamUUID = data.TeamUUID
	res.ImportTeamName = data.TeamName
	err = common.SetCacheInfo(data.Key, res)
	if err = log.InitTeamLogDir(data.TeamUUID); err != nil {
		RenderJSON(c, err, nil)
		return
	}
	RenderJSON(c, err, nil)
}

//func IssueTypeList(c *gin.Context) {
//	req := new(GetIssueTypeRequest)
//	if err := c.BindJSON(&req); err != nil {
//		return
//	}
//	acc := new(account.Account)
//	cacheInfo, err := cache.GetCacheInfo(req.Key)
//	if err != nil {
//		RenderJSON(c, err, nil)
//		return
//	}
//	acc.SetCurrentCache(cacheInfo)
//	typeList, err := acc.GetIssueTypeList()
//	if err != nil {
//		RenderJSON(c, err, nil)
//		return
//	}
//
//	issueTypes := make(map[string]bool)
//	for _, pid := range req.ProjectIDs {
//		for _, issueTypeID := range cacheInfo.ProjectIssueTypeMap[pid] {
//			issueTypes[issueTypeID] = true
//		}
//	}
//	list, err := issue_type.GetIssueTypeList(req.Key, typeList, issueTypes)
//	RenderJSON(c, err, map[string]interface{}{
//		"issue_types":    list,
//		"issue_type_map": cacheInfo.IssueTypeMap,
//	})
//}

type SaveIssueTypeListReq struct {
	Key          string                      `json:"key"`
	IssueTypeMap []types.BuiltinIssueTypeMap `json:"issue_type_map"`
}

func SaveIssueTypeList(c *gin.Context) {
	req := new(SaveIssueTypeListReq)
	if err := c.BindJSON(&req); err != nil {
		return
	}

	cacheInfo, err := common.GetCacheInfo(req.Key)
	if err != nil {
		return
	}
	cacheInfo.IssueTypeMap = req.IssueTypeMap
	err = common.SetCacheInfo(req.Key, cacheInfo)
	RenderJSON(c, err, nil)
}

func StartImport(c *gin.Context) {
	var data StartImportRequest
	if err := c.BindJSON(&data); err != nil {
		return
	}
	cookie := getCookie(c)
	importCache := common.ImportCacheMap.Get(cookie)
	importCache.PauseImportSignal = false
	importCache.StopImportSignal = false
	importCache.ImportTimeSecond = 0
	common.ImportCacheMap.Set(cookie, importCache)

	go func() {
		for {
			time.Sleep(time.Second)
			if importCache.CheckIsStop() {
				break
			}
		}
	}()

	go func() {
		importer.StartImport(data.Key, data.ProjectIDs, data.IssueTypeMap, data.Password)
	}()
	RenderJSON(c, nil, nil)
}

func Reset(c *gin.Context) {
	key := c.Param("key")
	info, err := common.GetCacheInfo(key)
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}

	info.ImportResult.Status = 0
	info.ResolveStatus = 0
	if err = common.SetCacheInfo(key, info); err != nil {
		RenderJSON(c, err, nil)
		return
	}
	RenderJSON(c, nil, nil)
}

func PauseImport(c *gin.Context) {
	cookie := getCookie(c)
	importCache := common.ImportCacheMap.Get(cookie)
	importCache.PauseImportSignal = true
	importCache.StopImportSignal = false
	importCache.ImportResult.Status = common.ImportStatusPause
	common.ImportCacheMap.Set(cookie, importCache)
	RenderJSON(c, nil, nil)
}

func ContinueImport(c *gin.Context) {
	cookie := getCookie(c)
	importCache := common.ImportCacheMap.Get(cookie)
	importCache.PauseImportSignal = false
	importCache.StopImportSignal = false
	importCache.ImportResult.Status = common.ImportStatusInProgress
	importCache.ImportResult.Status = common.ImportStatusPause
	common.ImportCacheMap.Set(cookie, importCache)
	RenderJSON(c, nil, nil)
}

//func StopImport(c *gin.Context) {
//	cookie := getCookie(c)
//	importCache := services.ImportCacheMap.Get(cookie)
//	importCache.StopImportSignal = true
//	importCache.PauseImportSignal = false
//	importCache.ImportResult.Status = common.ImportStatusCancel
//	services.ImportCacheMap.Set(cookie, importCache)
//	RenderJSON(c, accountInfo.InterruptImport(), nil)
//}

//func ImportProgress(c *gin.Context) {
//	cookie := getCookie(c)
//	importCache := services.ImportCacheMap.Get(cookie)
//	if importCache.LocalHome == "" {
//		RenderJSON(c, common.Errors(common.NotFoundError, nil), nil)
//		return
//	}
//	resp.BackupTime, _ = utils.GetFileModTime(common.GenBackupFilePath(importCache.LocalHome, importCache.BackupName))
//	resp.TeamName = info.ImportTeamName
//	resp.BackupName = "Jira " + info.ResolveResult.JiraVersion
//	if info.ImportResult == nil {
//		RenderJSON(c, nil, resp)
//		return
//	}
//	resp.StartTime = info.ImportResult.StartTime
//	resp.ExpectedTime = info.ImportResult.ExpectedTime
//	resp.SpentTime = services.ImportTimeSecond
//	resp.Status = info.ImportResult.Status
//	if resp.ExpectedTime <= resp.SpentTime {
//		add := resp.SpentTime / 10
//		if add == 0 {
//			add = 1
//		}
//		resp.ExpectedTime = resp.SpentTime + add
//	}
//	if resp.Status == common.ImportStatusDone {
//		resp.ExpectedTime = resp.SpentTime
//	}
//	if resp.Status == common.ImportStatusFail {
//		RenderJSON(c, common.Errors(common.ServerError, nil), nil)
//		return
//	}
//	RenderJSON(c, nil, resp)
//}

func GetAllImportLog(c *gin.Context) {
	key := c.Param("key")
	info, err := common.GetCacheInfo(key)
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
	info, err := common.GetCacheInfo(key)
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
	info, err := common.GetCacheInfo(key)
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
	info, err := common.GetCacheInfo(key)
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
	info, err := common.GetCacheInfo(key)
	if err != nil {
		RenderJSON(c, err, nil)
		return
	}
	RenderJSON(c, nil, info.ImportScope)
}
