// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
GdiDriver
	{
	New(.dc)
		{
		.setMapping()
		di = Object(
			cbSize: DOCINFO.Size(),
			lpszDocName: "Report")
		// Cannot switch to WithBkMode as the end point is ambiguous
		SetBkMode(.dc, TRANSPARENT)
		StartDoc(.dc, di)
		}

	setMapping()
		{
		SetupGdiDeviceScale(.dc)
		SetViewportOrg(.dc,
			-GetDeviceCaps(.dc, GDC.PHYSICALOFFSETX),
			-GetDeviceCaps(.dc, GDC.PHYSICALOFFSETY))
		}

	AddPage(dimens /*unused*/)
		{
		StartPage(.dc)
		.setMapping()
		}

	EndPage()
		{
		EndPage(.dc)
		}

	RightJustify(rect /*unused*/, data /*unused*/, flags)
		{
		flags |= DT.RIGHT
		return flags
		}

	CenterJustify(rect /*unused*/, data /*unused*/, flags)
		{
		flags |= DT.CENTER
		return flags
		}

	Finish(status)
		{
		super.Finish(status)
		EndDoc(.dc)
		return status
		}

	GetCopyDC()
		{ return .dc }
	GetDC()
		{ return .dc }
	}
