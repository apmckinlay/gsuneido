// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
SparklinePaint
	{
	New(@args)
		{
		super(@args)
		.pens = Object()
		}

	Draw(x, y, w, h, .dc)
		{
		super.Draw(x, y, w, h)
		}

	middleLinePen: 	false
	linePen:		false
	PreDraw(.circlePoints, .borderOnPoints, .thick)
		{
		if .middleLinePen is false
			.middleLinePen = ExtCreatePen(
				PS.GEOMETRIC|PS.SOLID|PS.ENDCAP_FLAT, .thick / 2,
				Object(lbstyle: BS.SOLID, lbColor: 0x4444BB, lbHatch: 0),
				0, 0)
		if .linePen is false
			.linePen = ExtCreatePen(
				PS.GEOMETRIC|PS.SOLID|PS.ENDCAP_ROUND, .thick,
				Object(lbstyle: BS.SOLID, lbColor: 0, lbHatch: 0),
				0, 0)
		}

	DrawRectangle(x1, y1, x2, y2)
		{
		DoWithHdcObjects(.dc, [.pen(PS.SOLID, 1, 0)])
			{
			Rectangle(.dc, x1, y1, x2, y2)
			}
		}

	pen(style, thick, color)
		{
		key = Object(String(style), String(thick), String(color)).Join('_')
		return .pens.GetInit(key, { CreatePen(style, thick, color) })
		}

	DrawGraphs(.x, .y, .w, .h)
		{
		super.DrawGraphs(.x, .y, .w, .h)
		}

	DrawNormalBand(checkedNormalRange, normalRangeColor)
		{
		WithHdcSettings(.dc, [.pen(PS.NULL, .thick, 0), brush: normalRangeColor])
			{
			y1Offset = .h - checkedNormalRange[0]
			y2Offset = .h - checkedNormalRange[1]
			Rectangle(.dc, .x, .y + y1Offset, .x + .w, .y + y2Offset)
			}
		}

	DrawMiddleLine()
		{
		DoWithHdcObjects(.dc, [.middleLinePen])
			{
			MoveTo(.dc, .x, .y + .h / 2)
			LineTo(.dc, .x + .w, .y + .h / 2)
			}
		}

	DrawLines(x, y, data, offset)
		{
		DoWithHdcObjects(.dc, [.linePen])
			{
			MoveTo(.dc, x, y)
			for value in data[1..]
				LineTo(.dc, x += offset, .y + .h - .Normalize(value))
			}
		}

	DrawPoint(x, y, half1, half2, color)
		{
		hdcSettings = Object(brush: color) // CLR
		if .borderOnPoints is false
			hdcSettings.Add(.pen(PS.SOLID, 1, color))
		WithHdcSettings(.dc, hdcSettings)
			{
			x1 = x - half1
			y1 = y - half1
			x2 = x + half2
			y2 = y + half2
			if .circlePoints
				Ellipse(.dc, x1, y1, x2, y2)
			else
				Rectangle(.dc, x1, y1, x2, y2)
			}
		}

	Destroy()
		{
		.pens.Each(DeleteObject).Delete(all:)
		if .middleLinePen isnt false
			DeleteObject(.middleLinePen)
		if .linePen isnt false
			DeleteObject(.linePen)
		.middleLinePen = .linePen = false
		}
	}
