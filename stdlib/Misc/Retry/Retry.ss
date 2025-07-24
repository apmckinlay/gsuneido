// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
// block must throw an exception to fail (can't just return false)
// e.g. Retry(){ if not YesNo("ok") throw "boom" }
// if retryException is passed in, it must match e exactly in order to force retry
class
	{
	CallClass(block, maxRetries = 10, minDelayMs = 2, retryException = '')
		{
		e = ''
		for i in .. maxRetries
			try
				{
				e = ''
				return block()
				}
			catch (e)
				{
				if retryException isnt '' and e isnt retryException
					throw e
				if i isnt maxRetries - 1
					RetrySleep(i, minDelayMs)
				}
		throw "Retry failed - too many retries, last error: " $ e
		}
	}
