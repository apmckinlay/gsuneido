// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New(.columns, .data = false, w = false, width = false, justify = "left", font = false)
		{
		super(false, w, width, justify, font)
		}
	GetSize(data = "")
		{
		return super.GetSize(.getData(data))
		}
	getData(data)
		{
		if .data isnt false
			data = .data
		if Object?(data)
			data = .BuildDesc(data)
		return data
		}
	Print(x, y, w, h, data = "")
		{
		super.Print(x, y, w, h, .getData(data))
		}
	DataToString(data, rec /*unused*/)
		{
		return .BuildDesc(data)
		}
	BuildDesc(data)
		{
		setval = Object()
		for rec in data
			for col in .columns
				if rec.Member?(col)
					setval.Add(Prompt(col) $ ': ' $ rec[col])
		return setval.Join(", ")
		}
	}
