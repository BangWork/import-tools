package controllers

import (
	"fmt"
	logApi "log"
	"net/http"
	"strconv"
	"time"

	"github.com/bangwork/import-tools/serve/models/ones"

	"github.com/bangwork/import-tools/serve/common"
	common2 "github.com/bangwork/import-tools/serve/models/common"
	"github.com/bangwork/import-tools/serve/services"
	"github.com/bangwork/import-tools/serve/services/auth"
	"github.com/bangwork/import-tools/serve/services/config"
	"github.com/bangwork/import-tools/serve/services/file"
	"github.com/bangwork/import-tools/serve/services/importer"
	"github.com/bangwork/import-tools/serve/services/importer/types"
	"github.com/bangwork/import-tools/serve/services/issue_type"
	"github.com/bangwork/import-tools/serve/services/log"
	"github.com/bangwork/import-tools/serve/services/project"
	"github.com/bangwork/import-tools/serve/services/team"
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
	if req.JiraLocalHome == "" {
		req.JiraLocalHome = common.GetJiraLocalHome()
	}
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
	h := getONESHeader(c)
	url := getONESUrl(c)

	info, err := file.GetONESDiskInfo(url, cookie, h)
	RenderJSON(c, err, info)
}

func TeamList(c *gin.Context) {
	h := getONESHeader(c)
	orgUUID := getOrgUUID(c)
	url := getONESUrl(c)
	l, err := team.GetImportHistory(orgUUID, url, h)
	RenderJSON(c, err, l)
}

func ProjectList(c *gin.Context) {
	cookie := getCookie(c)
	list, err := project.GetProjectList(cookie)
	RenderJSON(c, err, list)
}

func ProjectHistoryConfig(c *gin.Context) {
	orgUUID := getOrgUUID(c)
	h := getONESHeader(c)
	url := getONESUrl(c)
	list, err := config.GetHistoryProjectConfig(url, orgUUID, h)
	RenderJSON(c, err, list)
}

func SaveProjectHistoryConfig(c *gin.Context) {
	orgUUID := getOrgUUID(c)
	h := getONESHeader(c)
	url := getONESUrl(c)
	projectIDs := make([]string, 0)
	if err := c.BindJSON(&projectIDs); err != nil {
		return
	}
	err := config.SetHistoryProjectConfig(url, orgUUID, h, projectIDs)
	RenderJSON(c, err, nil)
}

func IssueTypeHistoryConfig(c *gin.Context) {
	orgUUID := getOrgUUID(c)
	h := getONESHeader(c)
	url := getONESUrl(c)
	list, err := config.GetHistoryIssueTypeConfig(url, orgUUID, h)
	RenderJSON(c, err, list)
}

func SaveIssueTypeHistoryConfig(c *gin.Context) {
	orgUUID := getOrgUUID(c)
	h := getONESHeader(c)
	url := getONESUrl(c)
	issueTypeMap := make([]*ones.IssueTypeMapConfig, 0)
	if err := c.BindJSON(&issueTypeMap); err != nil {
		return
	}
	err := config.SetHistoryIssueTypeConfig(url, orgUUID, h, issueTypeMap)
	RenderJSON(c, err, nil)
}

func CheckProjectDisk(c *gin.Context) {
	cookie := getCookie(c)
	h := getONESHeader(c)
	url := getONESUrl(c)

	req := make([]string, 0)
	if err := c.BindJSON(&req); err != nil {
		return
	}
	available, err := file.CheckProjectDiskAvailable(url, cookie, h, req)
	RenderJSON(c, err, map[string]bool{
		"ones_disk_available": available,
	})
}

func UnBoundONESIssueType(c *gin.Context) {
	h := getONESHeader(c)
	url := getONESUrl(c)
	teamUUID := getTeamUUID(c)
	cookie := getCookie(c)

	res, err := issue_type.GetUnBoundIssueTypes(cookie, url, teamUUID, h)
	RenderJSON(c, err, res)
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

func IssueTypeList(c *gin.Context) {
	projectIDs := make([]string, 0)
	if err := c.BindJSON(&projectIDs); err != nil {
		return
	}

	url := getONESUrl(c)
	teamUUID := getTeamUUID(c)
	h := getONESHeader(c)
	cookie := getCookie(c)

	list, err := issue_type.GetIssueTypeList(url, teamUUID, cookie, h, projectIDs)
	RenderJSON(c, err, list)
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
