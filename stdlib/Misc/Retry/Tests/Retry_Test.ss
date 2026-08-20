// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Main()
		{
		count = 0
		result = Retry(2, 1)
			{
			count++
			'hello world'
			}
		Assert(result is: 'hello world')
		Assert(count is: 1)

		count = 0
		block = {
			if count++ < 2
				throw 'ERROR!'
			'this is the result'
			}
		Assert({Retry(block, 2, 1)}
			throws: 'Retry failed - too many retries, last error: ERROR!')
		Assert(count is: 2)

		count = 0
		Assert(Retry(block, 3, 1) is: 'this is the result')
		Assert(count is: 3)
		}

	Test_MultiRetryException()
		{
		// object of string exceptions
		count = 0
		block = {
			if count++ < 2
				throw 'ERR_A'
			'ok'
			}
		Assert(Retry(block, 3, 1,
			retryException: Object('ERR_A')) is: 'ok')
		Assert(count is: 3)

		// unmatched exception should throw
		count = 0
		block = {
			count++
			throw 'ERR_C'
			}
		Assert({Retry(block, 3, 1,
			retryException: Object('ERR_A'))}
			throws: 'ERR_C')
		Assert(count is: 1)

		// multiple exceptions
		count = 0
		block = {
			count++
			if count is 1
				throw 'ERR_A'
			else if count is 2
				throw 'ERR_B'
			'ok'
			}
		Assert(Retry(block, 3, 1,
			retryException: Object('ERR_A', 'ERR_B')) is: 'ok')
		Assert(count is: 3)

		// mixed: plain string and object with minDelayMs
		count = 0
		block = {
			count++
			if count is 1
				throw 'ERR_A'
			else if count is 2
				throw 'ERR_B'
			'ok'
			}
		Assert(Retry(block, 3, 1,
			retryException: Object('ERR_A',
				Object('ERR_B', minDelayMs: 1))) is: 'ok')

		// string retryException still works (backward compat)
		count = 0
		block = {
			if count++ < 2
				throw 'SPECIFIC'
			'ok'
			}
		Assert(Retry(block, 3, 1,
			retryException: 'SPECIFIC') is: 'ok')
		Assert(count is: 3)
		}

	Test_MatchRetryException()
		{
		fn = Retry.Retry_matchRetryException
		// plain string retryException: match
		Assert(fn('ERR_A', 'ERR_A', 5) is: 5)

		// plain string retryException: no match
		Assert(fn('ERR_A', 'ERR_B', 5) is: false)

		// Object with plain member: match
		Assert(fn('ERR_A', #('ERR_A'), 5) is: 5)

		// Object with plain member: no match
		Assert(fn('ERR_A', #('ERR_B'), 5) is: false)

		// nested Object with minDelayMs: returns custom delay
		Assert(fn('ERR_B', #(('ERR_B', minDelayMs: 500)), 5) is: 500)

		// nested Object without minDelayMs: falls back to default
		Assert(fn('ERR_B', #(('ERR_B')), 5) is: 5)

		// empty Object: no match
		Assert(fn('ERR_A', #(), 5) is: false)

		// mixed: match plain string member
		Assert(fn('ERR_A', #('ERR_A', ('ERR_B', minDelayMs: 500)), 5) is: 5)

		// mixed: match nested Object member
		Assert(fn('ERR_B', #('ERR_A', ('ERR_B', minDelayMs: 500)), 5) is: 500)
		}
	}
