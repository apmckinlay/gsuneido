// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_calcPos()
		{
		posRect = Object(left: 100, right: 190, top: 300, bottom: 350)
		wa = #(left: 0, right: 1000, top: 0, bottom: 1000)
		width = 70
		height = 240
		f = { Dialog.Dialog_calcPos(posRect, width, height, wa) }

		Assert(f() is: #(x: 100, y: 350)) // below

		width = 950
		Assert(f() is: #(x: 50, y: 350)) // below, pushed left
		width = 70

		posRect.top = 900
		posRect.bottom = 950
		Assert(f() is: #(x: 100, y: 660)) // above

		posRect.top = 500
		posRect.bottom = 550
		height = 800
		Assert(f() is: #(x: 190, y: 150)) // to the right

		posRect.left = 900
		posRect.right = 990
		Assert(f() is: #(x: 830, y: 150)) // to the left
		}
	}