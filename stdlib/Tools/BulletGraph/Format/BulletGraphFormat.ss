// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
/*
	BulletGraphFormat - Written in February 2007 by Mauro Giubileo
	--------------------------------------------------------------

	Example Usage:
	Params(#(BulletGraph 24, satisfactory: 20, good: 25, target: 27, range: (0,30)))
*/
Format
	{
	standardizationFactor: 15 /* standardizes size arguments with the control */
	New(data = false, .satisfactory = 0, .good = 0, .target = 0, .range = #(0, 100),
		.color = 0x506363, .width = 128, .height = 32, .rectangle = true,
		.outside = 5, .vertical = false, .axis = false, .axisDensity = 5)
		{
		if .vertical and width is 128 and height is 32 /*= default vertical*/
			{ // swap w and h
			temp = .width
			.width = .height
			.height = temp
			}

		.Data = data
		.width *= .standardizationFactor
		.height *= .standardizationFactor
		.outside *= .standardizationFactor

		if .vertical
			{
			.drawAxisMethod = .drawVerticalAxis
			.drawBarMethod = .drawVerticalBars
			}
		else
			{
			.drawAxisMethod = .drawHorizontalAxis
			.drawBarMethod = .drawHorizontalBars
			}
		.setupColors()
		}

	setupColors()
		{
		colors = BulletGraphColors(.color)
		.badColor = colors.bad
		.satisfactoryColor = colors.satisfactory
		.goodColor = colors.good
		.valueColor = colors.value
		}

	GetSize(data /*unused*/ = '')
		{
		//Post: returns width, height and descent of the format as Object(w:, h:, d:)
		return Object(w: .width, h: .height, d: 0)
		}

	Print(x, y, w, h, data = '')
		{
		//Post: prints the format at position x,y with specified width and height
		.data = .Data isnt false
			? .Data
			: data
		.validateData()
		.normalization()
		.Draw(x, y, w, h)
		}

	validateData()
		{
		Assert(.data < .range[1])
		Assert(.satisfactory >= .range[0] and .satisfactory <= .range[1])
		Assert(.good >= .range[0] and .good <= .range[1])
		Assert(.target >= . range[0] and .target <= .range[1])
		Assert(.range[0] < .range[1])
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

	segments: 12
	actualBarRatio: 3
	targetIndicatorEnd: 3
	targetIndicatorRatio: 4
	Draw(.x, .y, .w, .h)
		{
		// check if outside area is too big in respect to width and height:
		if .w <= 0 or .h <= 0
			return

		.Rect = Object(x: .x, y: .y, w: .w, h: .h)

		// add outside spacing:
		.x += .outside
		.y += .outside
		.w -= .outside * 2
		.h -= .outside * 2
		.sz = Max(.w, .h) / 24  /*= font size */
		.fontOffset = .sz / 2

		Format.DoWithFont(Object(name: 'Arial', size: .sz / 15 /*= offset */))
			{|unused|
			.drawGraph(.x, .y, .w, .h)
			}
		}

	// ========================== Draw Vertical ==========================
	drawVerticalAxis()
		{
		.h -= .sz // this is for adding the text height
		incr = .h / (.axisDensity - 1)
		.axisSegmentSize = .w / .segments
		maxDigits = String(Max(Abs(.range[0]), Abs(.range[1]))).Size() + 1
		.x += maxDigits * .sz
		.y += .fontOffset
		.w -= maxDigits * .sz + .fontOffset + .axisSegmentSize
		.drawVerticalMarks(incr, maxDigits)
		.x += .fontOffset + .axisSegmentSize
		}

	drawVerticalMarks(incr, maxDigits)
		{
		textWidth = maxDigits * .sz
		textLength = .fontOffset + .axisSegmentSize
		lineX = .x + .fontOffset
		lineW = .x + textLength
		fontX = lineW - textWidth - textLength
		justify = 'right'
		for i in .. .axisDensity
			{
			text = Display((.range[1] * i * incr / .h + .range[0]).Round(1))
			y = .y + .h - i * incr
			_report.AddLine(lineX, y, lineW, y, thick: .5)
			Format.DoWithFont(false)
				{|font|
				_report.
					AddText(text, fontX, y - .sz / 2, textWidth, .h, font, :justify)
				}
			}
		}

	drawVerticalBars()
		{
		// Draw good values bar:
		.drawBar(.x, .y + .h, .x + .w, .y, .goodColor)

		// Draw satisfactory values bar:
		.drawBar(.x, .y + .h, .x + .w, .y + .h - .h * .good / .range[1],
			.satisfactoryColor)

		// Draw bad values bar:
		.drawBar(.x, .y + .h, .x + .w, .y + .h - .h * .satisfactory / .range[1],
			.badColor)

		// Draw actual value bar:
		if .data isnt false
			.drawBar(
				.x + .w / .actualBarRatio,
				.y + .h,
				.x + .w / .actualBarRatio * 2,
				.y + .h - .h * .data / .range[1],
				.valueColor)

		// Draw target value indicator:
		if .target > 0
			{
			x = .x + .w / .targetIndicatorRatio
			y = .y + .h - .h * .target / .range[1]
			w = .x + .w / .targetIndicatorRatio * .targetIndicatorEnd
			_report.AddLine(x, y, w, y,	.targetSize())
			}
		}

	// ========================= Draw Horizontal =========================
	drawHorizontalAxis()
		{
		.axisSegmentSize = .h / .segments
		.h -= .axisSegmentSize + .sz // this is for adding the axis height and text height
		maxDigits = String(Max(Abs(.range[0]), Abs(.range[1]))).Size() + 1
		.x += maxDigits * .fontOffset
		.w -= maxDigits * .sz
		incr = .w / (.axisDensity - 1)
		.drawHorizontalMarks(incr)
		}

	drawHorizontalMarks(incr)
		{
		lineY = .y + .h
		lineH = lineY + .h / .segments
		fontY = lineY + .h / .segments * 2
		for i in .. .axisDensity
			{
			text = Display((.range[1] * i * incr / .w + .range[0]).Round(1))
			x = .x + i * incr
			_report.AddLine(x, lineY, x, lineH, thick: .5)
			Format.DoWithFont(false)
				{|font|
				_report.
					AddText(text, x - .sz * (text.Size() - 1) / 2, fontY, .w, lineH, font)
				}
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
			.drawBar(
				.x,
				.y + .h / .actualBarRatio,
				.x + .w * .data / .range[1],
				.y + .h / .actualBarRatio * 2,
				.valueColor)

		// Draw target value indicator:
		if .target > 0
			{
			x = .x + .w * .target / .range[1]
			y = .y + .h / .targetIndicatorRatio
			h = .y + .h / .targetIndicatorRatio * .targetIndicatorEnd
			_report.AddLine(x, y, x, h,	.targetSize())
			}
		}

	// ========================== Draw General ===========================
	drawGraph(.x, .y, .w, .h)
		{
		if .axis
			(.drawAxisMethod)()
		(.drawBarMethod)()
		if .rectangle
			_report.AddRect(.x, .y, .w, .h, thick: 1)
		}

	drawBar(x1, y1, x2, y2, color)
		{
		_report.AddRect(x1, y1, x2 - x1, y2 - y1, thick: 0, fillColor: color)
		}

	targetSize()
		{
		return Max(0.1, (.w * 0.0005)) /*= thick */
		}
	}
