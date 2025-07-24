// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	UseDeepEquals: true

	New(.x, .y, .width, .height)
		{
		Assert(Number?(x) and Number?(y) and Number?(width) and Number?(height))
		}

	GetX()
		{
		return .x
		}
	GetY()
		{
		return .y
		}
	GetX2()
		{
		return .x + .width
		}
	GetY2()
		{
		return .y + .height
		}
	GetWidth()
		{
		return .width
		}
	GetHeight()
		{
		return .height
		}

	Set(x = false, y = false, width = false, height = false)
		{
		.assertSetValues(x, y, width, height)
		if x isnt false
			.x = x
		if y isnt false
			.y = y
		if width isnt false
			.width = width
		if height isnt false
			.height = height
		return this
		}

	assertSetValues(x, y, width, height)
		{
		Assert(x is false or Number?(x))
		Assert(y is false or Number?(y))
		Assert(width is false or Number?(width))
		Assert(height is false or Number?(height))
		}

	SetX(x)
		{
		Assert(x, isNumber:)
		.x = x
		}
	SetY(y)
		{
		Assert(y, isNumber:)
		.y = y
		}
	SetWidth(width)
		{
		Assert(width, isNumber:)
		.width = width
		}
	SetHeight(height)
		{
		Assert(height, isNumber:)
		.height = height
		}

	TranslateX(dx)
		{
		Assert(dx, isNumber:)
		.x += dx
		}
	TranslateY(dy)
		{
		Assert(dy, isNumber:)
		.y += dy
		}
	Translate(dx, dy)
		{
		Assert(dx, isNumber:)
		Assert(dy, isNumber:)
		.x += dx
		.y += dy
		}
	FromWindowsRect(winrc) // static method
		{
		return Rect(winrc.left, winrc.top,
			winrc.right - winrc.left, winrc.bottom - winrc.top)
		}
	ToWindowsRect()
		{
		return Object(left: .x, top: .y, right: .x + .width, bottom: .y + .height)
		}
	IntoWindowsRect(winrc)
		{
		winrc.left   = .x
		winrc.right  = .x + .width
		winrc.top    = .y
		winrc.bottom = .y + .height
		return winrc
		}
	Overlaps?(rect)
		{
		return .LinearOverlap?(.x, .x + .width, rect.GetX(), rect.GetX2()) and
			   .LinearOverlap?(.y, .y + .height, rect.GetY(), rect.GetY2())
		}
	OverlapsWindowsRect?(winrc)
		{
		return .LinearOverlap?(.x, .x + .width, winrc.left, winrc.right) and
			   .LinearOverlap?(.y, .y + .height, winrc.top, winrc.bottom)
		}
	Union(rect)
		{
		x = Min(.x, rect.GetX())
		y = Min(.y, rect.GetY())
		right = Max(.x + .width, rect.GetX2)
		bottom = Max(.y + .height, rect.GetY2)
		return Rect(x, y, right - x, bottom - y)
		}
	ContainsPoint?(pt)
		{
		x = pt.GetX()
		y = pt.GetY()
		return .x <= x and x <= .x + .width and
			   .y <= y and y <= .y + .height
		}
	TrapPoint(pt)
		{
		// if pt is within this Rect, returns pt
		// otherwise, returns a new point
		// which is the closest point within this Rect
		x = pt.GetX()
		y = pt.GetY()
		x2 = .x + .width
		y2 = .y + .height
		if .x <= x and x <= x2 and
		   .y <= y and y <= y2  // point is in rect
			return pt
		else
			return .closestPointWithinRect(x, y, x2, y2)
		}

	closestPointWithinRect(x, y, x2, y2)
		{
		if (x < .x) x = .x
		else if (x2 < x) x = x2
		if (y < .y) y = .y
		else if (y2 < y) y = y2
		return new Point(x, y)
		}

	LinearOverlap?(org1, dest1, org2, dest2)
		{
		// TODO[VCS]: I question whether overlap should return true when one
		//			  range just includes the rightmost edge of the other one
		Assert(org1 <= dest1 and org2 <= dest2)
		return org1 <= dest2 and dest1 >= org2
		// Same as: not (org1 > dest2 or dest1 < org2)
		}
	ToString()
		{
		return "Rect(" $ String(.x) $ ", " $ String(.y) $ ", " $
					String(.width) $ ", " $ String(.height) $ ")"
		}
	}