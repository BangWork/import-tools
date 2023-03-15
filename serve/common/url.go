package common

import (
	"fmt"
	"path"
)

const (
	//projectApiPrefix = "/"
	projectApiPrefix = "/project/api/project/"
)

func GenApiUrl(url, uri string) string {
	return fmt.Sprintf("%s%s%s", url, projectApiPrefix, uri)
}

func GenBackupFilePath(localHome, backupName string) string {
	return path.Join(localHome, JiraExportDir, backupName)
}

func GenExportPath(localHome string) string {
	return path.Join(localHome, JiraExportDir)
}

func GenAttachmentFilePath(localHome string) string {
	return path.Join(localHome, "/data/attachments")
}

func GenPrivatePath(shareDiskPath string) string {
	return path.Join(shareDiskPath, ShareDiskPathPrivate)
}
