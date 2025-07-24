// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
// if timing something small, the overhead of calling the block is significant
// could time an empty block and subtract that
class
	{
	CallClass(block, reps = 1, secs = false)
		{
		if secs is false
			return .run_reps(block, reps) // returns seconds
		else
			{
			r = .Secs(block, secs)
			return .Format(r)
			}
		}
	Format(r)
		{
		reps_per_sec = (r.reps / r.elapsed).RoundToPrecision(2)
		reps_per_sec = reps_per_sec >= 1
			? reps_per_sec $ " reps/sec = "
			: ""
		return reps_per_sec $ ReadableDuration(r.elapsed / r.reps) $ "/rep"
		}
	run_reps(block, reps)
		{
		start = .date()
		for ..reps
			block()
		return .date().MinusSeconds(start)
		}
	Secs(block, secs) // returns number of reps
		{
		start = .date()
		reps = 0
		n = 1
		do
			{
			for ..n
				block()
			reps += n
			elapsed = .date().MinusSeconds(start)
			avg = elapsed / reps
			if avg is 0
				n *= 100 /*= increase to get a measurable time */
			else
				{
				remaining = secs - elapsed
				n = (remaining / avg).Round(0)
				}
			} while n >= 1
		return Object(:reps, :elapsed)
		}
	date() // so test can override
		{
		return Date()
		}
	}
