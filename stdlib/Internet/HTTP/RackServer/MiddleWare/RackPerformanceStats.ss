// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
// Use this as the first middle-ware if we need to track the performance
RackComposeBase
	{
	max: 284804 /*= low boundary of the 19th bucket*/
	New(@args)
		{
		super(@args)
		.ResetStats()
		}

	ResetStats()
		{
		Suneido[.StatsMem] = Object().AddMany!(0, 20) /*=max bucket >= around 2^18 ms*/
		}

	StatsMem: 'RackPerfStats'
	Call(env)
		{
		start = Date()
		if -1 is result = .App(:env)
			return result

		bucket = .getBucket(.getDuration(start))
		Suneido[.StatsMem][bucket]++ // atomic
		return result
		}

	getBucket(duration)
		{
		if 0 is duration
			return 0

		if duration >= .max
			return 19 /*=max bucket*/

		return (duration.Log10() * 3.3/*= base 2 and 10 conversion*/).Ceiling()
		}

	getDuration(start)
		{
		return Date().MinusSeconds(start) * 1000 /*= to millisecond conversion*/
		}
	}