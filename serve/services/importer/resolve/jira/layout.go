package jira

import (
	"fmt"
	"strings"

	"github.com/bangwork/import-tools/serve/models/issuetype"
	"github.com/bangwork/import-tools/serve/services/importer/constants"
	"github.com/bangwork/import-tools/serve/services/importer/resolve"
	"github.com/bangwork/import-tools/serve/utils"
)

func (p *JiraResolver) prepareProjectIssueTypeLayout() error {
	var projectIssueTypeWithFields = make(map[string][]*resolve.ScopeConfigItem, 0)
	for k := range p.jiraProjectIssueTypeFieldScreen {
		parts := strings.Split(k, ":")
		if len(parts) != 4 {
			continue
		}
		projectID, issueTypeID, fieldLayoutID, FieldScreenSchemeID := parts[0], parts[1], parts[2], parts[3]
		if p.ignoreProjectID(projectID) {
			continue
		}
		k := fmt.Sprintf("%s:%s", projectID, issueTypeID)
		if _, ok := p.jiraProjectIssueTypeNameMap[k]; !ok {
			continue
		}

		a := &resolve.ScopeConfigItem{
			IssueTypeID:    issueTypeID,
			FieldConfigID:  fieldLayoutID,
			ScreenSchemeID: FieldScreenSchemeID,
		}

		projectIssueTypeWithFields[projectID] = append(projectIssueTypeWithFields[projectID], a)
	}

	for projectID, issueTypeWithFields := range projectIssueTypeWithFields {
		content := &resolve.ThirdProjectIssueTypeLayout{
			ProjectID:    projectID,
			ScopeConfigs: issueTypeWithFields,
		}
		p.jiraProjectIssueTypeLayoutSlice = append(p.jiraProjectIssueTypeLayoutSlice, content)
	}
	return nil
}

