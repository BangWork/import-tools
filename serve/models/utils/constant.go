package utils

const (
	// project
	ProjectStatusNormal  = 1
	ProjectStatusDeleted = 2
)

var (
	MapTeamStatusLabel = map[int]string{
		TeamStatusNormal:   LabelTeamStatusNormal,
		TeamStatusDisabled: LabelTeamStatusDisabled,
		TeamStatusPending:  LabelTeamStatusPending,
	}
)

// team
const (
	LabelTeamStatusNormal   = "normal"
	LabelTeamStatusDisabled = "disabled"
	LabelTeamStatusPending  = "pending"

	TeamStatusNormal   = 1
	TeamStatusDisabled = 2
	TeamStatusPending  = 3
)

var (
	MapOrgStautsLabel = map[int]string{
		OrgStatusNormal:   LabelOrgStatusNormal,
		OrgStatusDisabled: LabelOrgStatusDisabled,
		OrgStatusPending:  LabelOrgStatusPending,
	}
)

// organization
const (
	LabelOrgStatusNormal   = "normal"
	LabelOrgStatusDisabled = "disabled"
	LabelOrgStatusPending  = "pending"

	OrgStatusNormal   = 1
	OrgStatusDisabled = 2
	OrgStatusPending  = 3
)

// user domain
const (
	UserGroupStatusNormal  = 1
	UserGroupStatusDeleted = 2
)

const (
	SprintStatusCategoryToDo       = 1
	SprintStatusCategoryInProgress = 2
	SprintStatusCategoryDone       = 3
)
