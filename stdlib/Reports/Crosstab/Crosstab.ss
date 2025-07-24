// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(query, title = "", excludeFields = #(), columns = false)
		{
		if title is ""
			title = QueryGetTable(query, nothrow:)
		query = QueryStripSort(query)
		if columns is false
			columns = QueryColumns(query)
		excludeFields = ExcludeLowerFields(columns, excludeFields)
		selectFields = SelectFields(columns, excludeFields)
		prompts = selectFields.Prompts().Sort!()
		funcs = .crossTabFuncs()
		skipSave? = title is false
		return Object('Params'
			Object(.crossTabFormat, :query, :selectFields)
			title: "CrossTable",
			NoPresets: skipSave?,
			NoSaveLoadParams: skipSave?,
			name: "CrossTable - " $ title,
			validField: "crosstable_valid",
			Params: Object('Vert'
				Object('Pair', #(Static, Rows)
					Object('ChooseList', prompts, name: "Rows"))
				Object('Pair', #('Static', 'Columns')
					Object('ChooseList', prompts, name: "Columns"))
				Object('Pair', #('Static', 'Value')
				Object('ChooseList', prompts, name: "Value"))
				Object('Pair', #('Static', 'Function')
					Object('ChooseList', funcs,	name: "Function" mandatory:))
				xstretch: 0, ystretch: 0),
			printParams: #(Rows, Columns, Value, Function)
			)
		}

	crossTabFuncs()
		{
		funcs = Object()
		for m in Accumulator.Members()
			if m.Prefix?("F_")
				funcs.Add(m[2..])
		return funcs.Sort!()
		}

	crossTabFormat: InputFormat
		{
		New(query, selectFields)
			{
			super(Object('Crosstab', query,
				rows: .fieldname(_report.Params.Rows, selectFields),
				cols: .fieldname(_report.Params.Columns, selectFields),
				value: .fieldname(_report.Params.Value, selectFields),
				func: _report.Params.Function,
				sortcols:, :selectFields))
			}
		fieldname(prompt, selectFields)
			{
			return prompt is '' ? '' : selectFields.PromptToField(prompt)
			}
		}
	}
