package controllers

import "github.com/bangwork/import-tools/serve/services/importer/types"

type path struct {
	Path string `json:"path"`
}

type SetShareDiskRequest struct {
	Path         string `json:"path"`
	UseShareDisk bool   `json:"use_share_disk"`
}

type ChooseTeamRequest struct {
	Key      string `json:"key"`
	TeamUUID string `json:"team_uuid"`
	TeamName string `json:"team_name"`
}

type StartImportRequest struct {
	Key          string                      `json:"key"`
	Password     string                      `json:"password"`
	ProjectIDs   []string                    `json:"project_ids"`
	IssueTypeMap []types.BuiltinIssueTypeMap `json:"issue_type_map"`
}

type GetIssueTypeRequest struct {
	Key        string   `json:"key"`
	ProjectIDs []string `json:"project_ids"`
}

type ResolveProgressResponse struct {
	MultiTeam    bool   `json:"multi_team"`
	TeamName     string `json:"team_name"`
	OrgName      string `json:"org_name"`
	StartTime    int64  `json:"start_time"`
	Status       int    `json:"status"`
	ExpectedTime int64  `json:"expected_time"`
	SpentTime    int64  `json:"spent_time"`
	BackupName   string `json:"backup_name"`
}
