// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
function (block)
	{
	Thread()
		{
		reps = 0
		p = Thread.Profile({ reps = Timer.Secs(secs: 1, :block).reps })
		Defer({ ProfileResults(p, reps) })
		}
	}
