// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New(data = false, w = false, width = false, justify = "left",
		font = false, color = false, export = true, access = false)
		{
		super(data isnt false ? ScintillaRichStripHTML(data) : data, w,
			width, justify, font, color, export, access)
		}

	GetSize(data = "")
		{
		if data isnt ''
			data = ScintillaRichStripHTML(data)
		super.GetSize(data)
		}

	Print(x, y, w, h, data = "")
		{
		if data isnt ""
			data = ScintillaRichStripHTML(data)
		super.Print(x, y, w, h, data, debug:)
		}

	}