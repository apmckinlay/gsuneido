// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
ChooseTwoListControl
	{
	New(query = '', exclude = #(), extraFieldsPlugin = false, width = 50,
		mandatory = false, .listField = false, mandatoryListField = false)
		{
		super(.list(query, exclude, extraFieldsPlugin), :width, :mandatory,
			mandatory_list: .mandatoryList(mandatoryListField))
		.setList()
		}

	Getter_DialogControl()
		{
		.setList()
		super.Getter_DialogControl()
		}

	setList()
		{
		if .listField isnt false
			{
			.availableFields = .Send('GetField', .listField)
			.SetList(.promptsFromAvailableFields())
			}
		}

	promptsFromAvailableFields()
		{
		.promptMap = Object()
		for f in .availableFields
			.promptMap[f] = Datadict.GetFieldPrompt(f, .promptMap.Values())
		return .promptMap.Values()
		}

	availableFields: #()
	list(query, exclude, extraFieldsPlugin)
		{
		if .listField is false
			.availableFields = .selectFields(query, exclude, extraFieldsPlugin)
		return .promptsFromAvailableFields()
		}

	mandatoryList(mandatoryListField)
		{
		mandatoryList = mandatoryListField is false
			? #()
			: [][mandatoryListField]
		return mandatoryList.Map(SelectPrompt)
		}

	selectFields(query, exclude, extraFieldsPlugin)
		{
		fields = QueryColumns(query)
		if extraFieldsPlugin isnt false
			{
			Plugins().ForeachContribution(extraFieldsPlugin, 'fields')
				{ |x|
				fields.MergeUnion((x.func)())
				}
			}
		return SelectFields(fields, exclude, joins:).Fields.Values()
		}

	Get()
		{
		prompts = super.Get()
		fields = Object()
		for prompt in prompts.Split(',')
			for field in .availableFields
				if .promptMap[field] is prompt
					{
					fields.Add(field)
					break
					}
		return fields
		}

	GetSelectedList()
		{
		return super.Get().Split(',')
		}

	Set(fields)
		{
		// if listField is dependent on other record fields, need to refresh list since
		// the RecordControl may not have been "set" when this control was constructed
		.setList()
		super.Set(.convertToPrompts(fields))
		}

	convertToPrompts(fields)
		{
		if String?(fields)
			return fields

		prompts = Object()
		// excludeTags: esmpty to avoid ProgrammerErrors with fields selected then
		// subsequently marked internal. Field will be invalid and user can clean up
		for f in fields
			prompts.Add(.promptMap.GetDefault(
				f, Datadict.GetFieldPrompt(f, excludeTags: #())))
		return prompts.Join(',')
		}

	Valid?(forceCheck = false)
		{
		.setList()
		super.Valid?(forceCheck)
		}
	}
