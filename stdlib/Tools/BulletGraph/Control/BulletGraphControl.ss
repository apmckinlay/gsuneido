// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
/*
	BulletGraphControl - Written in February 2007 by Mauro Giubileo
	---------------------------------------------------------------

	Example Usage:
	Window( #(BulletGraph 24, satisfactory: 20, good: 25, target: 27, range: (0,30)) )
*/
WndProc
	{
	New(.data, .satisfactory = 0, .good = 0, .target = false, .range = #(0, 100),
		.color = 0x506363, width = 128, height = 32, .rectangle = true,
		.outside = 5, .vertical = false, .axis = false, .axisDensity = 5,
		.axisFormat = false, .selectedColor = false)
		{
		if .vertical and width is 128 and height is 32 /*= default vertical sizing*/
			{ // swap w and h
			temp = width
			width = height
			height = temp
			}
		.validateData()
		.CreateWindow('static', '', WS.VISIBLE | SS.NOTIFY)
		.SubClass()
		.setup()

		.Xmin = .width = ScaleWithDpiFactor(width)
		.Ymin = .height = ScaleWithDpiFactor(height)
		if .vertical
			{
			.drawMethod = .drawVertical
			.Xmin += .axisDetails.total.x
			.Ymin += (.axisDetails.text.y / 2).RoundUp(0)
			}
		else
			{
			.drawMethod = .drawHorizontal
			.Ymin += .axisDetails.total.y
			}
		}

	validateData()
		{
		Assert(.data < .range[1])
		Assert(.satisfactory >= .range[0] and .satisfactory <= .range[1])
		Assert(.good >= .range[0] and .good <= .range[1])
		Assert(.range[0] < .range[1])
		if .target isnt false
			Assert(.target >= .range[0] and .target <= .range[1])
		}

	setup()
		{
		.initResources()
		if .axis is true
			.WithSelectObject(.rsrc.font)
				{|hdc|
				maxValue = Max(Abs(.range[0]), Abs(.range[1]))
				text = .axisFormat isnt false
					? maxValue.Format(.axisFormat)
					: String(maxValue)
				.axisDetails = Object()
				GetTextExtentPoint32(hdc, text, text.Size() + 1,
					.axisDetails.text = Object())
				axisSize = ScaleWithDpiFactor(10 /*= axis size*/)
				.axisDetails.mark = Object(x: axisSize, y: axisSize)
				.axisDetails.total = Object(
					x: .axisDetails.SumWith({ it.x }),
					y: .axisDetails.SumWith({ it.y }))
				}
		else
			.axisDetails = Object().Set_default(#(x: 0, y: 0))
		.setupColors()
		.normalization()
		}

	rsrc: #()
	initResources()
		{
		.rsrc = Object()
		lf = Suneido.logfont.Copy()
		lf.lfHeight += 2 // Ensure font is smaller than the user font
		.rsrc.font = CreateFontIndirect(lf)
		.rsrc.axisPen = ExtCreatePen(PS.GEOMETRIC | PS.SOLID | PS.ENDCAP_ROUND, 1,
			Object(lbstyle: BS.SOLID, lbColor: 0x000000, lbHatch: 0), 0, 0)
		.rsrc.rectPen = CreatePen(PS.SOLID, 1, 0)
		.rsrc.clearPen = CreatePen(PS.SOLID, 1, GetSysColor(COLOR.TRIDFACE))
		}

	setupColors()
		{
		// Setup Colors
		// Extract single RGB components from color with a mask:
		r = (.color & 0xFF0000) / 0x10000 	/*= RGB: red value */
		g = (.color & 0x00FF00) / 0x100		/*= RGB: green value */
		b = (.color & 0x0000FF)				/*= RGB: blue value */

		.badColor = .bgr(r, g, b, 80 /*= fading rate*/)
		.satisfactoryColor = .bgr(r, g, b, 120 /*= fading rate*/)
		.goodColor = .bgr(r, g, b, 150 /*= fading rate*/)
		.valueColor = .bgr(r, g, b, 0)
		}

	bgr(r, g, b, fadingRate)
		{
		shade = Max(r, g, b)
		incr = (shade * fadingRate).PercentToDecimal().Int()
		r = Min(255 /*= max color*/, r + incr)
		g = Min(255 /*= max color*/, g + incr)
		b = Min(255 /*= max color*/, b + incr)
		return RGB(r, g, b)
		}

	normalization()
		{
		if .range.Readonly?()
			.range = .range.Copy()
		.range[1] -= .range[0]
		.data -= .range[0]
		.satisfactory -= .range[0]
		.good -= .range[0]
		if .target isnt false
			.target -= .range[0]
		}

	PAINT()
		{
		.dc = BeginPaint(.Hwnd, ps = Object())
		GetClientRect(.Hwnd, r = Object())
		.draw(r.left, r.top, r.right, r.bottom)
		EndPaint(.Hwnd, ps)
		return 0
		}

	axisSpacing: 2
	actualBarRatio: 3
	targetIndicatorEnd: 3
	targetIndicatorRatio: 4
	draw(.x, .y, .w, .h)
		{
		// Add outside spacing:
		.x += .outside
		.y += .outside
		.w -= .outside * 2
		.h -= .outside * 2
		(.drawMethod)()
		}

	// ========================== Draw Vertical ==========================
	drawVertical()
		{
		WithHdcSettings(.dc, [.rsrc.clearPen, brush: GetSysColor(COLOR.TRIDFACE)])
			{
			Rectangle(.dc, .x, .y - .axisDetails.text.y / 2, .w, .h)
			}
		.w -= .axisDetails.total.x
		.h -= .axisDetails.text.y / 2
		// Check if outside area is too big in respect to width and height:
		if .w > 0 and .h > 0
			.drawGraphs(
				{ .drawRect(.axisDetails.total.x) },
				{ .drawAxis(TA.RIGHT, .drawVerticalMarks) },
				.drawVerticalBars)
		}

	drawVerticalMarks()
		{
		incr = .h / (.axisDensity - 1)
		x = .x + .axisDetails.text.x
		lineX1 = x + .axisSpacing
		lineX2 = lineX1 + .axisDetails.mark.x
		for i in .. .axisDensity
			{
			// Create small reference segments for the axis:
			MoveTo(.dc, lineX1, y = (.y + .h - i * incr) - 1)
			LineTo(.dc, lineX2, y)

			// Print text beside the segments:
			value = (.range[1] * i * incr / .h + .range[0]).Round(0)
			text = .axisFormat is false
				? String(value)
				: value.Format(.axisFormat)
			TextOut(.dc, x, y - .axisDetails.text.y / 2, text, text.Size())
			}
		}

	drawVerticalBars()
		{
		// Draw good values bar:
		x = .x + .axisDetails.total.x
		.drawBar(x, .y + .h, x + .w, .y, .goodColor)

		// Draw satisfactory values bar:
		.drawBar(x, .y + .h, x + .w, .y + .h - .h * .good / .range[1], .satisfactoryColor)

		// Draw bad values bar:
		.drawBar(x, .y + .h, x + .w, .y + .h - .h * .satisfactory / .range[1], .badColor)

		// Draw actual value bar:
		if .data isnt false
			.drawBar(x + .w / .actualBarRatio, .y + .h,
				x + .w / .actualBarRatio * 2,
				.y + .h - .h * .data / .range[1], .barColor())

		// Draw target value indicator:
		if .target isnt false
			.drawVerticalTarget()
		}

	drawVerticalTarget()
		{
		x = .x + .axisDetails.total.x + .w / .targetIndicatorRatio
		y = .y + .h - .h * .target / .range[1]
		w = .x + .axisDetails.total.x + .w / .targetIndicatorRatio * .targetIndicatorEnd
		.drawBar(x, y, w, y - .targetSize, CLR.BLACK)
		}

	// ========================= Draw Horizontal =========================
	drawHorizontal()
		{
		WithHdcSettings(.dc, [.rsrc.clearPen, brush: GetSysColor(COLOR.TRIDFACE)])
			{
			Rectangle(.dc, .x - .axisDetails.text.x / 2, .y, .w + .axisDetails.text.x, .h)
			}
		.h -= .axisDetails.total.y
		// Check if outside area is too big in respect to width and height:
		if .w > 0 and .h > 0
			.drawGraphs(
				.drawRect,
				{ .drawAxis(TA.CENTER, { .drawHorizontalMarks() }) },
				.drawHorizontalBars)
		}

	drawHorizontalMarks()
		{
		incr = .w / (.axisDensity - 1)
		lineH1 = .y + .h
		lineH2 = lineH1 + .axisDetails.mark.y - .axisSpacing
		for i in .. .axisDensity
			{
			// Create small reference segments for the axis:
			MoveTo(.dc, x = (.x + i * incr) - 1, lineH1)
			LineTo(.dc, x, lineH2)

			// Print text beside the segments:
			value = (.range[1] * i * incr / .w + .range[0]).Round(1)
			text = .axisFormat is false
				? String(value)
				: value.Format(.axisFormat)
			TextOut(.dc, x, lineH2, text, text.Size())
			}
		}

	drawHorizontalBars()
		{
		// Draw good values bar:
		.drawBar(.x, .y, .x + .w, .y + .h, .goodColor)

		// Draw satisfactory values bar:
		.drawBar(.x, .y, .x + .w * .good / .range[1], .y + .h, .satisfactoryColor)

		// Draw bad values bar:
		.drawBar(.x, .y, .x + .w * .satisfactory / .range[1], .y + .h, .badColor)

		// Draw actual value bar:
		if .data isnt false
			.drawBar(.x, .y + .h / .actualBarRatio, .x + .w * .data / .range[1],
				.y + .h / .actualBarRatio * 2, .barColor())

		// Draw target value indicator:
		if .target isnt false
			.drawHorizontalTarget()
		}

	drawHorizontalTarget()
		{
		x = .x + .w * .target / .range[1]
		y = .y + .h / .targetIndicatorRatio
		h = .y + .h / .targetIndicatorRatio * .targetIndicatorEnd
		.drawBar(x, y, x + .targetSize, h, CLR.BLACK)
		}

	// ========================== Draw General ===========================
	drawGraphs(rectDraw, axisDraw, graphDraw)
		{
		if .rectangle
			(rectDraw)()
		if .axis
			(axisDraw)()
		(graphDraw)()
		}

	drawRect(xOffset = 0)
		{
		DoWithHdcObjects(.dc, [.rsrc.rectPen])
			{
			x = .x + xOffset
			Rectangle(.dc, x - 1, .y - 1, x + .w, .y + .h)
			}
		}

	drawAxis(align, block)
		{
		WithHdcSettings(.dc,
			[.rsrc.axisPen, .rsrc.font, SetTextAlign: align, SetBkMode: TRANSPARENT],
			block)
		}

	drawBar(x1, y1, x2, y2, color)
		{
		pen = CreatePen(PS.NULL, 1, color)
		WithHdcSettings(.dc, [pen, brush: color], { Rectangle(.dc, x1, y1, x2, y2) })
		DeleteObject(pen)
		}

	barColor()
		{
		return .selectedColor isnt false and .selected
			? .selectedColor
			: .valueColor
		}

	getter_targetSize()
		{
		return .targetSize = ScaleWithDpiFactor(3 /*=size*/)
		}

	// ======================== User Interaction =========================
	SetEnabled(enabled /*unused*/)
		{
		// Suppress to allow mouse calls without relying on setting Edit mode
		}

	MOUSEMOVE()
		{
		.Send('BulletGraph_Hover')
		return false
		}

	LBUTTONUP()
		{
		.Send('BulletGraph_Click')
		return false
		}

	selected: false
	Selected(selected)
		{
		if .selectedColor is false
			return
		repaint? = selected isnt .selected
		.selected = selected
		if repaint?
			.Repaint()
		}

	Destroy()
		{
		.rsrc.Each(DeleteObject)
		.rsrc = #()
		super.Destroy()
		}
	}
