package utils

import (
	"sort"

	"github.com/bangwork/import-tools/serve/utils"
)

type Tuple2_Int_Int struct {
	Ele_1 int64 `db:"_1"`
	Ele_2 int64 `db:"_2"`
}

type Tuple2_String_Int struct {
	Ele_1 string `db:"_1"`
	Ele_2 int64  `db:"_2"`
}

type Tuple2_Int_String struct {
	Ele_1 int64  `db:"_1"`
	Ele_2 string `db:"_2"`
}

type Tuple2_String_String struct {
	Ele_1 string `db:"_1"`
	Ele_2 string `db:"_2"`
}

type Tuple3_String_Int_Int struct {
	Ele_1 string `db:"_1"`
	Ele_2 int64  `db:"_2"`
	Ele_3 int64  `db:"_3"`
}
type Tuple3_String_String_Int struct {
	Ele_1 string `db:"_1"`
	Ele_2 string `db:"_2"`
	Ele_3 int64  `db:"_3"`
}

type Tuple3_Int_Int_Int struct {
	Ele_1 int64 `db:"_1"`
	Ele_2 int64 `db:"_2"`
	Ele_3 int64 `db:"_3"`
}

type Tuple3_String_String_String struct {
	Ele_1 string `db:"_1"`
	Ele_2 string `db:"_2"`
	Ele_3 string `db:"_3"`
}

type Tuple4_String_String_String_String struct {
	Ele_1 string `db:"_1"`
	Ele_2 string `db:"_2"`
	Ele_3 string `db:"_3"`
	Ele_4 string `db:"_4"`
}

type Tuple2_String_NString struct {
	Ele_1 string  `db:"_1"`
	Ele_2 *string `db:"_2"`
}

type Tuple2_NString_NString struct {
	Ele_1 *string `db:"_1"`
	Ele_2 *string `db:"_2"`
}

type Tuple3_String_NString_String struct {
	Ele_1 *string `db:"_1"`
	Ele_2 *string `db:"_2"`
	Ele_3 string  `db:"_3"`
}

type Tuple3_String_String_String_Sorter []*Tuple3_String_String_String

func (t Tuple3_String_String_String_Sorter) Len() int      { return len(t) }
func (t Tuple3_String_String_String_Sorter) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t Tuple3_String_String_String_Sorter) Less(i, j int) bool {
	return utils.CompareString(t[i].Ele_3, t[j].Ele_3)
}

type Tuple2_String_Int_Sorter []*Tuple2_String_Int

func (t Tuple2_String_Int_Sorter) Len() int      { return len(t) }
func (t Tuple2_String_Int_Sorter) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t Tuple2_String_Int_Sorter) Less(i, j int) bool {
	return t[i].Ele_2 < t[j].Ele_2
}

type NameWithSortKey struct {
	Name    string
	SortKey interface{}
}

type StatusWithSortKey struct {
	Status  int
	SortKey interface{}
}

func (k1 NameWithSortKey) Less(k2 NameWithSortKey) bool {
	switch k1.SortKey.(type) {
	case string:
		return utils.CompareString(k1.SortKey.(string), k2.SortKey.(string))
	case int64:
		return k1.SortKey.(int64) < k2.SortKey.(int64)
	default:
		panic("invalid sort key")
	}
}

type UUIDWithName struct {
	UUID string
	Name string
}

type userNameSorter struct {
	UserUUIDs       []string
	MapUserSortKeys map[string]NameWithSortKey
}

func SortUserUUIDsBySortKeysMap(userUUIDs []string, mapUserSortKeys map[string]NameWithSortKey) {
	sort.Sort(userNameSorter{
		UserUUIDs:       userUUIDs,
		MapUserSortKeys: mapUserSortKeys,
	})
}

func (m userNameSorter) Len() int { return len(m.UserUUIDs) }
func (m userNameSorter) Swap(i, j int) {
	m.UserUUIDs[i], m.UserUUIDs[j] = m.UserUUIDs[j], m.UserUUIDs[i]
}
func (m userNameSorter) Less(i, j int) bool {
	iUser := m.MapUserSortKeys[m.UserUUIDs[i]]
	jUser := m.MapUserSortKeys[m.UserUUIDs[j]]
	if iUser.SortKey == nil {
		return false
	} else if jUser.SortKey == nil {
		return true
	}
	return utils.CompareString(iUser.SortKey.(string), jUser.SortKey.(string))
}
