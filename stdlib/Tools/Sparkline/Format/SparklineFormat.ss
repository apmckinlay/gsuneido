// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
/*
	SparklineFormat - Written in February 2007 by Mauro Giubileo
	------------------------------------------------------------

	Example Usage:
	Params( #(Sparkline (10, 8, 12, 2, 4, 7, 9, 18, 12, 13, 16, 23, 23)
		normalRange: (10, 15)) )
*/
Format
	{
	New(data = #(), .valField = '', .width = 2000, .height = 500,
		.inside = 100, .thick = 14, .pointLineRatio = 4,
		.rectangle = true, .middleLine = false, .allPoints = false,
		.firstPoint = false, .lastPoint = true, .lines = true,
		.minPoint = false, .maxPoint = false, .normalRange = false,
		.normalRangeColor = 0xEEEEEE, .borderOnPoints = false, .circlePoints = false)
		{
		.Data = data
		if (String?(data) and data isnt '' and String?(valField) and valField isnt '')
			{
			.Data = Object()
			QueryApply(data)
				{|rec|
				.Data.Add(Number(rec[valField]))
				}
			}
		}

	GetSize(unused = '')
		{
		return Object(w: .width, h: .height, d: .height / 3 /*= offset factor*/)
		}

	Print(x, y, w, h, data = '')
		{
		if .Data isnt false
			data = .Data

		SparklinePaintWithReport(data, .inside, .thick, .rectangle,
			.middleLine, .allPoints, .firstPoint, .lastPoint, .lines, .minPoint,
			.maxPoint, .normalRange, .normalRangeColor, .borderOnPoints,
			.circlePoints, .pointLineRatio).
			Draw(x, y, w, h)
		}
	}
