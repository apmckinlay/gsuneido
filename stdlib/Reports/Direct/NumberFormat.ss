// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// converts "" to 0
OptionalNumberFormat
	{
	New(data = false, mask = false, width = false, w = false, font = false,
		color = false, justify = 'right', access = false)
		{
		super(data is false ? false : .Convert(data), :mask, :width, :w,
			:font, :color, :justify, :access)
		}
	Print(x, y, w, h, data = 0)
		{
		if Object?(data)
			data = .Data
		super.Print(x, y, w, h, .Convert(data))
		}
	ExportCSV(data = 0)
		{
		if Object?(data)
			data = .Data
		return super.ExportCSV(.Convert(data))
		}
	}