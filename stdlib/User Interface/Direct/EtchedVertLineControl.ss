// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Hwnd
	{
	Ystretch: 	1
	Name:		"EtchedVertLine"
	New(before = 2, after = 2)
		{
		.before = Max(0, before.Int())
		.CreateWindow("static", "", WS.VISIBLE | SS.ETCHEDVERT)
		.Xmin = 2 + .before + Max(0, after.Int())
		}
	GetReadOnly()			// read-only not applicable to etchedline
		{ return true }
	Resize(x, y, w, h)
		{
		super.Resize(x + .before, y, 2, h)
		}
	}
