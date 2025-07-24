// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
VertFormat
	{
	Xstretch: 1
	New(center, .printParams? = false)
		{
		super(@.format(center))
		}
	params: #(WrapItems)
	format(center)
		{
		vert = Object(.buildCanvas(center), "Vskip")
		if .printParams? is true
			{
			.params  = PrintParams()
			if .params isnt #('WrapItems')
				vert.Add("Vskip", .params, "Vskip")
			}
		return vert
		}
	buildCanvas(item)
		{
		if item is #()
			return #(Vert)

		return Object("DrawCanvasRow", CanvasGroup(DrawControl.BuildItems(item)),
			canvasWidth: item.GetDefault(#canvasWidth, false))
		}
	// only print the params
	ExportCSV(data /*unused*/ = '')
		{
		return Opt(_report.Construct(.params).ExportCSV(), '\n,\n')
		}
	}
