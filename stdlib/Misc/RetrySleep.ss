// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// i is the retry count (first try would be zero)
// min is the minimum number of milliseconds to sleep, override to make larger
function (i, min = 2)
	{
	// random exponential fallback
	Assert(min > 0)
	if Suneido.GetDefault(#SkipRetrySleep, false) // for tests
		return
	n = min << i
	Thread.Sleep(n + Random(n))
	}