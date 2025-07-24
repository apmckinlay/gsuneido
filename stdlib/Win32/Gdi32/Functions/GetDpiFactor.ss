// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	dpiFactor = 1
	WithDC(NULL)
		{|dc|
		dpiFactor = GetDeviceCaps(dc, GDC.LOGPIXELSY) / WinDefaultDpi
		}
	return dpiFactor
	}