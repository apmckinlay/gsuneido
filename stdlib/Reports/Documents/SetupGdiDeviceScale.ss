// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function (dc)
	{
	// set scale to 1440 units per inch
	SetMapMode(dc, MM.ISOTROPIC)
	tenInchTwips = 10.InchesInTwips()
	tenInchMM = 254
	SetWindowExt(dc,
		tenInchTwips * GetDeviceCaps(dc, GDC.HORZSIZE) / tenInchMM,
		tenInchTwips * GetDeviceCaps(dc, GDC.VERTSIZE) / tenInchMM)
	SetViewportExt(dc,
		GetDeviceCaps(dc, GDC.HORZRES),
		GetDeviceCaps(dc, GDC.VERTRES))
	}