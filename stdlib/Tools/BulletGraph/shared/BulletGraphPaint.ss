// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
// A class to share functionality between BulletGraph Control and Format
// - by Mauro Giubileo - added PDF code in March 2007
class
	{
	New(data, .satisfactory, .good, .target, .range, .color,
		.rectangle, .outside, .vertical, .axis, .axisDensity)
		{
		.SetData(data)
		Assert(satisfactory >= range[0] satisfactory <= range[1])
		Assert(good >= range[0] and good <= range[1])
		Assert(target >= range[0] and target <= range[1])
		Assert(range[0] < range[1])

		// some normalizations:
		.normRange = Object(0, .range[1])
		.normRange[1] = .range[1] - .range[0]
		.satisfactory -= .range[0]
		.good -= .range[0]
		.target -= .range[0]
		}

	SetData(.data)
		{
		Assert(.data < .range[1])
		.data -= .range[0]
		}

	Draw(.x, .y, .w, .h, dc = false)
		{
		.Rect = Object(x: .x, y: .y, w: .w, h: .h)

		// add outside spacing:
		.x += .outside
		.y += .outside
		.w -= .outside * 2
		.h -= .outside * 2
		.sz = Max(.w, .h) / 24  /*= font size */
		.fontOffset = .sz / 2

		.PreDraw(dc, .sz)
		.DrawGraphs(.x, .y, .w, .h)
		.PostDraw()
		}

	PreDraw(dc /*unused*/, sz /*unused*/)
		{ }

	DrawGraphs(.x, .y, .w, .h)
		{
		if .axis
			.drawAxisMarks()
		// check if outside area is too big in respect to width and height:
		if .w <= 0 or .h <= 0
			return
		.setupColors()
		if .rectangle
			.DrawEnclosingRectangle(.x, .y, .w, .h)
		.drawBars()
		}

	drawAxisMarks()
		{
		if .vertical
			.drawVerticalAxis()
		else
			.drawHorizontalAxis()
		}

	segments: 12
	drawVerticalAxis()
		{
		.h -= .sz // this is for adding the text height
		incr = .h / (.axisDensity - 1)
		.axisSegmentSize = .w / .segments
		maxDigits = String(Max(Abs(.range[0]), Abs(.range[1]))).Size() + 1
		.x += maxDigits * .sz
		.y += .fontOffset
		.w -= maxDigits * .sz + .fontOffset + .axisSegmentSize
		.DrawAxis(TA.RIGHT, { .drawVerticalMarks(incr, maxDigits) })
		.x += .fontOffset + .axisSegmentSize
		}

	DrawAxis(align /*unused*/, block)
		{
		block()
		}

	drawVerticalMarks(incr, maxDigits)
		{
		for i in .. .axisDensity
			{
			text = Display(((.range[1] - .range[0]) * i * incr / .h + .range[0]).Round(1))
			length = .fontOffset + .axisSegmentSize
			lineOb = Object(
				x: .x + .fontOffset,
				y: .y + .h - i * incr,
				w: .x + length,
				h: .h,
				:length)
			.DrawVerticalMarks(text, lineOb, maxDigits)
			}
		}

	DrawVerticalMarks()
		{
		throw 'must be implemented by derived class'
		}

	drawHorizontalAxis()
		{
		.axisSegmentSize = .h / .segments
		.h -= .axisSegmentSize + .sz // this is for adding the axis height and text height
		maxDigits = String(Max(Abs(.range[0]), Abs(.range[1]))).Size() + 1
		.x += maxDigits * .fontOffset
		.w -= maxDigits * .sz
		incr = .w / (.axisDensity - 1)
		.DrawAxis(TA.CENTER, { .drawHorizontalMarks(incr) })
		}

	drawHorizontalMarks(incr)
		{
		lineY = .y + .h
		lineH = lineY + .h / .segments
		fontY = lineY + .h / .segments * 2
		for i in .. .axisDensity
			{
			text = Display(((.range[1] - .range[0]) * i * incr / .w + .range[0]).Round(1))
			lineOb = Object(x: .x + i * incr, y: lineY, h: lineH, w: .w)
			.DrawHorizontalMarks(text, lineOb, fontY)
			}
		}

	DrawHorizontalMarks(text /*unused*/, lineOb /*unused*/, fontY /*unused*/)
		{
		throw 'must be implemented by derived class'
		}

	setupColors()
		{
		// Extract single RGB components from color with a mask:
		r = (.color & 0xFF0000) / 0x10000 	/*= RGB: red value */
		g = (.color & 0x00FF00) / 0x100		/*= RGB: green value */
		b = (.color & 0x0000FF)				/*= RGB: blue value */

		// percentages of fading:
		fadingRateBad = 80
		fadingRateSatisfactory = 120
		fadingRateGood = 150

		.badColor = .bgr(r, g, b, fadingRateBad)
		.satisfactoryColor = .bgr(r, g, b, fadingRateSatisfactory)
		.goodColor = .bgr(r, g, b, fadingRateGood)
		.valueColor = .bgr(r, g, b, 0)
		}

	colorMax: 255
	bgr(r, g, b, fadingRate)
		{
		shade = Max(r, g, b)
		incr = (shade * fadingRate).PercentToDecimal().Int()
		r = Min(.colorMax, r + incr)
		g = Min(.colorMax, g + incr)
		b = Min(.colorMax, b + incr)
		return RGB(r, g, b)
		}

	DrawEnclosingRectangle()
		{
		throw 'must be implemented by derived class'
		}

	drawBars()
		{
		if .vertical
			.drawVerticalBars()
		else
			.drawHorizontalBars()
		}

	targetIndicatorRatio: 4
	targetIndicatorEnd: 3
	drawVerticalBars()
		{
		// Draw good values bar:
		.DrawBar(.x, .y + .h, .x + .w, .y, .goodColor)

		if .data is false
			return

		// Draw satisfactory values bar:
		.DrawBar(.x, .y + .h, .x + .w, .y + .h - .h * .good / .normRange[1],
			.satisfactoryColor)

		// Draw bad values bar:
		.DrawBar(.x, .y + .h, .x + .w, .y + .h - .h * .satisfactory / .normRange[1],
			.badColor)

		// Draw actual value bar:
		.DrawBar(.x + .w / .actualBarRatio, .y + .h,
			.x + .w / .actualBarRatio * 2,
			.y + .h - .h * .data / .normRange[1], .valueColor)

		// Draw target value indicator:
		lineOb = Object(
			x: .x + .w / .targetIndicatorRatio
			y: .y + .h - .h * .target / .normRange[1]
			w: .x + .w / .targetIndicatorRatio * .targetIndicatorEnd)
		.DrawVerticalBars(lineOb)
		}

	DrawBar(x1 /*unused*/, y1 /*unused*/, x2 /*unused*/, y2 /*unused*/, color /*unused*/)
		{
		throw 'must be implemented by derived class'
		}

	DrawVerticalBars(lineOb /*unused*/)
		{
		throw 'must be implemented by derived class'
		}

	actualBarRatio: 3
	drawHorizontalBars()
		{
		// Draw good values bar:
		.DrawBar(.x, .y, .x + .w, .y + .h, .goodColor)

		if .data is false
			return

		// Draw satisfactory values bar:
		.DrawBar(.x, .y, .x + .w * .good / .normRange[1], .y + .h, .satisfactoryColor)

		// Draw bad values bar:
		.DrawBar(.x, .y, .x + .w * .satisfactory / .normRange[1], .y + .h, .badColor)

		// Draw actual value bar:
		.DrawBar(.x, .y + .h / .actualBarRatio, .x + .w * .data / .normRange[1],
			.y + .h / .actualBarRatio * 2, .valueColor)

		// Draw target value indicator:
		lineOb = Object(
			x: .x + .w * .target / .normRange[1]
			y: .y + .h / .targetIndicatorRatio
			h: .y + .h / .targetIndicatorRatio * .targetIndicatorEnd)
		.DrawHorizontalBars(lineOb)
		}

	DrawHorizontalBars(lineOb /*unused*/)
		{
		throw 'must be implemented by derived class'
		}

	PostDraw()
		{ }
	}