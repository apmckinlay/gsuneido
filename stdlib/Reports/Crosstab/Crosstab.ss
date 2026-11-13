// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(query, title = "", excludeFields = #(), columns = false)
		{
		if title is ""
			title = QueryGetTable(query, nothrow:)
		query = Query.StripSort(query)
		if columns is false
			columns = QueryColumns(query)
		excludeFields = ExcludeLowerFields(columns, excludeFields)
		selectFields = SelectFields(columns, excludeFields)
		prompts = selectFields.Prompts().Sort!()
		funcs = .crossTabFuncs()
		skipSave? = title is false
		return [#Params,
			[.crossTabFormat, :query, :selectFields]
			title: #CrossTable,
			NoPresets: skipSave?,
			NoSaveLoadParams: skipSave?,
			name: "CrossTable - " $ title,
			validField: #crosstable_valid,
			Params: [#Vert,
				[#Pair, #(Static, Rows),
					[#ChooseList, prompts, name: #Rows]],
				[#Pair, #(Static, Columns),
					[#ChooseList, prompts, name: #Columns]],
				[#Pair, #(Static, Value),
					[#ChooseList, prompts, name: #Value]],
				[#Pair, #(Static, Function),
					[#ChooseList, funcs, name: #Function, mandatory:]]
				xstretch: 0, ystretch: 0],
			printParams: #(Rows, Columns, Value, Function)]
		}

	crossTabFuncs()
		{
		funcs = Object()
		for m in Accumulator.Members()
			if m.Prefix?(#F_)
				funcs.Add(m[2..])
		return funcs.Sort!()
		}

	crossTabFormat: InputFormat
		{
		New(query, selectFields)
			{
			super([#Crosstab, query,
				rows: .fieldname(_report.Params.Rows, selectFields),
				cols: .fieldname(_report.Params.Columns, selectFields),
				value: .fieldname(_report.Params.Value, selectFields),
				func: _report.Params.Function,
				sortcols:, :selectFields])
			}

		fieldname(prompt, selectFields)
			{
			return prompt is "" ? "" : selectFields.PromptToField(prompt)
			}
		}
	}
