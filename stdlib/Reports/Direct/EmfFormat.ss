// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	New(.file, w, h)
		{
		.w = w.InchesInTwips()
		.h = h.InchesInTwips()
		.hemf = GetEnhMetaFile(file)
		if .hemf is NULL
			throw "EmfFormat: can't open: " $ file
		}
	GetSize(data /*unused*/ = false)
		{
		return Object(w: .w, h: .w, d: 0)
		}
	Print(x, y, w /*unused*/, h /*unused*/, data /*unused*/ = false)
		{
		hdc = _report.GetDC()
		if not PlayEnhMetaFile(hdc, .hemf, Object(left: x, top: y, right: .w, bottom: .h))
			throw "EmfFormat: error displaying: " $ .file
		}
	}