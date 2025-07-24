// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
function (hdc, w, h, block)
	{
	hdcTemp = CreateCompatibleDC(hdc)
	bmpTemp = CreateCompatibleBitmap(hdc, w, h)
	DoWithHdcObjects(hdcTemp, [bmpTemp])
		{
		block(hdcTemp)
		}
	DeleteObject(bmpTemp)
	DeleteDC(hdcTemp)
	}