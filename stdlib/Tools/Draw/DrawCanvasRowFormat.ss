// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
DrawItemFormat
	{
	New(.item, .canvasWidth, .canvasHeight = false)
		{
		super(item)
		.width = .canvasWidth.InchesInCanvasUnit()
		.height = .canvasHeight is false ? false : .canvasHeight.InchesInCanvasUnit()
		}
	GetSize(data /*unused*/ = false)
		{
		if .item is #()
			return Object(w: 0, h: 0, d: 0)

		.rect = .item.BoundingRect()
		// this is the conversion between resizing from DrawControl to format size
		.conversion = 17 / .item.ScaleBy

		.w = .width * .item.ScaleBy * .conversion
		w = _report.GetDimens().W

		if .height is false
			{
			.h = Max(.rect.y1, .rect.y2) * .conversion
			h = .h * _report.GetDimens().W / .w
			}
		else
			{
			.h = .height * .item.ScaleBy * .conversion
			h = _report.GetDimens().H
			}

		return Object(:w, :h, d: 0)
		}

	Data: false
	Print(x, y, w, h, data = false)
		{
		if .item is #()
			return

		if .Data isnt false
			data = .Data

		ratioX = w / .w * .conversion
		ratioY = h / .h * .conversion
		ratio = Min(ratioX, ratioY)

		dx = x
		dy = y

		.Draw(dx, dy, ratio, data)
		}

	Header()
		{
		return false
		}
	}