func (p *JiraResolver) prepareProjectIssueTypeWithFields() error {
	var localScreen = make(map[string][]*resolve.IssueTypeWithFields, 0)

	// local screen
	for k := range p.jiraProjectIssueTypeFieldScreen {
		parts := strings.Split(k, ":")
		if len(parts) != 4 {
			continue
		}
		projectID, issueTypeID, fieldLayoutID, FieldScreenSchemeID := parts[0], parts[1], parts[2], parts[3]
		if p.ignoreProjectID(projectID) {
			continue
		}
		k := fmt.Sprintf("%s:%s", projectID, issueTypeID)
		if _, ok := p.jiraProjectIssueTypeNameMap[k]; !ok {
			continue
		}

		fieldLayout, ok := p.jiraFieldLayoutIDWithLayout[fieldLayoutID]
		if !ok {
			continue
		}

		// Get all the attribute information of the configuration scheme through the id of fieldLayout, there will be whether the attribute is required and hidden information
		items, ok := p.jiraFieldLayoutIDWithItem[fieldLayout.ID]
		if !ok {
			continue
		}

		//pp.Println(issueTypeID, fieldLayout.Name, fieldScreenSchemeName)

		// field config fields map[FieldIdentifier] = FieldLayoutItem
		var mapFieldIdentifierWithItem = make(map[string]*fieldLayoutItem)
		for _, item := range items {
			mapFieldIdentifierWithItem[item.FieldIdentifier] = item
		}

		// Use FieldScreenSchemeID to obtain the creation, modification, and details screen corresponding to the issue type
		screenItems, ok := p.jiraFieldScreenSchemeIDWithItem[FieldScreenSchemeID]
		if !ok {
			continue
		}

		// 0: Create, 1：Edit, 2：View
		var operations = []string{"0", "1", "2"}
		var defaultFieldScreenSchemeItem *fieldScreenSchemeItem
		//
		var allScreenItems = make([]*fieldScreenSchemeItem, 0)

		var existOperations = make([]string, 0)
		for _, s := range screenItems {
			if s.Operation != "" {
				existOperations = append(existOperations, s.Operation)
				allScreenItems = append(allScreenItems, s)
			}

			if s.Operation == "" {
				defaultFieldScreenSchemeItem = s
			}
		}

		//pp.Println("defaultFieldScreenSchemeItem: ", defaultFieldScreenSchemeItem)

		// 算出没有配置 screen 的 operation
		_, unmapedOperations := utils.StringArrayDifference(operations, existOperations)
		//pp.Println("unmapedOperations: ", unmapedOperations)

		if defaultFieldScreenSchemeItem != nil && len(unmapedOperations) != 0 {
			for _, op := range unmapedOperations {
				item := &fieldScreenSchemeItem{
					ID:                  defaultFieldScreenSchemeItem.ID,
					Operation:           op,
					FieldScreenID:       defaultFieldScreenSchemeItem.FieldScreenID,
					FieldScreenSchemeID: defaultFieldScreenSchemeItem.FieldScreenSchemeID,
				}
				allScreenItems = append(allScreenItems, item)
			}
		}

		//pp.Println("allScreenItems: ", allScreenItems)

		var mapFields = make(map[string]*resolve.Field, 0)

		for _, s := range allScreenItems {
			// According to FieldScreenID, go to the corresponding field layout tabs (that is, screen configuration)
			tabIDs, ok := p.jiraFieldScreenIDWithTabIds[s.FieldScreenID]
			if !ok {
				continue
			}

			for _, tabID := range tabIDs {
				tabItems, ok := p.jiraFieldScreenTabIDWithItem[tabID]
				if !ok {
					continue
				}

				for _, tabItem := range tabItems {
					fieldItem, ok := mapFieldIdentifierWithItem[tabItem.FieldIdentifier]
					if !ok {
						// In special cases, the added custom field will not be recorded in the field layout if it has not been modified in the field config
						fieldItem = &fieldLayoutItem{
							ID:              "",
							FieldLayoutID:   "",
							FieldIdentifier: tabItem.FieldIdentifier,
							IsHidden:        "false",
							IsRequired:      "false",
						}
					}

					if fieldItem.IsHidden == "true" {
						continue
					}

					var required bool
					if fieldItem.IsRequired == "true" {
						required = true
					}

					field := &resolve.Field{
						FieldIdentifier: fieldItem.FieldIdentifier,
						Required:        required,
					}

					mapFields[fieldItem.FieldIdentifier] = field
				}
			}
		}

		addFieldByIdentifier("versions", mapFieldIdentifierWithItem, mapFields)
		addFieldByIdentifier("assignee", mapFieldIdentifierWithItem, mapFields)
		addFieldByIdentifier("attachment", mapFieldIdentifierWithItem, mapFields)
		addFieldByIdentifier("comment", mapFieldIdentifierWithItem, mapFields)
		addFieldByIdentifier("description", mapFieldIdentifierWithItem, mapFields)
		addFieldByIdentifier("duedate", mapFieldIdentifierWithItem, mapFields)
		addFieldByIdentifier("fixVersions", mapFieldIdentifierWithItem, mapFields)
		addFieldByIdentifier("issuetype", mapFieldIdentifierWithItem, mapFields)
		addFieldByIdentifier("labels", mapFieldIdentifierWithItem, mapFields)
		addFieldByIdentifier("issuelinks", mapFieldIdentifierWithItem, mapFields)
		addFieldByIdentifier("priority", mapFieldIdentifierWithItem, mapFields)
		addFieldByIdentifier("resolution", mapFieldIdentifierWithItem, mapFields)
		addFieldByIdentifier("summary", mapFieldIdentifierWithItem, mapFields)
		addFieldByIdentifier("components", mapFieldIdentifierWithItem, mapFields)
		addFieldByIdentifier("reporter", mapFieldIdentifierWithItem, mapFields)

		sprintField := &resolve.Field{
			FieldIdentifier: "sprint",
			Required:        false,
		}
		mapFields[sprintField.FieldIdentifier] = sprintField

		idField := &resolve.Field{
			FieldIdentifier: "id",
			Required:        false,
		}
		mapFields[idField.FieldIdentifier] = idField

		var fields = make([]*resolve.Field, 0)
		for _, f := range mapFields {
			fields = append(fields, f)
		}

		d := &resolve.IssueTypeWithFields{
			IssueTypeID: issueTypeID,
			Fields:      fields,
		}

		localScreen[projectID] = append(localScreen[projectID], d)

		//if issueTypeName == "Task" && fieldLayout.Name == "testFieldConfiguration"{
		//	for _, f := range mapFields {
		//		pp.Println(f.FieldIdentifier)
		//	}
		//	pp.Println("len: ", len(mapFields))
		//	//pp.Println(mapFields)
		//	pp.Println("------------------------------------------------------------------")
		//}
	}

	for projectID, issueTypeWithFields := range localScreen {
		content := &resolve.ThirdProjectIssueTypeWithFields{
			ProjectID:           projectID,
			IssueTypeWithFields: issueTypeWithFields,
		}
		p.jiraProjectIssueTypeWithFieldsSlice = append(p.jiraProjectIssueTypeWithFieldsSlice, content)

		content = &resolve.ThirdProjectIssueTypeWithFields{
			ProjectID: projectID,
			IssueTypeWithFields: []*resolve.IssueTypeWithFields{
				{
					IssueTypeID: "",
					Fields: []*resolve.Field{
						{
							FieldIdentifier: customFieldReleaseStartDate,
							Required:        false,
						},
					},
					IssueTypeDetailType: issuetype.DetailTypePublish,
				},
			},
		}
		p.jiraProjectIssueTypeWithFieldsSlice = append(p.jiraProjectIssueTypeWithFieldsSlice, content)
	}
	return nil
}

