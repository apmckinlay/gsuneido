// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	// interface:
	Test_Main()
		{
		pt = Point(20, 30)
		Assert(pt.GetX() is: 20)
		Assert(pt.GetY() is: 30)
		Assert(pt.ToWindowsPoint() is: Object(x: 20, y: 30))
		Assert(Display(pt) is: "Point(20, 30)")
		}
	Test_Translate()
		{
		pt = Point(0, 0)
		Assert(pt.Translate(0, 0) is: pt)
		Assert(pt is: Point(0, 0))
		pt = pt.Translate(1, 0)
		Assert(pt is: Point(1, 0))
		pt = pt.Translate(0, 1)
		Assert(pt is: Point(1, 1))
		pt = pt.Translate(-0.5, -0.5)
		Assert(pt is: Point(0.5, 0.5))
		}
	}
