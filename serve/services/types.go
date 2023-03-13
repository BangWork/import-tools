package services

type StartResolveRequest struct {
	JiraLocalHome string `json:"jira_local_home"`
	BackupName    string `json:"backup_name"`
}

type LoginRequest struct {
	URL      string `json:"url"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ConfirmImportRequest struct {
	ImportType      int    `json:"import_type"`
	ServerID        string `json:"server_id"`
	Version         string `json:"version"`
	ResourceUUID    string `json:"resource_uuid"`
	ImportUUID      string `json:"import_uuid"`
	FromImportTools bool   `json:"from_import_tools"`
}

type QueueListResponse struct {
	BatchTasks []*QueueStruct `json:"batch_tasks"`
}

type QueueStruct struct {
	Extra     string `json:"extra"`
	JobStatus string `json:"job_status"`
}

type IssueTypeListResponse struct {
	MigratedList     []*MigratedList     `json:"migrated_list"`
	ReadyMigrateList []*ReadyMigrateList `json:"ready_migrate_list"`
}

type MigratedList struct {
	IssueTypeID       string `json:"third_issue_type_id"`
	IssueTypeName     string `json:"third_issue_type_name"`
	ONESIssueTypeName string `json:"ones_issue_type_name"`
	Action            string `json:"action"`
}

type ReadyMigrateList struct {
	IssueTypeID   string `json:"third_issue_type_id"`
	IssueTypeName string `json:"third_issue_type_name"`
	Type          int    `json:"type"`
	Available     bool   `json:"available"`
}
