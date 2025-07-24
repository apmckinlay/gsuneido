// Copyright (C) 2002 Axon Development Corporation All rights reserved worldwide.
ChooseListControl
	{
	sf: false
	New(fields = #(), table = false, exclude_fields = #(), .listField = false,
		.mandatory = false, width = 10)
		{
		super(.getlist(fields, table, exclude_fields), listSeparator: '', :width)
		.refreshListFromSelectFields()
		}

	refreshListFromSelectFields(rec = false)
		{
		if 0 isnt sf = .Send('FieldPrompt_GetSelectFields', :rec)
			{
			.sf = sf
			.setSelectFieldMap()
			}
		}

	setSelectFieldMap()
		{
		list = Object()
		.map = Object()
		for prompt in .sf.Prompts()
			{
			list.Add(prompt)
			.map[prompt] = .sf.Fields[prompt]
			}
		list.Sort!()
		.SetList(list)
		}

	getlist(fields, table, exclude_fields)
		{
		.validateParams(fields, table, .listField)
		list = Object()
		.map = Object()
		for field in fields
			.add_field(field, list)

		.addFieldsFromQuery(table, exclude_fields, list)

		if .listField isnt false
			{
			if 0 is fieldList = .Send('GetField', .listField)
				fieldList = Record()[.listField]
			if Object?(fieldList)
				fieldList.Each({|x| .add_field(x, list)})
			}

		return list.Sort!()
		}

	validateParams(fields, table, listField)
		{
		if listField isnt false and (not fields.Empty?() or table isnt false)
			{
			SuneidoLog('ERROR: FieldPromptControl - ' $
				'listField cannot be used in combination with table or fields', calls:)
			return false // return values are for test
			}
		return true // return values are for test
		}

	addFieldsFromQuery(table, exclude_fields, list)
		{
		if table is false
			return
		for field in QuerySelectColumns(table)
			if .allowAddField?(field, exclude_fields, .map)
				.add_field(field, list)
		}

	allowAddField?(field, exclude_fields, map)
		{
		if exclude_fields.Has?(field) or map.Has?(field)
			return false

		ddVals = Datadict(field, #(ExcludeSelect))
		return not ddVals.Member?('ExcludeSelect') or not ddVals.ExcludeSelect
		}

	add_field(field, list)
		{
		prompt = SelectFields.GetFieldPrompt(field, .map.Members())
		if prompt is field or prompt is "" or Customizable.DeletedField?(field)
			return
		list.Add(prompt)
		.map[prompt] = field
		}

	GetFields()
		{
		return .map.Values()
		}

	SetFieldMap(fields)
		{
		list = Object()
		.map = Object()
		for field in fields
			.add_field(field, list)
		list.Sort!()
		.SetList(list)
		}

	Set(value)
		{
		.refreshList()
		prompt = .map.Find(value)
		.Field.Set(prompt is false ? value : prompt)
		}

	RefreshListFromSelectFields(rec = false)
		{
		fieldBefore = .Get()
		fieldText = .Field.Get()
		.refreshListFromSelectFields(rec)
		fieldAfter = .Get()
		if fieldText isnt '' and fieldAfter isnt fieldBefore
			.NewValue(false)
		}

	refreshList()
		{
		if .listField isnt false and 0 isnt fieldList = .Send('GetField', .listField)
			.SetFieldMap(fieldList)
		.refreshListFromSelectFields()
		}

	Get()
		{
		return .map.GetDefault(.Field.Get(), '')
		}
	NewValue(value /*unused*/)
		{
		.Field.Dirty?(true)
		.Send("NewValue", .Get())
		}
	SetMandatory(.mandatory = false)
		{
		}
	Valid?()
		{
		prompt = .Field.Get()
		return prompt is '' ? not .mandatory : .map.Member?(prompt)
		}
	On_DropDown()
		{
		.refreshList()
		super.On_DropDown()
		}

	FieldSetFocus()
		{
		.refreshList()
		super.FieldSetFocus()
		}

	AliasParamSelectField(field = false)
		{
		if field is false
			field = .Get()
		if .sf is false or not .sf.NameAbbrev?(field)
			return field
		return field.Has?('_abbrev') ? field.Replace('_abbrev', '_num') : 'string'
		}
	}
