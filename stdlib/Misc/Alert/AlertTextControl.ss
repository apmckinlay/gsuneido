// Copyright (C) 2016 Axon Development Corporation All rights reserved worldwide.
// specifically written for AlertMessageBox
// could be generalized
StaticControl
	{
	Name: 'AlertText'
	New(.text)
		{
		super(.text, whitebgnd:)
		.flags = DT.WORDBREAK | DT.EDITCONTROL | DT.NOPREFIX | DT.EXPANDTABS
		.setSize()
		.AdjustControlSize()
		}
	setSize()
		{
		xmin = ScaleWithDpiFactor(200) /*= xmin*/
		xmax = ScaleWithDpiFactor(800) /*= xmax*/
		yThreshold = ScaleWithDpiFactor(400) /*= yThreshold*/
		.WithDC()
			{|dc|
			DoWithHdcObjects(dc, [.GetFont()])
				{
				DrawText(dc, .text, -1, r = [right: xmax / 2], .flags | DT.CALCRECT)
				if r.bottom > yThreshold
					DrawText(dc, .text, -1, r = [right: xmax], .flags | DT.CALCRECT)
				}
			}
		.Xmin = Max(r.right, xmin)
		.Ymin = r.bottom
		}
	ERASEBKGND()
		{ return 1 }
	}