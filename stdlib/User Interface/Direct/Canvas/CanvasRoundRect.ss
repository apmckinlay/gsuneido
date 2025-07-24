// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
CanvasItem
	{
	relativeSizeDivisor: 3
	New(x1, y1, x2, y2, .width = false, .height = false, .name = '')
		{
		.sortPoints(x1, y1, x2, y2)
		.posLocked? = false
		if width is false
			.width = Min(.x2 - .x1, .y2 - .y1) / .relativeSizeDivisor
		if height is false
			.height = .width
		}
	sortPoints(x1, y1, x2, y2)
		{
		.x1 = Min(x1, x2)
		.y1 = Min(y1, y2)
		.x2 = Max(x1, x2)
		.y2 = Max(y1, y2)
		}
	Paint()
		{
		_report.AddRoundRect(.x1, .y1, .x2-.x1, .y2-.y1, .width, .height, .GetThick(),
			.GetColor(), .GetLineColor())
		}
	BoundingRect()
		{
		return Object(x1: .x1, y1: .y1, x2: .x2, y2: .y2)
		}
	SetSize(x1, y1, x2, y2)
		{
		.x1 = x1
		.y1 = y1
		.x2 = x2
		.y2 = y2
		}
	ResetSize()
		{
		result = ResetSizeControl(0,
			Object(x1: .x1, y1: .y1, x2: .x2, y2: .y2, width: .width, height: .height))
		if result is false
			return
		x1 = Number(result.x1)
		y1 = Number(result.y1)
		x2 = Number(result.x2)
		y2 = Number(result.y2)
		.sortPoints(x1, y1, x2, y2)
		.width = Number(result.width)
		.height = Number(result.height)
		}
	StringToSave()
		{
		'CanvasRoundRect(x1: ' $ Display(.x1) $ ', y1: ' $ Display(.y1) $
			', x2: ' $ Display(.x2) $ ', y2: ' $ Display(.y2) $
			', width: ' $ Display(.width) $ ', height: ' $ Display(.height) $ ')'
		}
	ObToSave()
		{
		return Object('CanvasRoundRect', .x1, .y1, .x2, .y2, .width, .height)
		}
	Resize(origx, origy, x, y)
		{
		if .Resizing?(.x1, origx)
			.x1 = x
		if .Resizing?(.y1, origy)
			.y1 = y
		if .Resizing?(.x2, origx)
			.x2 = x
		if .Resizing?(.y2, origy)
			.y2 = y
		.sortPoints(.x1, .y1, .x2, .y2)
		.width = Min(.x2 - .x1, .y2 - .y1) / .relativeSizeDivisor
		.height = .width
		}
	Scale(by)
		{
		.x1 *= by
		.x2 *= by
		.y1 *= by
		.y2 *= by
		.sortPoints(.x1, .y1, .x2, .y2)
		.width = Min(.x2 - .x1, .y2 - .y1) / .relativeSizeDivisor
		.height = .width
		}
	Move(dx, dy)
		{
		if .posLocked?
			return
		.x1 += dx
		.x2 += dx
		.y1 += dy
		.y2 += dy
		}
	GetName()
		{ return .name }
	GetSuJSObject()
		{
		return Object('SuCanvasRoundRect', .x1, .y1, .x2, .y2, .width, .height, id: .Id)
		}
	GetCoordinates()
		{
		return Object(.x1, .y1, .x2, .y2, .width, .height)
		}
	ToggleLock()
		{
		.posLocked? = not .posLocked?
		}
	}
