package objectlinktype

const (
	LinkModelTwoWayMany = 212
	LinkModelManyToMany = 202

	columns = "uuid, team_uuid, name, name_pinyin, link_model, source_type, source_condition, link_out_desc, link_out_desc_pinyin, " + "target_type, target_condition, link_in_desc, link_in_desc_pinyin, built_in, create_time, status"
)
