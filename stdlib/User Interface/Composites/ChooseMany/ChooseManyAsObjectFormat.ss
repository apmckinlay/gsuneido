// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New(.idField, .displayField, .query, .data = false, w = false, width = false,
		justify = "left", font = false,	.delimiter = ', ')
		{
		super(false, w, width, justify, font)
		}
	GetSize(data = "")
		{
		return super.GetSize(.getData(data))
		}
	getData(data, rec = false)
		{
		if .data isnt false
			data = .data
		if data isnt ""
			data = .BuildDesc(data, rec)
		return data
		}
	Print(x, y, w, h, data = "", rec = false)
		{
		super.Print(x, y, w, h, .getData(data, rec))
		}
	DataToString(data, rec)
		{
		return .BuildDesc(data, rec)
		}
	BuildDesc(data, rec = false)
		{
		if not Object?(data)
			return ''
		if .idField is .displayField
			return data.Join(.delimiter)

		descs = Object()
		for id in data.Copy()
			if false isnt rec = Query1(.query $ ' where ' $ .idField $
				' is ' $ Display(id))
				descs.Add(rec[.displayField])
		setval = descs.Join(.delimiter)
		return setval
		}
	}