package types

type ImportTask struct {
	Key             string
	Password        string
	UserUUID        string `json:"user_uuid"`
	ImportType      int    `json:"import_type"`
	Status          int    `json:"status"`
	ImportFromLocal bool   `json:"import_from_local"`
	LocalFilePath   string `json:"local_file_path"`
	AttachmentsPath string
	ImportTeamUUID  string

	SelectedProjectIDs map[string]bool       `json:"selected_project_ids"`
	BuiltinIssueTypes  []BuiltinIssueTypeMap `json:"builtin_issue_types"`
	MapFilePath        map[string]string
}

type ThirdResourceIDMap struct {
	ResourceID string
	OnesUUID   string
}

type BuiltinIssueTypeMap struct {
	IssueTypeID string `json:"id"`
	BuiltinType int    `json:"type"`
}
