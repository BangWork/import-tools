package file

import (
	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/models/ones"
	"github.com/bangwork/import-tools/serve/utils"
)

func CheckProjectDiskAvailable(url, orgUUID, cookie string, h map[string]string, projectOriginalKeys []string) (bool, error) {
	importCache := common.ImportCacheMap.Get(cookie)
	var attachmentSizes int64
	for _, originalKey := range projectOriginalKeys {
		attachmentSize, err := utils.GetDirSize(common.GenProjectAttachmentFilePath(importCache.LocalHome, originalKey))
		if err != nil {
			return false, err
		}
		attachmentSizes += attachmentSize
	}

	orgConfig, err := ones.GetOrgConfig(url, orgUUID, h)
	if err != nil {
		return false, err
	}
	return common.CheckDiskAvailable(attachmentSizes, orgConfig.FileDiskCapacity), nil
}

func GetONESDiskInfo(url, orgUUID, cookie string, h map[string]string) (resp *ones.ResolveResultResponse, err error) {
	importCache := common.ImportCacheMap.Get(cookie)
	if importCache.ResolveResult != nil && importCache.ResolveResult.JiraVersion == "" {
		importCache.ResolveResult.JiraVersion = "Cloud"
	}
	orgConfig, err := ones.GetOrgConfig(url, orgUUID, h)
	if err != nil {
		return
	}
	diskCapacity := orgConfig.FileDiskCapacity
	fileStorage := orgConfig.FileStorage

	attachmentSize, err := utils.GetDirSize(common.GenAttachmentFilePath(importCache.LocalHome))
	if err != nil {
		return
	}
	available := common.CheckDiskAvailable(attachmentSize, orgConfig.FileDiskCapacity)

	resp = new(ones.ResolveResultResponse)
	resp.ResolveResult = importCache.ResolveResult
	resp.ONESDiskAvailable = available
	resp.ONESFileStorage = fileStorage
	resp.ONESDiskCapacity = diskCapacity
	return
}
