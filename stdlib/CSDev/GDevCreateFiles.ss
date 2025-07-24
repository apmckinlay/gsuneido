// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
class
	{
	ServerBat: 'Start G Server.bat'
	CallClass(quiet? = false, defaultPort = false)
		{
		port = defaultPort
			? 3147 /*= default port */
			: 4000 + Random(5000) /*= from 4000 to 9000,
				to avoid default 3147 and handles http port */
		PutFile('gsuneido.args', '-s -p ' $ port $ ' CSDevServerWindow()')
		if not quiet?
			Print('Created File: gsuneido.args')
		}
	}