func (p *JiraResolver) prepareIssueTypeLayout() error {
	projectIssueTypeFieldLayout := make(map[string]string)
	issueTypeWithFieldLayout := make(map[string]string)
	defaultFieldLayout := p.DefaultFieldLayout()

	//pp.Println("defaultFieldLayout", defaultFieldLayout)

	for projectID, _ := range p.jiraProjectIDNameMap {
		fieldLayoutSchemeID, ok := p.jiraProjectWithFieldLayoutScheme[projectID]
		if !ok {
			// 处理使用默认的 FieldLayoutScheme, fieldSystem Default Field Configuration
			issueTypes, ok := p.jiraProjectWithIssueTypeMap[projectID]
			if !ok {
				continue
			}

			for _, issueTypeId := range issueTypes {
				k1 := fmt.Sprintf("%s:%s", issueTypeId, defaultFieldLayout.ID)
				k2 := fmt.Sprintf("%s:%s", projectID, k1)
				issueTypeWithFieldLayout[k1] = k1
				projectIssueTypeFieldLayout[k2] = k2
			}
		} else {
			entitys, ok := p.jiraFieldLayoutSchemeIDWithEntity[fieldLayoutSchemeID]
			if !ok {
				continue
			}

			unmappedIssueTypeIDs, ok := p.unmappedIssueTypeIDsFromEntitys(projectID, entitys)
			if !ok {
				continue
			}

			for _, e := range entitys {
				if e.IssueTypeID == "" {
					// Default Used for all unmapped issue types
					for _, issueTypeID := range unmappedIssueTypeIDs {
						var k1 string
						if e.FieldLayoutID == "" {
							//	没有 fieldlayout 属性默认使用 Default Field Configuration
							//pp.Println(e)
							k1 = fmt.Sprintf("%s:%s", issueTypeID, defaultFieldLayout.ID)

						} else {
							// 有 fieldlayout 属性值
							if _, ok := p.jiraFieldLayoutIDWithLayout[e.FieldLayoutID]; !ok {
								continue
							}
							k1 = fmt.Sprintf("%s:%s", issueTypeID, e.FieldLayoutID)
						}

						issueTypeWithFieldLayout[k1] = k1

						k2 := fmt.Sprintf("%s:%s", projectID, k1)
						projectIssueTypeFieldLayout[k2] = k2
					}
				} else {
					var k1 string

					// 有 issue type 属性值
					if e.FieldLayoutID == "" {
						//defaultFieldLayout.ID
						//	没有 fieldlayout 属性默认使用 Default Field Configuration
						k1 = fmt.Sprintf("%s:%s", e.IssueTypeID, defaultFieldLayout.ID)

					} else {
						// 有 fieldlayout 属性值
						if _, ok := p.jiraFieldLayoutIDWithLayout[e.FieldLayoutID]; !ok {
							continue
						}

						k1 = fmt.Sprintf("%s:%s", e.IssueTypeID, e.FieldLayoutID)
					}

					issueTypeWithFieldLayout[k1] = k1

					k2 := fmt.Sprintf("%s:%s", projectID, k1)
					projectIssueTypeFieldLayout[k2] = k2
				}
			}
		}

		issueTypeScreenSchemeID, ok := p.jiraProjectWithIssueTypeScreenScheme[projectID]
		if !ok {
			continue
		}

		entitys2, ok := p.jiraIssueTypeScreenSchemeIDWithEntity[issueTypeScreenSchemeID]
		if !ok {
			continue
		}

		unmappedIssueTypeIDs2, ok := p.unmappedIssueTypeIDsFromIssueTypeScreenSchemeEntitys(projectID, entitys2)
		if !ok {
			continue
		}

		for _, e := range entitys2 {
			if e.IssueTypeID == "" {
				for _, issueTypeID := range unmappedIssueTypeIDs2 {
					if _, ok := p.jiraFieldScreenSchemeIDWithName[e.FieldScreenSchemeID]; !ok {
						continue
					}

					prefix := fmt.Sprintf("%s:%s:", projectID, issueTypeID)
					k := getKeyByPrefix(projectIssueTypeFieldLayout, prefix)
					if len(k) == 0 {
						continue
					}

					k2 := fmt.Sprintf("%s:%s", k, e.FieldScreenSchemeID)
					p.jiraProjectIssueTypeFieldScreen[k2] = k2
				}
			} else {
				if _, ok := p.jiraFieldScreenSchemeIDWithName[e.FieldScreenSchemeID]; !ok {
					continue
				}

				// 检查项目中是否有应用了这个 issue type
				projectIssueTypeKey := fmt.Sprintf("%s:%s", projectID, e.IssueTypeID)
				if _, ok := p.jiraProjectIssueTypeNameMap[projectIssueTypeKey]; !ok {
					continue
				}

				prefix := fmt.Sprintf("%s:%s:", projectID, e.IssueTypeID)
				k := getKeyByPrefix(projectIssueTypeFieldLayout, prefix)
				if len(k) == 0 {
					continue
				}

				k2 := fmt.Sprintf("%s:%s", k, e.FieldScreenSchemeID)
				p.jiraProjectIssueTypeFieldScreen[k2] = k2
			}
		}
	}

	var issueTypeFieldScreen = make(map[string]string)

	for k := range p.jiraProjectIssueTypeFieldScreen {
		parts := strings.Split(k, ":")
		if len(parts) != 4 {
			continue
		}
		k2 := fmt.Sprintf("%s:%s:%s", parts[1], parts[2], parts[3])
		issueTypeFieldScreen[k2] = k2
	}

	// 循环所有的 issue type + field Layout + field screen scheme 的组合
	// k 的组成 "{issue type id}:{field layout id}:{field screen scheme id}" 的组合，例如 "10000:10001:10002"
	for k := range issueTypeFieldScreen {
		parts := strings.Split(k, ":")
		if len(parts) != 3 {
			continue
		}
		issueTypeID, fieldLayoutID, FieldScreenSchemeID := parts[0], parts[1], parts[2]

		// 取 field layout，field layout 就是 工作项类型对应的属性配置方案
		fieldLayout, ok := p.jiraFieldLayoutIDWithLayout[fieldLayoutID]
		if !ok {
			continue
		}

		// 通过 fieldLayout 的 id 取配置方案的所有属性信息，里面会有属性是否必填和隐藏信息
		items, ok := p.jiraFieldLayoutIDWithItem[fieldLayout.ID]
		if !ok {
			continue
		}

		fieldScreenSchemeName, ok := p.jiraFieldScreenSchemeIDWithName[FieldScreenSchemeID]
		if !ok {
			continue
		}

		//pp.Println(issueTypeID, fieldLayout.Name, fieldScreenSchemeName)

		var mapFieldIdentifierWithItem = make(map[string]*fieldLayoutItem)

		// 过滤掉在 field layout 中隐藏的属性
		for _, item := range items {
			mapFieldIdentifierWithItem[item.FieldIdentifier] = item
		}

		// 通过 FieldScreenSchemeID 获取 issue type 对应的创建，修改，详情 screen
		screenItems, ok := p.jiraFieldScreenSchemeIDWithItem[FieldScreenSchemeID]
		if !ok {
			continue
		}

		// 0: Create, 1：Edit, 2：View
		var operations = []string{"0", "1", "2"}
		var defaultFieldScreenSchemeItem *fieldScreenSchemeItem
		//
		var allScreenItems = make([]*fieldScreenSchemeItem, 0)

		var existOperations = make([]string, 0)
		for _, s := range screenItems {
			if s.Operation != "" {
				existOperations = append(existOperations, s.Operation)
				allScreenItems = append(allScreenItems, s)
			}

			if s.Operation == "" {
				defaultFieldScreenSchemeItem = s
			}
		}

		//pp.Println("defaultFieldScreenSchemeItem: ", defaultFieldScreenSchemeItem)

		// 算出没有配置 screen 的 operation
		_, unmapedOperations := utils.StringArrayDifference(operations, existOperations)

		//pp.Println("unmapedOperations: ", unmapedOperations)

		if defaultFieldScreenSchemeItem != nil && len(unmapedOperations) != 0 {
			for _, op := range unmapedOperations {
				item := &fieldScreenSchemeItem{
					ID:                  defaultFieldScreenSchemeItem.ID,
					Operation:           op,
					FieldScreenID:       defaultFieldScreenSchemeItem.FieldScreenID,
					FieldScreenSchemeID: defaultFieldScreenSchemeItem.FieldScreenSchemeID,
				}
				allScreenItems = append(allScreenItems, item)
			}
		}

		//pp.Println("allScreenItems: ", allScreenItems)

		var mapCreateIssueFields = make(map[string]*resolve.Field, 0)
		var mapViewIssueFields = make(map[string]*resolve.Field, 0)

		for _, s := range allScreenItems {
			// 根据 FieldScreenID 去到对应的  field layout tabs (就是 screen 配置）
			tabIDs, ok := p.jiraFieldScreenIDWithTabIds[s.FieldScreenID]
			if !ok {
				continue
			}

			for _, tabID := range tabIDs {
				// 用 tab id 取 tab 对应配置的属性
				tabItems, ok := p.jiraFieldScreenTabIDWithItem[tabID]
				if !ok {
					continue
				}

				for _, tabItem := range tabItems {
					fieldItem, ok := mapFieldIdentifierWithItem[tabItem.FieldIdentifier]
					if !ok {
						// 特殊情况，添加的 custom field，如果没有在 field config 中修改过，不会在 field layout 中有记录
						fieldItem = &fieldLayoutItem{
							ID:              "",
							FieldLayoutID:   "",
							FieldIdentifier: tabItem.FieldIdentifier,
							IsHidden:        "false",
							IsRequired:      "false",
						}
					}

					if fieldItem.IsHidden == "true" {
						continue
					}

					var required bool
					if fieldItem.IsRequired == "true" {
						required = true
					}

					field := &resolve.Field{
						FieldIdentifier: fieldItem.FieldIdentifier,
						Required:        required,
					}

					if s.Operation == "0" {
						mapCreateIssueFields[fieldItem.FieldIdentifier] = field

						//createFields = append(createFields, field)
					} else {
						mapViewIssueFields[fieldItem.FieldIdentifier] = field
					}
				}
			}
		}

		addFieldByIdentifier("versions", mapFieldIdentifierWithItem, mapViewIssueFields)
		addFieldByIdentifier("assignee", mapFieldIdentifierWithItem, mapViewIssueFields)
		addFieldByIdentifier("attachment", mapFieldIdentifierWithItem, mapViewIssueFields)
		addFieldByIdentifier("comment", mapFieldIdentifierWithItem, mapViewIssueFields)
		addFieldByIdentifier("description", mapFieldIdentifierWithItem, mapViewIssueFields)
		addFieldByIdentifier("duedate", mapFieldIdentifierWithItem, mapViewIssueFields)
		addFieldByIdentifier("fixVersions", mapFieldIdentifierWithItem, mapViewIssueFields)
		addFieldByIdentifier("issuetype", mapFieldIdentifierWithItem, mapViewIssueFields)
		addFieldByIdentifier("labels", mapFieldIdentifierWithItem, mapViewIssueFields)
		addFieldByIdentifier("issuelinks", mapFieldIdentifierWithItem, mapViewIssueFields)
		addFieldByIdentifier("priority", mapFieldIdentifierWithItem, mapViewIssueFields)
		addFieldByIdentifier("resolution", mapFieldIdentifierWithItem, mapViewIssueFields)
		addFieldByIdentifier("summary", mapFieldIdentifierWithItem, mapViewIssueFields)
		addFieldByIdentifier("reporter", mapFieldIdentifierWithItem, mapViewIssueFields)

		addFieldByIdentifier("assignee", mapFieldIdentifierWithItem, mapCreateIssueFields)
		addFieldByIdentifier("description", mapFieldIdentifierWithItem, mapCreateIssueFields)
		addFieldByIdentifier("duedate", mapFieldIdentifierWithItem, mapCreateIssueFields)
		addFieldByIdentifier("issuetype", mapFieldIdentifierWithItem, mapCreateIssueFields)
		addFieldByIdentifier("priority", mapFieldIdentifierWithItem, mapCreateIssueFields)
		addFieldByIdentifier("summary", mapFieldIdentifierWithItem, mapCreateIssueFields)
		addFieldByIdentifier("reporter", mapFieldIdentifierWithItem, mapCreateIssueFields)

		sprintField := &resolve.Field{
			FieldIdentifier: "sprint",
			Required:        false,
		}
		mapViewIssueFields[sprintField.FieldIdentifier] = sprintField
		mapCreateIssueFields[sprintField.FieldIdentifier] = sprintField

		idField := &resolve.Field{
			FieldIdentifier: "id",
			Required:        false,
		}
		mapViewIssueFields[idField.FieldIdentifier] = idField

		var viewFields = make([]*resolve.Field, 0)
		for _, f := range mapViewIssueFields {
			viewFields = append(viewFields, f)
		}

		var createFields = make([]*resolve.Field, 0)
		for _, f := range mapCreateIssueFields {
			createFields = append(createFields, f)
		}

		sv := &resolve.ThirdIssueTypeLayout{
			IssueTypeID:       issueTypeID,
			FieldConfigID:     fieldLayout.ID,
			FieldConfigName:   fieldLayout.Name,
			ScreenSchemeID:    FieldScreenSchemeID,
			ScreenSchemeName:  fieldScreenSchemeName,
			CreateIssueConfig: createFields,
			ViewIssueConfig:   viewFields,
		}
		p.jiraIssueTypeLayoutSlice = append(p.jiraIssueTypeLayoutSlice, sv)
	}

	return nil
}

