// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New(.data = false, font = false, justify = 'right')
		{
		super(false, width: 5, :justify, :font)
		// width of 5 to fit heading of 'Time'
		}
	WidthChar: '9'
	Print(x, y, w, h, data = "")
		{
		if .data isnt false
			data = .data
		super.Print(x, y, w, h, Number?(data) ? data.Pad(3, '0') : data)
		}
	}