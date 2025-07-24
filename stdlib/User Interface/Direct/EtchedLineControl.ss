// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Hwnd
	{
	Xstretch: 	0
	Name:		"EtchedLine"
	New(before = 2, after = 2)
		{
		.before = Max(0, ScaleWithDpiFactor(before))
		.after = Max(0, ScaleWithDpiFactor(after))
		.CreateWindow("static", "", WS.VISIBLE | SS.ETCHEDHORZ)
		.Ymin = 2 + .before + .after
		}
	GetReadOnly()			// read-only not applicable to etchedline
		{ return true }
	Resize(x, y, w, h)
		{
		super.Resize(x, y + .before, w, 2)
		}
	}
