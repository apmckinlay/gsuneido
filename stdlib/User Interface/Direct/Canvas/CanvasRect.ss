// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
CanvasItem
	{
	New(x1, y1, x2, y2, .name = '')
		{
		.sortPoints(x1, y1, x2, y2)
		.posLocked? = false
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
		_report.AddRect(.x1, .y1, .x2-.x1, .y2-.y1, .GetThick(), .GetColor(),
			.GetLineColor())
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
		result = ResetSizeControl(0, Object(x1: .x1, y1: .y1, x2: .x2, y2: .y2))
		if (result is false)
			return
		x1 = Number(result.x1)
		y1 = Number(result.y1)
		x2 = Number(result.x2)
		y2 = Number(result.y2)
		.sortPoints(x1, y1, x2, y2)
		}
	StringToSave()
		{
		'CanvasRect(x1: ' $ Display(.x1) $ ', y1: ' $ Display(.y1) $
			', x2: ' $ Display(.x2) $ ', y2: ' $ Display(.y2) $ ')'
		}
	ObToSave()
		{
		return Object('CanvasRect', .x1, .y1, .x2, .y2)
		}
	Resize(origx, origy, x, y, literal = false)
		{
		if literal
				{
				.x1 = origx
				.y1 = origy
				.x2 = x
				.y2 = y
				return
				}
		if .Resizing?(.x1, origx)
			.x1 = x
		if .Resizing?(.y1, origy)
			.y1 = y
		if .Resizing?(.x2, origx)
			.x2 = x
		if .Resizing?(.y2, origy)
			.y2 = y
		.sortPoints(.x1, .y1, .x2, .y2)
		}
	Scale(by)
		{
		.x1 *= by
		.x2 *= by
		.y1 *= by
		.y2 *= by
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
		return Object('SuCanvasRect', .x1, .y1, .x2, .y2, id: .Id)
		}
	ToggleLock()
		{
		.posLocked? = not .posLocked?
		}
	}
