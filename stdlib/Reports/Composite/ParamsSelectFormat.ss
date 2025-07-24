// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
HorzFormat
	{
	New(data, name, format = false)
		{
		super(@.layout(data, name, format))
		.data = data
		}
	layout(data, name, format)
		{
		items = Object()
		format = .createFormat(data, name, format)
		if data.operation.Has?("range")
			{
			format2 = format.Copy()
			format2.data = data.value2
			text = data.operation is 'not in range' ? 'Excluding From' : 'From'
			items.Add(Object('Text', TranslateLanguage(text) $ ' '))
			items.Add(format)
			items.Add(Object('Text', ' ' $ TranslateLanguage('To') $ ' '))
			items.Add(format2)
			}
		else if data.operation.Has?("in list")
			items.Add(Object('Text', TranslateLanguage(data.operation) $
				' (' $ ParamsChooseListControl.DisplayValues(data.value,
					name).Join(', ') $ ')'))
		else if (data.operation is "empty" or data.operation is "not empty")
			items.Add(Object('Text' TranslateLanguage(data.operation)))
		else
			{
			if data.operation is "equals"
				{
				if data.value is ""
					format = #(Text "\"\"")
				}
			else
				items.Add(Object('Text', TranslateLanguage(data.operation $ " ")))
			items.Add(format)
			}
		items.Add('Hskip')
		return items
		}

	createFormat(data, name, format)
		{
		dd = format is false
				? Datadict(name)
				// format will not be false for formulas from reporter reports
				// being printed in the params
				: format.Compile()
		format = dd.Format.Copy()
		format.Delete('width')
		if String?(format[0]) and format[0].Has?('Wrap')
			format[0] = 'Text'
		if format[0] is 'CheckMark'
			format[0] = 'Boolean'

		map = GetContributions('ParamSelectFormatMap')
		for mem in map.Members()
			if format[0] is mem
				format = map[mem]
		format.data = dd.Encode(data.value,
			paramSelect: data.GetDefault('operation', false))
		return format
		}

	ExportCSV(data = '')
		{
		return super.ExportCSV(data).Tr('\n', '')
		}
	}
