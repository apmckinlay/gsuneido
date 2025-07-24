// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function (msg, params = '', caughtMsg = '')
	{
	if Suneido.User is 'default'
		throw msg
	SuneidoLog('ERROR: (CAUGHT) ' $ msg, calls:, :params, :caughtMsg)
	ProgrammerErrorExtra(msg)
	}