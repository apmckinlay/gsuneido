// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
CheckMarkBaseFormat
	{
	New(.data = '', .w = false, width = 4)
		{
		super(width)
		}

	Print(x, y, w, h, data = "")
		{
		if .data isnt ""
			data = .data
		if data is ''
			return
		else if not Boolean?(data)
			{
			.PrintInvalidData(x, y, w, h, data)
			return
			}

		super.PrintWithBox(x, y, w, h, data)
		}

	ExportCSV(data = "")
		{
		if .data isnt ""
			data = .data
		return super.ExportCSV(data)
		}
	}