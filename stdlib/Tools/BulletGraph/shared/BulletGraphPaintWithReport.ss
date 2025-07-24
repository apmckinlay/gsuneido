// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
BulletGraphPaint
	{
	oldFont: false
	PreDraw(dc, .sz)
		{
		if dc isnt false
			_report.SetDC(dc)
		.font = Object(name: 'Arial', size: sz / 15) /*= font size*/
		.oldFont = _report.SelectFont(.font)
		}

	DrawVerticalMarks(text, lineOb, maxDigits)
		{
		textWidth = maxDigits * .sz
		fontX = lineOb.w - textWidth
		_report.AddLine(lineOb.x, lineOb.y, lineOb.w, lineOb.y, thick: .5)

		Format.DoWithFont(false)
			{ |font|
			_report.AddText(text, fontX - lineOb.length,
				lineOb.y - .sz / 2, textWidth, lineOb.h, font, justify: 'right')
			}
		}

	DrawHorizontalMarks(text, lineOb, fontY)
		{
		_report.AddLine(lineOb.x, lineOb.y, lineOb.x, lineOb.h, thick: .5)
		Format.DoWithFont(false)
			{ |font|
			_report.AddText(text,
				lineOb.x - .sz * (text.Size() - 1) / 2, fontY, lineOb.w, lineOb.h, font)
			}
		}

	DrawEnclosingRectangle(.x, .y, .w, .h)
		{
		_report.AddRect(.x, .y, .w, .h, thick: .5)
		}

	DrawBar(x1, y1, x2, y2, color)
		{
		_report.AddRect(x1, y1, x2 - x1, y2 - y1, thick: 0, fillColor: color)
		}

	DrawVerticalBars(lineOb)
		{
		_report.AddLine(lineOb.x, lineOb.y, lineOb.w, lineOb.y,
			Max(0.1, .w * 0.001) /*= thick */)
		}

	DrawHorizontalBars(lineOb)
		{
		_report.AddLine(lineOb.x, lineOb.y, lineOb.x, lineOb.h,
			Max(0.1 : (.w * 0.0005)) /*= thick */)
		}

	PostDraw()
		{
		if .oldFont is false
			return
		_report.SelectFont(.oldFont)
		.oldFont = false
		}
	}