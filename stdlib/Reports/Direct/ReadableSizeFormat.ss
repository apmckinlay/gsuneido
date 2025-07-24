// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New(data = false, width = false, w = false,
		.font = false, justify = 'right', color = false, export = true,
		access = false)
		{
		super(data is false ? false :
			(data is "" ? data : .Convert(data)),
			:width, :w, :font, :justify, :color, :export, :access)
		}
	Convert(data)
		{
		try
			return ReadableSize(Number(data))
		catch
			return data
		}
	GetSize(data = 0)
		{
		if Object?(data)
			data = .Data
		return super.GetSize(.Convert(data))
		}
	Print(x, y, w, h, data = 0)
		{
		if Object?(data)
			data = .Data
		if data isnt ''
			data = .Convert(data)
		super.Print(x, y, w, h, data)
		}
	ExportCSV(data = 0)
		{
		data = .getData(data)
		if not String?(data)
			data = Display(data)
		return data.Tr(',')
		}
	getData(data)
		{
		if .Data isnt false
			return .Data
		if Object?(data)
			data = .Data
		if data isnt ''
			data = .Convert(data)
		return data
		}
	}
