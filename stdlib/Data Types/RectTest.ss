// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// SuJsWebTest
Test
	{
	// interface:
	Test_Main()
		{
		rect1 = Rect(5, 45, 105, 85)
		Assert(rect1.GetX() is: 5)
		Assert(rect1.GetY() is: 45)
		Assert(rect1.GetWidth() is: 105)
		Assert(rect1.GetHeight() is: 85)

		Assert(rect1.Overlaps?(rect1))

		rect2 = Rect(200, 67, 654, 123)
		Assert(rect2.ToWindowsRect()
			is: Object(left: 200, top: 67, right: 200 + 654, bottom: 67 + 123))
		.testBoth(rect1, rect2, false)

		rect3 = Rect(0, 0, 10, 10)
		.testBoth(rect3, Rect(-5, -5, 5, 6), true)
		.testBoth(rect3, Rect(-5, -5, 1, 1), false)
		.testBoth(rect3, Rect(5, 5, 7, 5), true)
		.testBoth(rect3, Rect(5, 5, 1, 1), true)
		.testBoth(rect3, Rect(8, 8, 2, 2), true)
		.testBoth(rect3, Rect(6, 6, 3, 6), true)

		rect3.Set(10, 15, 20, 25)
		Assert(rect3.GetX() is: 10)
		Assert(rect3.GetY() is: 15)
		Assert(rect3.GetWidth() is: 20)
		Assert(rect3.GetHeight() is: 25)

		rect3.Set(x: 40)
		Assert(rect3.GetX() is: 40)
		Assert(rect3.GetY() is: 15)
		Assert(rect3.GetWidth() is: 20)
		Assert(rect3.GetHeight() is: 25)

		rect3.Set(y: 90)
		Assert(rect3.GetX() is: 40)
		Assert(rect3.GetY() is: 90)
		Assert(rect3.GetWidth() is: 20)
		Assert(rect3.GetHeight() is: 25)

		rect3.Set(width: 67)
		Assert(rect3.GetX() is: 40)
		Assert(rect3.GetY() is: 90)
		Assert(rect3.GetWidth() is: 67)
		Assert(rect3.GetHeight() is: 25)

		rect3.Set(height: 30)
		Assert(rect3.GetX() is: 40)
		Assert(rect3.GetY() is: 90)
		Assert(rect3.GetWidth() is: 67)
		Assert(rect3.GetHeight() is: 30)

		rect3.SetX(10)
		Assert(rect3.GetX() is: 10)
		Assert(rect3.GetY() is: 90)
		Assert(rect3.GetWidth() is: 67)
		Assert(rect3.GetHeight() is: 30)

		rect3.SetWidth(2 * rect3.GetWidth())
		Assert(rect3.GetX() is: 10)
		Assert(rect3.GetY() is: 90)
		Assert(rect3.GetWidth() is: 134)
		Assert(rect3.GetHeight() is: 30)

		rect3.SetY(rect3.GetY() / 2)
		Assert(rect3.GetX() is: 10)
		Assert(rect3.GetY() is: 45)
		Assert(rect3.GetWidth() is: 134)
		Assert(rect3.GetHeight() is: 30)

		rect3.SetHeight(rect3.GetHeight() - 1)
		Assert(rect3.GetX() is: 10)
		Assert(rect3.GetY() is: 45)
		Assert(rect3.GetWidth() is: 134)
		Assert(rect3.GetHeight() is: 29)
		}
	testBoth(rect1, rect2, overlap)
		{
		Assert(rect1.Overlaps?(rect2) is: overlap)
		Assert(rect2.Overlaps?(rect1) is: overlap)
		}
	Test_Translate()
		{
		rc = Rect.FromWindowsRect(Object(left: 1, top: 1, right: 3, bottom: 3))
		rc.TranslateX(-1)
		Assert(rc.GetX() is: 0)
		Assert(rc.GetY() is: 1)
		Assert(rc.GetWidth() is: 2)
		Assert(rc.GetHeight() is: 2)
		rc.TranslateY(1)
		Assert(rc.GetX() is: 0)
		Assert(rc.GetY() is: 2)
		Assert(rc.GetWidth() is: 2)
		Assert(rc.GetHeight() is: 2)
		rc.Translate(-1, -2)
		Assert(rc.GetX() is: -1)
		Assert(rc.GetY() is: 0)
		Assert(rc.GetWidth() is: 2)
		Assert(rc.GetHeight() is: 2)
		}
	Test_ToWindowsRect()
		{
		rc = Rect(0, 1, 2, 3)
		rcwin = rc.ToWindowsRect()
		Assert(rcwin is: #(left:0, top:1, right:2, bottom:4))
		rc.Set(x:-1, height:10)
		rc.IntoWindowsRect(rcwin)
		Assert(rcwin is: #(left:-1, top:1, right:1, bottom: 11))
		}
	Test_WindowsRectOverlap()
		{
		// THESE ARE THE RECTANGLES
		// ----------------------------------------------
		//     +----------------+    +--------------+
		//     |      RC_TL     |    |     RC_TR    |
		//     +---------------XXXXXXX--------------+
		//                     XXXXXXX
		//                     XXXXXXX------------+
		//     +----------------+   |    RC_BR    |
		//     |      RC_BL     |   +-------------+
		//     +----------------+
		// ----------------------------------------------
		rc_tl = Rect(10, 10, 20, 2)
		rc_tr = Rect(35, 10, 15, 2)
		rc_bl = Rect(10, 20, 20, 2)
		rc_br = Rect(34, 19, 15, 2)
		rc_mid = Rect(28, 11, 7, 8)
		// Preliminary: make sure none of the corner rectangles overlap each
		//				other (except they should overlap themselves)
		cornerRects = Object(rc_tl, rc_tr, rc_bl, rc_br)
		for (rc1 in cornerRects)
			{
			for (rc2 in cornerRects)
				{
				if rc1 is rc2
					{
					Assert(rc1.Overlaps?(rc2) is: true)
					}
				else
					{
					Assert(rc1.Overlaps?(rc2) is: false)
					Assert(rc2.Overlaps?(rc1) is: false)
					}
				}
			}
		// Test the middle rectangle to see whether it overlaps as expected.
		expectedOverlaps = #(true, true, false, true)
		for (k = 0; k < 4; ++k)
			{
			Assert(expectedOverlaps[k] is: cornerRects[k].Overlaps?(rc_mid))
			Assert(expectedOverlaps[k] is: rc_mid.Overlaps?(cornerRects[k]))
			}
		// Now, let's convert the middle rectangle to a Windows RECT structure
		// and check to make sure the results are as expected
		winrc_mid = rc_mid.ToWindowsRect()
		Assert(rc_mid.OverlapsWindowsRect?(winrc_mid))
		for (k = 0; k < 4; ++k)
			{
			Assert(expectedOverlaps[k] is: cornerRects[k].OverlapsWindowsRect?(winrc_mid))
			Assert(Rect.FromWindowsRect(winrc_mid).Overlaps?(cornerRects[k])
				is: expectedOverlaps[k])
			}
		}
	Test_ContainsPoint()
		{
		// General testing
		rc = Rect(0, 0, 10, 10)
		pts = Object(
			Object(new Point(0, 0), true),
			Object(new Point(0, 1), true),
			Object(new Point(1, 0), true),
			Object(new Point(5, 5), true),
			Object(new Point(0, 10), true),
			Object(new Point(10, 0), true),
			Object(new Point(10, 10), true),
			Object(new Point(-1, -1), false)
			Object(new Point(-1, 0), false),
			Object(new Point(0, -1), false),
			Object(new Point(11, 11), false),
			Object(new Point(11, 0), false),
			Object(new Point(0, 11), false),
			Object(new Point(10, 11), false),
			Object(new Point(11, 10), false)
		)
		for (tuple in pts)
			{
			Assert(tuple[1] is: rc.ContainsPoint?(tuple[0]))
			}
		}
	Test_TrapPoint()
		{
		rc = Rect(-1, -1, 2, 1)
		pts = Object(
			Object(new Point(-1, -1)),
			Object(new Point(-1, 0)),
			Object(new Point(1, 0)),
			Object(new Point(1, -1)),
			Object(new Point(0, 0)),
			Object(new Point(-0.5, -0.5)),
			// bottom left corner
			Object(new Point(-2, -2), new Point(-1, -1)),
			Object(new Point(-2, -1), new Point(-1, -1)),
			Object(new Point(-1, -2), new Point(-1, -1))
			// top left corner
			Object(new Point(-2, 0), new Point(-1, 0)),
			Object(new Point(-1, 1), new Point(-1, 0)),
			Object(new Point(-2, 1), new Point(-1, 0)),
			// top right corner
			Object(new Point(2, 1), new Point(1, 0)),
			Object(new Point(2, 0), new Point(1, 0)),
			Object(new Point(1, 1), new Point(1, 0)),
			// bottom right corner
			Object(new Point(2, -1), new Point(1, -1)),
			Object(new Point(2, -2), new Point(1, -1)),
			Object(new Point(1, -2), new Point(1, -1)),
			// right edge
			Object(new Point(2, -0.5), new Point(1, -0.5))
		)
		for (tuple in pts)
			{
			input    = tuple[0]
			expected = 1 < tuple.Size() ? tuple[1] : input
			output   = rc.TrapPoint(input)
			Assert(expected is: output)
			Assert(rc.ContainsPoint?(output))
			}
		}

	Test_Set()
		{
		r = Rect(1,2,3,4)
		Assert({ r.Set(x: 'test') } throws: 'Assert FAILED')
		Assert({ r.Set(y: true) } throws: 'Assert FAILED')
		Assert({ r.Set(width: #20190101) } throws: 'Assert FAILED')
		Assert({ r.Set(height: '50') } throws: 'Assert FAILED')

		r.Set(10, 20, 100, 200)
		Assert(r is: Rect(10, 20, 100, 200))

		r.Set(height: 250)
		Assert(r is: Rect(10, 20, 100, 250))

		r.Set(height: 450, width: 150, y: 5, x: 1)
		Assert(r is: Rect(1, 5, 150, 450))
		}
	}
