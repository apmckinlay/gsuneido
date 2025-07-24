// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New(.data = "", w = false, width = false, justify = "left",
		font = false, yesno = false, .strings = #("true", "false"),
		color = false, .showEmpty? = false)
		{
		super(false, w, width, justify, font, color)
		if yesno
			.strings = #("yes", "no")
		}
	GetSize(data = "")
		{
		data = .tostr(data, .strings)
		return super.GetSize(data)
		}
	Print(x, y, w, h, data = "")
		{
		data = .tostr(data, .strings)
		super.Print(x, y, w, h, data)
		}
	ExportCSV(data = "")
		{
		data = .tostr(data, #("Yes", "No"))
		return .CSVExportString(data)
		}
	tostr(data, strings)
		{
		if .data isnt ""
			data = .data
		if .showEmpty? and data is ''
			return ''
		data = BooleanOrEmpty?(data) ? strings[data is true ? 0 : 1] : data
		return TranslateLanguage(data)
		}
	}
