// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (s)
	{
	if Suneido.User is 'default'
		Alert(s, "WARNING", flags: MB.ICONWARNING)
	else
		SuneidoLog("warning: " $ s)
	}