package importer

import (
	"log"
	"runtime/debug"
	"time"

	"github.com/bangwork/import-tools/serve/services/account"

	"github.com/bangwork/import-tools/serve/services"

	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/services/cache"
	"github.com/bangwork/import-tools/serve/services/importer/constants"
	"github.com/bangwork/import-tools/serve/services/importer/sync"
	"github.com/bangwork/import-tools/serve/services/importer/types"
)

func StartResolve(account *account.Account) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[importer] err: %s\n%s", p, debug.Stack())
		}
	}()
	services.StopResolveSignal = false
	task := &types.ImportTask{
		UserUUID:          account.AuthHeader[common.UserID],
		ImportType:        constants.ImportTypeJira,
		LocalFilePath:     common.GenBackupFilePath(account.LocalHome, account.BackupName),
		AttachmentsPath:   common.GenAttachmentFilePath(account.LocalHome),
		BuiltinIssueTypes: nil,
	}

	if err := sync.NewImporter(task).Resolve(); err != nil {
		return
	}
}

func StopResolve() error {
	services.StopResolveSignal = true
	return nil
}

func StartImport(projectIDs []string, builtinIssueTypeMap []types.BuiltinIssueTypeMap, password string) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[importer] err: %s\n%s", p, debug.Stack())
			cacheInfo, err := cache.GetCacheInfo()
			if err != nil {
				log.Println("get cache fail", err)
				return
			}
			cacheInfo.ImportResult = &cache.ImportResult{
				Status: common.ImportStatusFail,
			}
			if err = cache.SetCacheInfo(cacheInfo); err != nil {
				return
			}
		}
	}()

	cacheInfo, err := cache.GetCacheInfo()
	if err != nil {
		log.Println("get cache fail", err)
		return
	}
	cacheInfo.ImportResult = &cache.ImportResult{
		StartTime: time.Now().Unix(),
		Status:    common.ImportStatusInProgress,
		TeamName:  cacheInfo.TeamName,
	}
	cacheInfo.ImportScope = &cache.ResolveResult{
		JiraVersion:     cacheInfo.ResolveResult.JiraVersion,
		ProjectCount:    int64(len(projectIDs)),
		IssueCount:      common.Calculating,
		MemberCount:     common.Calculating,
		AttachmentSize:  common.Calculating,
		AttachmentCount: common.Calculating,
	}
	if err := cache.SetCacheInfo(cacheInfo); err != nil {
		log.Println("set cache fail", err)
		return
	}
	cache.SetExpectTimeCache()
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
	}

	if err := sync.NewImporter(task).Import(); err != nil {
		return
	}
}
