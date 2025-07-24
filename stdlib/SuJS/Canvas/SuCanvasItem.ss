// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
class
	{
	New()
		{
		.canvas = _canvas
		.Id = _spec.id
		}

	Getter_Driver()
		{
		return .canvas.Driver
		}

	Getter_Canvas()
		{
		return .canvas
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

	InHandleArea(handlex, handley, x, y)
		{
		detectBoundary = .handleBoundary + 1
		return x >= handlex - detectBoundary and x <= handlex + detectBoundary and
			y >= handley - detectBoundary and y <= handley + detectBoundary
		}

	BoundingRect()
		{
		return #(x1: 0, x2: 0, y1: 0, y2: 0)
		}

	handles: #()
	boundingRect: false
	boundingLine: false
	Select()
		{
		.paintHandles()
		}

	Unselect()
		{
		.removeHandles()
		}

	handleBoundary: 3
	paintHandles()
		{
		.removeHandles()
		.ForeachHandle()
			{ |x, y|
			.handles.Add(Object(
				el: .Driver.AddRect(x - .handleBoundary, y - .handleBoundary,
					.handleBoundary * 2, .handleBoundary * 2, 1, #black),
				:x,
				:y))
			}
		}

	removeHandles()
		{
		for handle in .handles
			.Driver.Remove(handle.el)
		.handles = Object()
		}

	updateHandles()
		{
		if .handles.NotEmpty?()
			{
			.removeHandles()
			.paintHandles()
			}
		}

	Move(dx, dy)
		{
		for handle in .handles
			{
			handle.x += dx
			handle.y += dy
			.Driver.MoveRect(handle.el, dx, dy)
			}
		.Send('Move', dx, dy)
		}

	Resizing?(prev, orig)
		{
		return prev >= orig - 4 and prev <= orig + 4 /*= offset of resize*/
		}

	Resize(origx, origy, x, y)
		{
		.updateHandles()
		.Send('Resize', origx, origy, x, y)
		}

	MoveToBack()
		{
		for el in .GetElements()
			.Driver.MoveToBack(el)
		}

	MoveToFront()
		{
		for el in .GetElements()
			.Driver.MoveToFront(el)
		}

	GetElements()
		{
		return #()
		}

	ResetSize(@unused)
		{
		.updateHandles()
		.RemoveBoundingRect()
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
		}


	PaintBoundingRect(r)
		{
		if r.x1 is r.x2 or r.y1 is r.y2
			{
			.removeBoundingRect()
			.paintBoundingLine(r)
			}
		else
			{
			.removeBoundingLine()
			.paintBoundingRect(r)
			}
		}

	paintBoundingRect(r)
		{
		if .boundingRect is false
			{
			.boundingRect = .canvas.Driver.AddRect(
				r.x1, r.y1, r.x2 - r.x1, r.y2 - r.y1, 1)
			.boundingRect.SetAttribute('stroke-dasharray', '5,5')
			}
		else
			.Driver.ResizeRect(.boundingRect, r.x1, r.y1, r.x2 - r.x1, r.y2 - r.y1)
		}

	paintBoundingLine(r)
		{
		if .boundingLine is false
			{
			.boundingLine = .Driver.AddLine(r.x1, r.y1, r.x2, r.y2, 1)
			.boundingLine.SetAttribute('stroke-dasharray', '5,5')
			}
		else
			.Driver.ResizeLine(.boundingLine, r.x1, r.y1, r.x2, r.y2)
		}

	RemoveBoundingRect()
		{
		.removeBoundingLine()
		.removeBoundingRect()
		}

	removeBoundingRect()
		{
		if .boundingRect isnt false
			{
			.Driver.Remove(.boundingRect)
			.boundingRect = false
			}
		}

	removeBoundingLine()
		{
		if .boundingLine isnt false
			{
			.Driver.Remove(.boundingLine)
			.boundingLine = false
			}
		}

	Remove()
		{
		for el in .GetElements()
			.Driver.Remove(el)
		.removeHandles()
		.RemoveBoundingRect()
		}

	AfterEdit(@unused)
		{
		.updateHandles()
		.RemoveBoundingRect()
		}

	Send(@args)
		{
		.canvas.Event('ToItem', .Id, args)
		}
	}
