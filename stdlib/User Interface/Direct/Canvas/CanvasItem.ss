// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Grouped?: false
	Scale(unused)
		{
		throw ".Scale(by) must be defined by sub-classes of CanvasItem"
		}
	StringToSave()
		{
		throw ".StringToSave() must be defined by sub-classes of CanvasItem"
		}
	ObToSave()
		{
		throw ".ObToSave() must be defined by sub-classes of CanvasItem"
		}
	Getter_ScaleBy()
		{
		// need the win32 check because this is called by DrawItem.BuildItems
		// which is called by demo data when it's running generate invoice reports
		return .ScaleBy = Sys.Win32?() ? GetDpiFactor() : 1
		}
	Paint()
		{
		}
	handleBoundary: 3
	HandleBoundary()
		{
		return .handleBoundary * GetDpiFactor()
		}
	PaintHandles(hdc)
		{
		WithHdcSettings(hdc, [GetStockObject(SO.BLACK_BRUSH), SetROP2: R2.NOTXORPEN])
			{
			.ForeachHandle()
				{ |x,y|
				Rectangle(hdc, left: x - .HandleBoundary(), right: x + .HandleBoundary(),
					top: y - .HandleBoundary(), bottom: y + .HandleBoundary())
				}
			}
		}
	ForeachHandle(block)
		{
		r = .BoundingRect()
		block(r.x1, r.y1)
		block(r.x1, r.y2)
		block(r.x2, r.y1)
		block(r.x2, r.y2)
		x = (r.x1 + r.x2) / 2
		y = (r.y1 + r.y2) / 2
		block(x, r.y1)
		block(x, r.y2)
		block(r.x1, y)
		block(r.x2, y)
		.handles = Object(Object(x: r.x1, y: r.y1), Object(x: r.x1, y: r.y2),
			Object(x: r.x2, y: r.y1), Object(x: r.x2, y: r.y2),
			Object(:x, y: r.y1), Object(:x, y: r.y2),
			Object(x: r.x1, :y), Object(x: r.x2, :y))
		}
	color: 0xffffff /*=CLR.WHITE*/
	SetColor(color)
		{
		.color = color
		return this
		}
	GetColor()
		{
		return .color
		}
	lin_color: 0
	SetLineColor(color)
		{
		.lin_color = color
		return this
		}
	GetLineColor()
		{
		return .lin_color
		}
	thick: 1
	SetThick(thick)
		{
		.thick = thick
		return this
		}
	GetThick()
		{
		return .thick
		}
	handles: #()
	GetHandles()
		{
		return .handles
		}
	InHandleArea(handlex, handley, x, y)
		{
		detectBoundary = .HandleBoundary() + 1
		return x >= handlex - detectBoundary and x <= handlex + detectBoundary and
			y >= handley - detectBoundary and y <= handley + detectBoundary
		}
	Contains(x, y)
		{
		r = .BoundingRect()
		return r.x1 <= x and x <= r.x2 and r.y1 <= y and y <= r.y2
		}
	Overlaps?(x1, y1, x2, y2)
		{
		r = .BoundingRect()
		return Rect.LinearOverlap?(x1, x2, r.x1, r.x2) and
			   Rect.LinearOverlap?(y1, y2, r.y1, r.y2)
		}
	IsHandle?(x, y)
		{
		for handle in .handles
			if (.InHandleArea(handle.x, handle.y, x, y))
				return true
		return false
		}
	BoundingRect()
		{
		return #(x1: 0, x2: 0, y1: 0, y2: 0)
		}
	GetItems()
		{
		return #()
		}
	Edit()
		{
		}
	GetResource()
		{
		return #()
		}
	FromSVG(ob/*unused*/)
		{
		return ''
		}
	Destroy()
		{
		}
	GetName()
		{ return '' }
	ReverseScale(by)
		{
		reverseBy = 1 / by
		.Scale(reverseBy)
		}
	Resizing?(prev, orig)
		{
		return prev >= orig - 4 and prev <= orig + 4 /*= offset of resize*/
		}
	SetupScale()
		{
		if .ScaleBy isnt 1
			.Scale(.ScaleBy)
		return this
		}
	ToString()
		{
		if .ScaleBy is 1
			return .StringToSave()

		.ReverseScale(.ScaleBy)
		res = .StringToSave()
		.Scale(.ScaleBy)
		return res
		}

	Get()
		{
		if .ScaleBy is 1
			return .ObToSave()

		.ReverseScale(.ScaleBy)
		res = .ObToSave()
		.Scale(.ScaleBy)
		return res
		}

	GetSuJSObject()
		{
		return false
		}

	GetCoordinates()
		{
		if false is ob = .GetSuJSObject()
			return false
		return ob[1::4/*=x1, y1, x2, y2*/]
		}
	}