func (p *JiraResolver) DefaultFieldLayout() fieldLayout {
	var r fieldLayout
	for _, layout := range p.jiraFieldLayoutIDWithLayout {
		if layout.Type == "default" {
			r = *layout
			break
		}
	}
	return r
}

func (p *JiraResolver) prepareFieldScreenLayoutItem() error {
	for {
		element, err := p.nextElement("FieldScreenLayoutItem")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}

		m := new(fieldScreenLayoutItem)
		m.ID = getAttributeValue(element, "id")
		m.FieldIdentifier = getAttributeValue(element, "fieldidentifier")
		m.FieldScreenTabID = getAttributeValue(element, "fieldscreentab")
		m.Sequence = getAttributeValue(element, "sequence")
		p.jiraFieldScreenTabIDWithItem[m.FieldScreenTabID] = append(p.jiraFieldScreenTabIDWithItem[m.FieldScreenTabID], m)
	}
	return nil
}
func (p *JiraResolver) prepareFieldLayoutItem() error {
	for {
		element, err := p.nextElement("FieldLayoutItem")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}

		m := new(fieldLayoutItem)
		m.ID = getAttributeValue(element, "id")
		m.FieldLayoutID = getAttributeValue(element, "fieldlayout")
		m.FieldIdentifier = getAttributeValue(element, "fieldidentifier")
		m.IsHidden = getAttributeValue(element, "ishidden")
		m.IsRequired = getAttributeValue(element, "isrequired")
		p.jiraFieldLayoutIDWithItem[m.FieldLayoutID] = append(p.jiraFieldLayoutIDWithItem[m.FieldLayoutID], m)
	}
	return nil
}
func (p *JiraResolver) prepareIssueTypeScreenSchemeEntity() error {
	for {
		element, err := p.nextElement("IssueTypeScreenSchemeEntity")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}

		m := new(issueTypeScreenSchemeEntity)
		m.ID = getAttributeValue(element, "id")
		m.IssueTypeScreenSchemeID = getAttributeValue(element, "scheme")
		m.FieldScreenSchemeID = getAttributeValue(element, "fieldscreenscheme")
		m.IssueTypeID = getAttributeValue(element, "issuetype")
		p.jiraIssueTypeScreenSchemeIDWithEntity[m.IssueTypeScreenSchemeID] = append(p.jiraIssueTypeScreenSchemeIDWithEntity[m.IssueTypeScreenSchemeID], m)
	}
	return nil
}
func (p *JiraResolver) prepareFieldConfigScheme() error {
	for {
		element, err := p.nextElement("FieldConfigScheme")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}

		id := getAttributeValue(element, "id")
		field := getAttributeValue(element, "fieldid")
		if strings.HasPrefix(field, constants.CustomFieldPrefix) {
			p.jiraFieldConfigSchemeMap[id] = field
		}
	}
	return nil
}
func (p *JiraResolver) prepareFieldScreenSchemeItem() error {
	for {
		element, err := p.nextElement("FieldScreenSchemeItem")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}

		m := new(fieldScreenSchemeItem)
		m.ID = getAttributeValue(element, "id")
		m.FieldScreenSchemeID = getAttributeValue(element, "fieldscreenscheme")
		m.FieldScreenID = getAttributeValue(element, "fieldscreen")
		m.Operation = getAttributeValue(element, "operation")
		p.jiraFieldScreenSchemeIDWithItem[m.FieldScreenSchemeID] = append(p.jiraFieldScreenSchemeIDWithItem[m.FieldScreenSchemeID], m)
	}
	return nil
}

