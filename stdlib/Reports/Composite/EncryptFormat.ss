// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New(.data = false, w = false, width = false, justify = "left", font = false,
		color = false)
		{
		super(:w, :width, :justify, :font, :color)
		}
	GetSize(data = "")
		{
		if .data isnt false
			data = .data
		return super.GetSize(.Convert(data))
		}
	Print(x, y, w, h, data = "")
		{
		if .data isnt false
			data = .data
		super.Print(x, y, w, h, .Convert(data))
		}
	ExportCSV(data = '')
		{
		if .data isnt false
			data = .data
		super.ExportCSV(.Convert(data))
		}
	Convert(data)
		{
		if not String?(data)
			data = String(data)
		return data.Xor(EncryptControlKey())
		}
	DataToString(data, rec /*unused*/)
		{
		return .Convert(data)
		}
	}
