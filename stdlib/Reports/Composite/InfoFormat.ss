// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New(@args)
		{
		super(@args)
		if .Data isnt false
			.Data = StripInfoLabel(.Data)
		}

	Print(x, y, w, h, data = "")
		{
		if .Data is false and String?(data)
			data = StripInfoLabel(data)
		super.Print(x, y, w, h, data)
		}
	}