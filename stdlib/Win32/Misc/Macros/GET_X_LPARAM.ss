// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function(x)
	{
	// defined in windowsx.h as #define GET_X_LPARAM(lp)  ((int)(short)LOWORD(lp))
	return LOSWORD(x)
	}