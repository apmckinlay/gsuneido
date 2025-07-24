// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
/*
	BulletGraphFormat - Written in February 2007 by Mauro Giubileo
	--------------------------------------------------------------

	Example Usage:
	Params( #(BulletGraph 24, satisfactory: 20, good: 25, target: 27, range: (0,30)) )
*/
Format
	{
	New(data = false, .satisfactory = 0, .good = 0, .target = 0, .range = #(0, 100),
		.color = 0x506363, width = 1280, height = 320, .rectangle = true,
		.outside = 50, .vertical = false, .axis = false, .axisDensity = 5)
		{
		if (.vertical and width is 1280 and height is 320) /*= default vertical*/
			{ // swap w and h
			temp = width
			width = height
			height = temp
			}
		.Data = data
		.width = width
		.height = height
		}

	GetSize(data /*unused*/ = "")
		{
		//Post: returns width, height and descent of the format as Object(w:, h:, d:)
		return Object(w: .width, h: .height, d: 0)
		}

	Print(x, y, w, h, data = "")
		{
		//Post: prints the format at position x,y with specified width and height
		if (.Data isnt false)
			data = .Data

		BulletGraphPaintWithReport(data, .satisfactory, .good, .target, .range, .color,
			.rectangle, .outside, .vertical, .axis, .axisDensity).
			Draw(x, y, w, h)
		}
	}