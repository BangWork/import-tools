package config

import (
	"github.com/bangwork/import-tools/serve/common"
	"github.com/bangwork/import-tools/serve/utils"
)

type UserJiraConfig struct {
	ShowJiraLocalHome    bool     `json:"show_jira_local_home"`
	LastJiraLocalHome    string   `json:"last_jira_local_home"`
	DefaultJiraLocalHome string   `json:"default_jira_local_home"`
	LastBackupName       string   `json:"last_backup_name"`
	BackupList           []string `json:"backup_list"`
}

func GetUserJiraConfig(cookie string) (*UserJiraConfig, error) {
	importCache := common.ImportCacheMap.Get(cookie)
	r := new(UserJiraConfig)
	r.DefaultJiraLocalHome = common.DefaultJiraLocalHome
	r.ShowJiraLocalHome = common.GetJiraLocalHome() == ""
	r.LastJiraLocalHome = common.GetJiraLocalHome()
	r.LastBackupName = importCache.BackupName
	r.BackupList = make([]string, 0)

	if r.LastJiraLocalHome != "" {
		res, err := utils.ListZipFile(common.GenExportPath(r.LastJiraLocalHome))
		if err != nil {
			return nil, err
		}
		r.BackupList = res
	}
	return r, nil
}
