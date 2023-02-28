package services

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
	JiraList []*JiraIssueType `json:"jira_list"`
	ONESList []*ONESIssueType `json:"ones_list"`
}

type JiraIssueType struct {
	IssueTypeID    string `json:"third_issue_type_id"`
	IssueTypeName  string `json:"third_issue_type_name"`
	ONESDetailType int    `json:"ones_detail_type"`
}

type ONESIssueType struct {
	Name       string `json:"name"`
	DetailType int    `json:"detail_type"`
}
