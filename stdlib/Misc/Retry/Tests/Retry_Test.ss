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
	}
