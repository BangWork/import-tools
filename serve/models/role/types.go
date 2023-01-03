package role

const (
	RoleStatusNormal = 1
)

type Role struct {
	UUID            string `db:"uuid"`
	TeamUUID        string `db:"team_uuid"`
	Name            string `db:"name"`
	NamePinyin      string `db:"name_pinyin"`
	BuiltIn         bool   `db:"built_in"`
	IsProjectMember bool   `db:"is_project_member"`
	CreateTime      int64  `db:"create_time"`
	Status          int    `db:"status"`
}

type RoleConfig struct {
	ContainerUUID string `db:"container_uuid" json:"container_uuid"`
	RoleUUID      string `db:"role_uuid" json:"role_uuid"`
	CreateTime    int64  `db:"create_time" json:"create_time"`
}

type RoleMember struct {
	ContainerUUID string `db:"container_uuid"`
	RoleUUID      string `db:"role_uuid"`
	UserUUID      string `db:"user_uuid"`
}

type ProjectRoleKey struct {
	ProjectUUID string `db:"project_uuid"`
	RoleUUID    string `db:"role_uuid"`
}

func (p *ProjectRoleKey) StringKey() string {
	return p.ProjectUUID + "." + p.RoleUUID
}

func UniqueProjectRoleKeys(keys []*ProjectRoleKey) []*ProjectRoleKey {
	newKeys := []*ProjectRoleKey{}
	m := map[string]bool{}
	for _, key := range keys {
		if key == nil {
			continue
		}
		strKey := key.StringKey()
		if ok := m[strKey]; !ok {
			newKeys = append(newKeys, key)
			m[strKey] = true
		}
	}
	return newKeys
}
