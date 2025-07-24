// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
MapButtonControl
	{
	New(multiLoc = false, onlyLatLong? = false)
		{
		super(multiLoc, onlyLatLong?)
		.Top = .GetChild().Top
		}
	Controls: (Button 'Map' tip: 'Right click to choose map source')
	}