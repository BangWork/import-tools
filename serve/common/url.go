package common

import "fmt"

const (
	//projectApiPrefix = "/"
	projectApiPrefix = "/project/api/project/"
)

func GenApiUrl(url, uri string) string {
	return fmt.Sprintf("%s%s%s", url, projectApiPrefix, uri)
}

func GenBackupFilePath(localHome, backupName string) string {
	return fmt.Sprintf("%s/%s/%s", localHome, JiraExportDir, backupName)
}

func GenExportPath(localHome string) string {
	return fmt.Sprintf("%s/%s", localHome, JiraExportDir)
}

func GenAttachmentFilePath(localHome string) string {
	return fmt.Sprintf("%s/data/attachments", localHome)
}

func GenPrivatePath(shareDiskPath string) string {
	return fmt.Sprintf("%s/%s", shareDiskPath, ShareDiskPathPrivate)
}
