// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
// block must throw an exception to fail (can't just return false)
// e.g. Retry(){ if not YesNo("ok") throw "boom" }
// if retryException is passed in, it must match e exactly in order to force retry
// retryException can be a string or an Object of exceptions with optional minDelayMs:
//   e.g. Object('Error1', Object('Error2', minDelayMs: 500))
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
				delay = minDelayMs
				if retryException isnt '' and
					false is delay = .matchRetryException(e, retryException, minDelayMs)
					throw e
				if i isnt maxRetries - 1
					RetrySleep(i, delay)
				}
		throw "Retry failed - too many retries, last error: " $ e
		}

	matchRetryException(e, retryException, defaultMinDelayMs)
		{
		if not Object?(retryException)
			return e is retryException
				? defaultMinDelayMs
				: false

		for exception in retryException
			{
			if Object?(exception)
				{
				if e is exception[0]
					return exception.GetDefault('minDelayMs', defaultMinDelayMs)
				}
			else if e is exception
				return defaultMinDelayMs
			}
		return false
		}
	}
