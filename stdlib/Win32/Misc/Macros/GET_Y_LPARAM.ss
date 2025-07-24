// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function(x)
	{
	// defined in windowsx.h as #define GET_Y_LPARAM(lp)  ((int)(short)HIWORD(lp))
	return HISWORD(x)
	}
