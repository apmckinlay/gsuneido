// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New(data = false, .mask = false, width = false, w = false,
		.font = false, justify = 'right', color = false, export = true,
		access = false)
		{
		super(data is false ? false :
			(mask is false or data is "" ? data : .Convert(data, mask)),
			width: .getwidth(width, mask), :w,
			:font, :justify, :color, :export, :access)
		}
	WidthChar: '9'
	getwidth(width, mask)
		{
		if width is false and mask isnt false
			width = .EvalMask(mask).Size()
		return width
		}
	EvalMask(mask)
		{
		if String?(mask) and mask =~ "^-?[A-Z][a-zA-Z_0-9]*$"
			mask = mask.Extract('-?') $ Global(mask.Replace('-', ''))
		return mask
		}
	Convert(data, mask = false)
		{
		try
			data = (data is false) ? 0 : Number(data)
		catch
			data = String(data)
		mask = .EvalMask(mask)
		if IsInf?(data)
			data = ''
		if mask isnt false and Number?(data)
			data = data.Format(mask)
		return data
		}
	GetSize(data = 0)
		{
		if Object?(data)
			data = .Data
		return super.GetSize(.Convert(data, .mask))
		}
	Print(x, y, w, h, data = 0)
		{
		if Object?(data)
			data = .Data
		if .mask isnt false and data isnt ''
			data = .Convert(data, .mask)
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
		if .mask isnt false and data isnt ''
			data = .Convert(data, .mask)
		return data
		}
	}