func (p *JiraResolver) prepareFieldScreenTab() error {
	for {
		element, err := p.nextElement("FieldScreenTab")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}

		id := getAttributeValue(element, "id")
		fieldScreenID := getAttributeValue(element, "fieldscreen")
		p.jiraFieldScreenIDWithTabIds[fieldScreenID] = append(p.jiraFieldScreenIDWithTabIds[fieldScreenID], id)
	}
	return nil
}

func (p *JiraResolver) prepareFieldScreenScheme() error {
	for {
		element, err := p.nextElement("FieldScreenScheme")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}

		id := getAttributeValue(element, "id")
		name := getAttributeValue(element, "name")
		p.jiraFieldScreenSchemeIDWithName[id] = name
	}
	return nil
}

func (p *JiraResolver) prepareFieldLayoutSchemeEntity() error {
	for {
		element, err := p.nextElement("FieldLayoutSchemeEntity")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}

		m := new(fieldLayoutSchemeEntity)
		m.ID = getAttributeValue(element, "id")
		m.FieldLayoutSchemeID = getAttributeValue(element, "scheme")
		m.IssueTypeID = getAttributeValue(element, "issuetype")
		m.FieldLayoutID = getAttributeValue(element, "fieldlayout")
		if _, found := p.jiraFieldLayoutSchemeIDWithEntity[m.FieldLayoutSchemeID]; !found {
			p.jiraFieldLayoutSchemeIDWithEntity[m.FieldLayoutSchemeID] = make([]*fieldLayoutSchemeEntity, 0)
		}
		p.jiraFieldLayoutSchemeIDWithEntity[m.FieldLayoutSchemeID] = append(p.jiraFieldLayoutSchemeIDWithEntity[m.FieldLayoutSchemeID], m)
	}
	return nil
}

