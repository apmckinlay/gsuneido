// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
CanvasItem
	{
	New(.left, .top, .right, .bottom, .xStartArc, .yStartArc, .xEndArc, .yEndArc,
		.name = '')
		{
		.posLocked? = false
		}
	Paint()
		{
		_report.AddArc(.left, .top, .right, .bottom, .xStartArc, .yStartArc, .xEndArc,
			.yEndArc, .GetThick(), .GetLineColor())
		}
	BoundingRect()
		{
		return Object(
			x1: Min(.xStartArc, .xEndArc), y1: Min(.yStartArc, .yEndArc),
			x2: Max(.xStartArc, .xEndArc), y2: Max(.yStartArc, .yEndArc))
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
			Object(left: .left, top: .top, right: .right, bottom: .bottom,
				xStartArc: .xStartArc, yStartArc: .yStartArc, xEndArc: .xEndArc,
				yEndArc: .yEndArc))
		if (result is false)
			return
		.left = Number(result.left)
		.top = Number(result.top)
		.right = Number(result.right)
		.bottom = Number(result.bottom)
		.xStartArc = Number(result.xStartArc)
		.yStartArc = Number(result.yStartArc)
		.xEndArc = Number(result.xEndArc)
		.yEndArc = Number(result.yEndArc)
		}
	StringToSave()
		{
		'CanvasArc(left: ' $ Display(.left) $ ', top: ' $ Display(.top) $
			', right: ' $ Display(.right) $ ', bottom: ' $ Display(.bottom) $
			', xStartArc: ' $ Display(.xStartArc) $
			', yStartArc: ' $ Display(.yStartArc) $
			', xEndArc: ' $ Display(.xEndArc) $ ', yEndArc: ' $ Display(.yEndArc) $
			')'
		}
	ObToSave()
		{
		return Object('CanvasArc', .left, .top, .right, .bottom,
			.xStartArc, .yStartArc, .xEndArc, .yEndArc)
		}
	Resize(origx, origy, x, y)
		{
		if .Resizing?(.left, origx)
			.left = x
		if .Resizing?(.top, origy)
			.top = y
		if .Resizing?(.right, origx)
			.right = x
		if .Resizing?(.bottom, origy)
			.bottom = y
		.xStartArc = .left
		.xEndArc = .right
		.yStartArc = .top
		.yEndArc = .bottom
		}
	Scale(by)
		{
		.left *= by
		.top *= by
		.right *= by
		.bottom *= by

		.xStartArc = .left
		.xEndArc = .right
		.yStartArc = .top
		.yEndArc = .bottom
		}
	Move(dx, dy)
		{
		if .posLocked?
			return
		.left += dx
		.right += dx
		.xStartArc += dx
		.xEndArc += dx
		.top += dy
		.bottom += dy
		.yStartArc += dy
		.yEndArc += dy
		}
	GetName()
		{ return .name }

	GetSuJSObject()
		{
		return Object('SuCanvasArc', .left, .top, .right, .bottom,
			.xStartArc, .yStartArc, .xEndArc, .yEndArc, id: .Id)
		}

	GetCoordinates()
		{
		return Object(.left, .top, .right, .bottom,
			.xStartArc, .yStartArc, .xEndArc, .yEndArc)
		}
	ToggleLock()
		{
		.posLocked? = not .posLocked?
		}
	}