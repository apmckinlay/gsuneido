// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// this is the server side of the workstation's long poll
	// this method itself is not long polling (it's just straight polling)
	// with current settings, poll duration will be ~ 40 seconds (iterations * wait)
	pollIterations: 400
	pollWaitMs: 99
	CallClass(args, fn, conditionFn defaultFn)
		{
		// TODO: change to loop until 40 seconds has elapsed, not a fixed number
		checkingTime = Date()
		.pollIterations.Times
			{
			result = fn(args, checkingTime)
			if conditionFn(result)
				return Json.Encode(result)
			Thread.Sleep(.pollWaitMs)
			}
		return Json.Encode(defaultFn())
		}
	}