func (p *JiraResolver) prepareFieldLayout() error {
	for {
		element, err := p.nextElement("FieldLayout")
		if err != nil {
			return err
		}
		if element == nil {
			break
		}

		m := new(fieldLayout)
		m.ID = getAttributeValue(element, "id")
		m.Name = getAttributeValue(element, "name")
		m.Type = getAttributeValue(element, "type")
		p.jiraFieldLayoutIDWithLayout[m.ID] = m
	}
	return nil
}

func getKeyByPrefix(m map[string]string, prefix string) string {
	var r string
	for k := range m {
		if strings.HasPrefix(k, prefix) {
			r = k
		}
	}
	return r
}

func (p *JiraResolver) unmappedIssueTypeIDsFromIssueTypeScreenSchemeEntitys(projectID string, entitys []*issueTypeScreenSchemeEntity) ([]string, bool) {
	var newIssueTypeIDs = make([]string, 0)
	for _, e := range entitys {
		if e.IssueTypeID != "" {
			newIssueTypeIDs = append(newIssueTypeIDs, e.IssueTypeID)
		}
	}

	issueTypeIDs, ok := p.jiraProjectWithIssueTypeMap[projectID]
	if !ok {
		return nil, false
	}
	_, deletions := utils.StringArrayDifference(issueTypeIDs, newIssueTypeIDs)
	return deletions, true
}

func (p *JiraResolver) unmappedIssueTypeIDsFromEntitys(projectID string, entitys []*fieldLayoutSchemeEntity) ([]string, bool) {
	var newIssueTypeIDs = make([]string, 0)
	for _, e := range entitys {
		if e.IssueTypeID != "" {
			newIssueTypeIDs = append(newIssueTypeIDs, e.IssueTypeID)
		}
	}

	issueTypeIDs, ok := p.jiraProjectWithIssueTypeMap[projectID]
	if !ok {
		return nil, false
	}
	_, deletions := utils.StringArrayDifference(issueTypeIDs, newIssueTypeIDs)
	return deletions, true
}

func addFieldByIdentifier(identifier string, mapFieldIdentifierWithItem map[string]*fieldLayoutItem, mapIssueFields map[string]*resolve.Field) {
	field, ok := mapFieldIdentifierWithItem[identifier]
	if ok && field.IsHidden == "false" {
		var required bool
		if field.IsRequired == "true" {
			required = true
		}

		d := &resolve.Field{
			FieldIdentifier: field.FieldIdentifier,
			Required:        required,
		}

		mapIssueFields[identifier] = d
	}
}
