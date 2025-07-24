// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
function (prefix, block, asErratic = #())
	{
	try
		block()
	catch (e)
		{
		level = "ERROR"
		if asErratic.Any?({ e =~ it })
			level = "ERRATIC"
		SuneidoLog(level $ ": (CAUGHT) " $ prefix $ ': ' $ e,
			caughtMsg: 'unattended; no msg for user; multiple possible sources')
		}
	}