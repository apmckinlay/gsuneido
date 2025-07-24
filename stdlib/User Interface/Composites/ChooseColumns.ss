// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(columns, logNoPromptAsError = false)
		{
		.promptMap = Datadict.GetPromptMap(columns, logAsError: logNoPromptAsError)
		}

	AvailableList()
		{
		defaultColumnObs = Object()
		for col in .promptMap.Members()
			defaultColumnObs.Add(
				Object(column: .promptMap[col],
					col_width: .getDefaultWidth(col)))
		return defaultColumnObs
		}

	// expecting object like #(#(column: field_name, col_width: NUMBER))
	// returns   object like #(#(field_name, width: NUMBER), field_name2)
	GetSaveData(listOb)
		{
		data = Object()
		for col in listOb
			{
			fieldName = .getFieldNameByPrompt(col.column)
			if col.col_width is .getDefaultWidth(fieldName)
				data.Add(fieldName)
			else
				data.Add(Object(fieldName, width: col.col_width))
			}
		return data
		}

	getFieldNameByPrompt(prompt)
		{
		field = .promptMap.Find(prompt)
		return field isnt false ? field : ''
		}

	defaultWidth: 10
	getDefaultWidth(col, _report = false)
		{
		ddFormat = Datadict(col).Format
		if ddFormat.Member?('width')
			return ddFormat.width

		fmt = Global(ddFormat[0].RemoveSuffix('Format') $ 'Format')
		if fmt.Method?('GetDefaultWidth')
			{
			if report is false
				_report = ReportInstance()
			if false isnt w = _report.Construct(ddFormat).GetDefaultWidth()
				return w
			}

		return .defaultWidth
		}

	// expecting object like #(#(field_name, width: NUMBER), field_name2)
	// returns   object like #(#(column: field_name, col_width: NUMBER))
	SetSaveList(data)
		{
		listOb = Object()
		for col in data
			{
			if String?(col)
				col = Object(col)
			listOb.Add(.convertSavedCol(col))
			}
		return listOb
		}

	convertSavedCol(colOb)
		{
		fieldName = colOb[0]
		width = colOb.GetDefault('width', .getDefaultWidth(fieldName))
		return Object(
			column: .promptMap.GetDefault(fieldName, fieldName),
			col_width: width)
		}

	ValidList(listOb)
		{
		for colOb in listOb
			if not Number?(colOb.col_width) or colOb.col_width <= 0
				return colOb.column $ ' has an invalid Width'
		return ''
		}

	FindColumn(cols, col)
		{
		return cols.FindIf({ .GetFieldName(it) is col })
		}

	GetFieldName(col)
		{
		if not Object?(col)
			return col
		if col.Member?('field')
			return col.field
		return col[0]
		}
	}