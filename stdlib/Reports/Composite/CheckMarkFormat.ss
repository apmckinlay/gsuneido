// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
CheckMarkBaseFormat
	{
	New(data = '', width = 4)
		{
		super(width)
		.data = data
		}

	Print(x, y, w, h, data = "")
		{
		if .data isnt ""
			data = .data
		if data is true
			super.Print(x, y, w, h)
		else if data isnt false and data isnt ""
			.PrintInvalidData(x, y, w, h, data)
		}
	ExportCSV(data = "")
		{
		if (.data isnt "")
			data = .data
		return super.ExportCSV(data)
		}
	Variable?()
		{ return false }
	}
