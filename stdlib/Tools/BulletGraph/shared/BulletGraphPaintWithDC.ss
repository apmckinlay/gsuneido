// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
BulletGraphPaint
	{
	rsrc: #()
	PreDraw(.dc, .sz)
		{
		if .rsrc.NotEmpty?()
			return
		.rsrc = Object()
		// Create the font for text below the axis:
		lf = Object(
			lfFaceName: 'Arial',
			lfHeight: -sz * GetDeviceCaps(dc, GDC.LOGPIXELSY) / 72 /*= dpi*/,
			lfWeight: '',
			lfUnderline: false,
			lfEscapement: 0,
			lfOrientation: 0,
			lfCharSet: CHARSET[GetLanguage().charset])
		lf.lfHeight = lf.lfHeight.Round(0)
		.rsrc.font = CreateFontIndirect(lf)

		.rsrc.axisPen = ExtCreatePen(PS.GEOMETRIC | PS.SOLID | PS.ENDCAP_ROUND, 1,
			Object(lbstyle: BS.SOLID, lbColor: 0x000000, lbHatch: 0), 0, 0)

		.rsrc.rectPen = CreatePen(PS.SOLID, 1, 0)

		.rsrc.barsPen = ExtCreatePen(PS.GEOMETRIC | PS.SOLID | PS.ENDCAP_ROUND,
			Max(.Rect.w * 0.01, 2), /*= width */
			Object(lbstyle: BS.SOLID, lbColor: 0x000000, lbHatch: 0), 0, 0)

		.rsrc.clearPen = CreatePen(PS.SOLID, 1, GetSysColor(COLOR.TRIDFACE))
		}

	DrawGraphs(x, y, w, h)
		{
		WithHdcSettings(.dc, [.rsrc.clearPen, brush: GetSysColor(COLOR.TRIDFACE)])
			{ Rectangle(.dc, .Rect.x, .Rect.y, .Rect.w, .Rect.h) }
		super.DrawGraphs(x, y, w, h)
		}

	DrawAxis(align, block)
		{
		WithHdcSettings(.dc,
			[.rsrc.axisPen, .rsrc.font, SetTextAlign: align, SetBkMode: TRANSPARENT],
			block)
		}

	DrawVerticalMarks(text, lineOb, maxDigits /*unused*/)
		{
		lineOb.y -= 1
		// Create small reference segments for the axis:
		MoveTo(.dc, lineOb.x + 2, lineOb.y)
		LineTo(.dc, lineOb.w, lineOb.y)

		// Print text beside the segments:
		TextOut(.dc, lineOb.x, lineOb.y - .sz / 1.25, text, text.Size()) /*= offset */
		}

	DrawHorizontalMarks(text, lineOb, fontY)
		{
		lineOb.x -= 1
		// Create small reference segments for the axis:
		MoveTo(.dc, lineOb.x, lineOb.y)
		LineTo(.dc, lineOb.x, lineOb.h + 2)

		// Print text beside the segments:
		TextOut(.dc, lineOb.x, fontY, text, text.Size())
		}

	DrawEnclosingRectangle(x, y, w, h)
		{
		DoWithHdcObjects(.dc, [.rsrc.rectPen])
			{ Rectangle(.dc, x - 1, y - 1, x + w, y + h) }
		}

	DrawBar(x1, y1, x2, y2, color)
		{
		pen = CreatePen(PS.NULL, 1, color)
		WithHdcSettings(.dc, [pen, brush: color], { Rectangle(.dc, x1, y1, x2, y2) })
		DeleteObject(pen)
		}

	DrawVerticalBars(lineOb)
		{
		.drawBars(lineOb.x, lineOb.y, lineOb.w, lineOb.y)
		}

	drawBars(x1, y1, x2, y2)
		{
		DoWithHdcObjects(.dc, [.rsrc.barsPen])
			{
			MoveTo(.dc, x1, y1)
			LineTo(.dc, x2, y2)
			}
		}

	DrawHorizontalBars(lineOb)
		{
		.drawBars(lineOb.x, lineOb.y, lineOb.x, lineOb.h)
		}

	Destroy()
		{
		.rsrc.Each(DeleteObject)
		.rsrc = #()
		}
	}