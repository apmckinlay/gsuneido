// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function (hwnd, msg, title, flags)
	{
	SuneidoLog("ERROR: ALERT - " $ msg, calls:, params: Object(:title, :hwnd, :flags))
	return false // prevent returning suneidolog record
	}
