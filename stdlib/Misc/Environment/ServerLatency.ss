// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// NOTE: running this locally will be so fast it'll often give zero
// will also depend on how busy the suneido server is

// returns the average of the lesser of the two types of request
function ()
	{
	total = 0
	reps = 5
	reps.Times
		{
		// use Database.CurrentSize() for a simple request to db server
		t1 = Timer({ Database.CurrentSize() })
		t2 = 10
		// use Built as a simple request that doesn't use database
		try t2 = Timer({ RunOnHttp(HttpPort(), 'Built', timeoutConnect: .1) })
		total += Min(t1, t2)
		}
	return (total.SecondsInMs() / reps).Round(0)
	}