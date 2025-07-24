// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
/*
	BulletGraphControl - Written in February 2007 by Mauro Giubileo
	---------------------------------------------------------------

	Example Usage:
	Window( #(BulletGraph 24, satisfactory: 20, good: 25, target: 27, range: (0,30)) )
*/
WndProc
	{
	New(data = false, .satisfactory = 0, .good = 0, .target = 0, .range = #(0, 100),
		.color = 0x506363, width = 128, height = 32, .rectangle = true,
		.outside = 5, .vertical = false, .axis = false, .axisDensity = 5)
		{
		if .vertical and width is 128 and height is 32 /*= default vertical */
			{ // swap w and h
			temp = width
			width = height
			height = temp
			}
		width = ScaleWithDpiFactor(width)
		height = ScaleWithDpiFactor(height)

		.Xmin = width
		.Ymin = height
		.Data = data
		.width = width
		.height = height

		.CreateWindow("static", "", WS.VISIBLE)
		.SubClass()

		.paintGraph = BulletGraphPaintWithDC(.Data,
			.satisfactory, .good, .target, .range, .color,
			.rectangle, .outside, .vertical, .axis, .axisDensity)
		}

	PAINT()
		{
		dc = BeginPaint(.Hwnd, ps = Object())
		GetClientRect(.Hwnd, r = Object())
		.paintGraph.Draw(r.left, r.top, r.right, r.bottom, :dc)
		EndPaint(.Hwnd, ps)
		return 0
		}

	SetData(value)
		{
		.paintGraph.SetData(.Data = value)
		InvalidateRect(.Hwnd, NULL, true)
		UpdateWindow(.Hwnd)
		}

	DESTROY()
		{
		.paintGraph.Destroy()
		return 0
		}
	}
