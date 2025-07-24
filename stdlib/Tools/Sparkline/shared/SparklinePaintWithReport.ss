// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
SparklinePaint
	{
	PreDraw(.circlePoints, .borderOnPoints, .thick)
		{ }

	DrawRectangle(x1, y1, x2, y2)
		{
		_report.AddRect(x1, y1, x2 - x1, y2 - y1, thick: 1)
		}

	DrawGraphs(.x, .y, .w, .h)
		{
		super.DrawGraphs(.x, .y, .w, .h)
		}

	DrawNormalBand(checkedNormalRange, normalRangeColor)
		{
		bandY = .y + (.h - checkedNormalRange[1])
		bandH = checkedNormalRange[1] - checkedNormalRange[0]
		_report.AddRect(.x, bandY, .w, bandH, thick: 0, fillColor: normalRangeColor)
		}

	DrawMiddleLine()
		{
		lineY = .y + .h / 2
		_report.AddLine(.x, lineY, .x + .w, lineY, .thick, color: 0x4444BB)
		}

	DrawLines(x, y, data, offset)
		{
		for value in data[1..]
			{
			prevX = x
			prevY = y
			x += offset
			y = .y + .h - .Normalize(value)
			_report.AddLine(prevX, prevY, x, y, .thick)
			}
		}

	DrawPoint(x, y, half1, half2, color)
		{
		thick = .borderOnPoints ? 5 : 0 /*= offset*/
		if .circlePoints
			_report.AddCircle(x, y, half2, :thick, fillColor: color)
		else
			_report.AddRect(x - half1, y - half1, half1 + half2,
				half1 + half2, :thick, fillColor: color)
		}
	}
