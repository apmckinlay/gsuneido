// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// FIXME: pdf output is not working
Format
	{
	New(.item)
		{
		}
	GetSize(data /*unused*/ = false)
		{
		if .item is #()
			return Object(w: 0, h: 0, d: 0)
		.rect = .item.BoundingRect()
		// this is the conversion between resizing from DrawControl to format size
		.conversion = 17 / .item.ScaleBy
		h = (.rect.y1 - .rect.y2).Abs() * .conversion
		w = (.rect.x1 - .rect.x2).Abs() * .conversion

		// to avoid canvas size too large for page error
		h = Min(h, _report.GetDimens().H)
		w = Min(w, _report.GetDimens().W)
		return Object(:w, :h, d: 0)
		}
	Print(x, y, w, h, data = false)
		{
		if .item is #()
			return

		ratioX = w / (.rect.x2 - .rect.x1).Abs()
		ratioY = h / (.rect.y2 - .rect.y1).Abs()
		ratio = Min(ratioX, ratioY)

		dx = x - Min(.rect.x1, .rect.x2) * ratio
		dy = y - Min(.rect.y1, .rect.y2) * ratio
		.Draw(dx, dy, ratio, data)
		}
	Draw(dx, dy, ratio, data)
		{
		items = .item.GetItems()
		oldThick = false
		for item in items
			{
			if Sys.SuneidoJs?()
				{
				oldThick = item.GetThick()
				item.SetThick(10/*= default thick of Format*/)
				}
			item.Scale(ratio, print?:)
			item.Move(dx, dy)
			}
		.item.Paint(:data)
		for item in items
			{
			if oldThick isnt false
				item.SetThick(oldThick)
			item.Move(-dx, -dy)
			item.ReverseScale(ratio, print?:)
			}
		}
	}
