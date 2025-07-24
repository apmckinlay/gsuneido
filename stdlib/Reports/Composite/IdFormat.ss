// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New(.data = false, w = false, width = false, justify = "left",
		font = false, .query = "", .numField = "", .nameField = false,
		.abbrevField = false, .abbrev = false, access = false)
		{
		super(false, w, width, justify, font, :access)
		}
	GetSize(data = "")
		{
		return super.GetSize(.getPrintData(data))
		}
	Print(x, y, w, h, data = "")
		{
		super.Print(x, y, w, h, .getPrintData(data))
		}
	ExportCSV(data = '')
		{
		return .CSVExportString(.getPrintData(data))
		}
	getPrintData(data)
		{
		if .data isnt false
			data = .data
		if data isnt ""
			data = .GetIdName(data)
		return data
		}
	DataToString(data, rec /*unused*/)
		{
		return .GetIdName(data)
		}
	GetIdName(data)
		{
		query = QueryAddWhere(.query,
			' where ' $ .numField $ ' is ' $ Display(data))
		if false isnt cache = Suneido.GetDefault('ReportQueryCache', false)
			rec = cache.Get(query)
		else
			rec = Query1(query)
		if rec is false
			return String?(data) ? data : '???'
		if .abbrev is true and ('' isnt abbrev = .getField(rec, abbrev?:))
			return abbrev
		return .getField(rec)
		}

	getField(rec, abbrev? = false)
		{
		field = abbrev? ? .abbrevField : .nameField
		field = field is false
			? .numField.Replace('_num$', abbrev? ? '_abbrev' : '_name')
			: field
		return rec[field]
		}
	}