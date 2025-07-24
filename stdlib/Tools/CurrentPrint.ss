// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(row, hwnd, query, name = false, excludeFields = #())
		{
		if name is false
			name = query
		sf = SelectFields(QueryColumns(query), excludeFields, joins: false)
		return ToolDialog(hwnd, .Params(row, sf, name))
		}

	Params(row, sf, name)
		{
		field_ob = Object()
		for f in sf.Prompts().Sort!()
			field_ob[sf.PromptToField(f)] = f.Trim()
		return Object('Params'
			Object(AccessCurrentPrint, row, sf.Fields)
			title: 'Print Current Record',
			name: 'Current Print - ' $ name,
			Params:	Object('Form'
				Object('Static', 'Columns' group: 0)
				Object('ChooseTwoList' field_ob,
					title: 'Columns',
					name: 'choosefields',
					group: 1) 'nl' ))
		}
	}