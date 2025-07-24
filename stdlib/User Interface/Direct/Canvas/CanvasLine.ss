// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
CanvasItem
	{
	New(.x1, .y1, .x2, .y2, .name = '')
		{
		.posLocked? = false
		}
	Paint()
		{
		_report.AddLine(.x1, .y1, .x2, .y2, .GetThick(), .GetLineColor())
		}
	BoundingRect()
		{
		return Object(x1: Min(.x1, .x2), y1: Min(.y1, .y2),
			x2: Max(.x1, .x2), y2: Max(.y1, .y2))
		}
	handles: #()
	ForeachHandle(block)
		{
		block(.x1, .y1)
		block(.x2, .y2)
		.handles = Object(Object(x: .x1, y: .y1), Object(x: .x2, y: .y2))
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
		if false is size = ResetSizeControl(0, Object(x1: .x1, y1: .y1, x2: .x2, y2: .y2))
			return
		.x1 = Number(size.x1)
		.y1 = Number(size.y1)
		.x2 = Number(size.x2)
		.y2 = Number(size.y2)
		}
	StringToSave()
		{
		'CanvasLine(x1: ' $ Display(.x1) $ ', y1: ' $ Display(.y1) $
			', x2: ' $ Display(.x2) $ ', y2: ' $ Display(.y2) $ ')'
		}
	ObToSave()
		{
		return Object('CanvasLine', .x1, .y1, .x2, .y2)
		}
	IsHandle?(x, y)
		{
		for handle in .handles
			if .InHandleArea(handle.x, handle.y, x, y)
				return true
		return false
		}
	Resize(origx, origy, x, y)
		{
		if .Resizing?(.x1, origx) and .Resizing?(.y1, origy)
			{
			.x1 = x
			.y1 = y
			}
		if .Resizing?(.x2, origx) and .Resizing?(.y2, origy)
			{
			.x2 = x
			.y2 = y
			}
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
		{
		return .name
		}
	GetSuJSObject()
		{
		return Object('SuCanvasLine', .x1, .y1, .x2, .y2, id: .Id)
		}
	ToggleLock()
		{
		.posLocked? = not .posLocked?
		}
	}
