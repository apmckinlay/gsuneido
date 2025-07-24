// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
//
// Note: This is specifically comparing against true
// to handle blocks/functions that return e.g. true or error string e.g. for return throw
// This is still compatible with returning true/false
class
	{
	CallClass(maxretries, min, block)
		{
		Assert(block isnt false)
		result = false
		for (i = 1; true isnt (result = block(count: i)) and i < maxretries; ++i)
			.retrySleep(i, min)

		if String?(result)
			result $= " (too many retries)"
		return result
		}

	retrySleep(i, min) // for test
		{
		RetrySleep(i, min)
		}
	}
