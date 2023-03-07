package importer

import (
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/services/importer/constants"
	"github.com/bangwork/import-tools/serve/services/importer/sync"
	"github.com/bangwork/import-tools/serve/services/importer/types"
)

func StartResolve(cookie, userUUID, localHome, backupName string) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[importer] err: %s\n%s", p, debug.Stack())
		}
	}()
	if len(cookie) == 0 || len(userUUID) == 0 || len(localHome) == 0 || len(backupName) == 0 {
		panic(fmt.Sprintf("params err：%s/%s/%s/%s", cookie, userUUID, localHome, backupName))
	}
	cacheInfo := common.ImportCacheMap.Get(cookie)
	cacheInfo.BackupName = backupName
	cacheInfo.LocalHome = localHome
	cacheInfo.StopResolveSignal = false
	cacheInfo.ResolveStartTime = time.Now().Unix()
	common.ImportCacheMap.Set(cookie, cacheInfo)
	task := &types.ImportTask{
		UserUUID:          userUUID,
		ImportType:        constants.ImportTypeJira,
		LocalFilePath:     common.GenBackupFilePath(localHome, backupName),
		AttachmentsPath:   common.GenAttachmentFilePath(localHome),
		BuiltinIssueTypes: nil,
		Cookie:            cookie,
		BackupName:        backupName,
	}

	if err := sync.NewImporter(task).Resolve(); err != nil {
		return
	}
}

func StopResolve(cookie string) error {
	cacheInfo := common.ImportCacheMap.Get(cookie)
	cacheInfo.StopResolveSignal = true
	common.ImportCacheMap.Set(cookie, cacheInfo)
	return nil
}

func StartImport(key string, projectIDs []string, builtinIssueTypeMap []types.BuiltinIssueTypeMap, password string) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[importer] err: %s\n%s", p, debug.Stack())
			cacheInfo, err := common.GetCacheInfo(key)
			if err != nil {
				log.Println("get cache fail", err)
				return
			}
			cacheInfo.ImportResult = &common.ImportResult{
				Status: common.ImportStatusFail,
			}
			if err = common.SetCacheInfo(key, cacheInfo); err != nil {
				return
			}
		}
	}()

	cacheInfo, err := common.GetCacheInfo(key)
	if err != nil {
		log.Println("get cache fail", err)
		return
	}
	cacheInfo.ImportResult = &common.ImportResult{
		StartTime: time.Now().Unix(),
		Status:    common.ImportStatusInProgress,
		TeamName:  cacheInfo.TeamName,
	}
	cacheInfo.ImportScope = &common.ResolveResult{
		JiraVersion:     cacheInfo.ResolveResult.JiraVersion,
		ProjectCount:    int64(len(projectIDs)),
		IssueCount:      common.Calculating,
		MemberCount:     common.Calculating,
		AttachmentSize:  common.Calculating,
		AttachmentCount: common.Calculating,
	}
	if err := common.SetCacheInfo(key, cacheInfo); err != nil {
		log.Println("set cache fail", err)
		return
	}
	common.SetExpectTimeCache(key)
	mapProjectID := map[string]bool{}
	for _, id := range projectIDs {
		mapProjectID[id] = true
	}

	task := &types.ImportTask{
		UserUUID:           cacheInfo.ImportUserUUID,
		ImportType:         constants.ImportTypeJira,
		LocalFilePath:      common.GenBackupFilePath(cacheInfo.LocalHome, cacheInfo.BackupName),
		BuiltinIssueTypes:  builtinIssueTypeMap,
		SelectedProjectIDs: mapProjectID,
		MapFilePath:        cacheInfo.MapFilePath,
		ImportTeamUUID:     cacheInfo.ImportTeamUUID,
		AttachmentsPath:    common.GenAttachmentFilePath(cacheInfo.LocalHome),
		Password:           password,
		Key:                key,
	}

	if err := sync.NewImporter(task).Import(); err != nil {
		return
	}
}
