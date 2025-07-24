// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
function (seconds, block)
	{
	forever
		{
		start = Date()
		block()
		duration = Date().MinusSeconds(start)
		if duration < seconds
			Thread.Sleep((seconds - duration) * 1.SecondsInMs())
		}
	}