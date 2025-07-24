// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// immutable
class
	{
	UseDeepEquals: true
	New(.x, .y)
		{
		Assert(Number?(x) and (Number?(y)))
		}
	GetX()
		{
		return .x
		}
	GetY()
		{
		return .y
		}
	ToWindowsPoint()
		{
		return Object(x: .x, y: .y)
		}
	FromWindowsPoint(winpt) // static, returns new Point
		{
		return new Point(winpt.x, winpt.y)
		}
	Translate(dx, dy)
		{
		return new Point(.x + dx, .y + dy)
		}
	ToString()
		{
		return "Point(" $ String(.x) $ ", " $ String(.y) $ ")"
		}
	}
