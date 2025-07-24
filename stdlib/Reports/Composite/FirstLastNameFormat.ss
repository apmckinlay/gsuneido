// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New(.data = false, w = false, width = false, justify = "left", font = false)
		{
		super(false, w, width, justify, font)
		}
	convert(name)
		{
		name = NameSplit(name, split_on: ',')
		return Join(' ', name.first, name.last)
		}
	GetSize(data = "")
		{
		if (.data isnt false)
			data = .data
		return super.GetSize(.convert(data))
		}
	Print(x, y, w, h, data = "")
		{
		if (.data isnt false)
			data = .data
		super.Print(x, y, w, h, .convert(data))
		}
	}