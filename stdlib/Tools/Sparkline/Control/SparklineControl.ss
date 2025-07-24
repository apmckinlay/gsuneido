// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
/*
	SparklineControl - Written in February 2007 by Mauro Giubileo
	-------------------------------------------------------------

	Example Usage:
	Window(#(Sparkline (10, 8, 12, 2, 4, 7, 9, 18, 12, 13, 16, 23, 23)
		normalRange: (10, 15)) )
*/
WndProc
	{
	New(data = #(), .valField = '', .width = 128, .height = 30, .inside = 8, .thick = 1,
		.pointLineRatio = 5, .rectangle = true, .middleLine = false, .allPoints = false,
		.firstPoint = false, .lastPoint = true, .lines = true, .minPoint = false,
		.maxPoint = false, .normalRange = false, .normalRangeColor = 0xEEEEEE,
		.borderOnPoints = false, .circlePoints = false)
		{
		width = ScaleWithDpiFactor(width)
		height = ScaleWithDpiFactor(height)
		.Xmin = width
		.Ymin = height
		.Top = height / 2 + 5 /* = align with normal size text */
		.Data = .GetData(data, .valField)
		.CreateWindow('static', '', WS.VISIBLE)
		.SubClass()

		.paintGraph = SparklinePaintWithDC(data, .inside, .thick, .rectangle,
			.middleLine, .allPoints, .firstPoint, .lastPoint, .lines, .minPoint,
			.maxPoint, .normalRange, .normalRangeColor, .borderOnPoints,
			.circlePoints, .pointLineRatio)
		}

	GetData(data, valField)
		{
		if String?(data) and data isnt '' and String?(valField) and valField isnt ''
			return QueryList(data, valField).Map!(Number)
		return data
		}

	PAINT()
		{
		dc = BeginPaint(.Hwnd, ps = Object())
		GetClientRect(.Hwnd, r = Object())
		.paintGraph.Draw(r.left, r.top, r.right, r.bottom, :dc)
		EndPaint(.Hwnd, ps)
		return 0
		}

	DESTROY()
		{
		.paintGraph.Destroy()
		return 0
		}
	}
