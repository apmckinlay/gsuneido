// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	// NOTE: see TestTokenBucket for functional testing
	Test_Consume()
		{
		testCl = TokenBucket
			{
			TokenBucket_refill() {}
			}
		cl = testCl(capacity: 5, refillRatePerSec: 1)
		Assert(cl.Consume())
		Assert(cl.Consume())
		Assert(cl.Consume())
		Assert(cl.Consume())
		Assert(cl.Consume())
		Assert(cl.Consume() is: false)
		Assert(cl.TokenBucket_tokens is: 0)
		}

	Test_refill()
		{
		_now = Date()
		testCl = TokenBucket
			{
			TokenBucket_now() { return _now }
			}
		cl = testCl(capacity: 5, refillRatePerSec: 1)

		// refill on empty bucket but not enough time has passed to earn anything
		cl.TokenBucket_tokens = 0
		cl.TokenBucket_lastRefilled = _now.Minus(milliseconds: 111)
		cl.TokenBucket_refill()
		Assert(cl.TokenBucket_tokens is: 0)

		// bucket empty, enough time passed to earn 1 token
		cl.TokenBucket_lastRefilled = _now.Minus(milliseconds: 1500)
		cl.TokenBucket_refill()
		Assert(cl.TokenBucket_tokens is: 1)

		// enough time passed to fill bucket completely, but should not over-fill
		cl.TokenBucket_lastRefilled = _now.Minus(seconds: 10)
		cl.TokenBucket_refill()
		Assert(cl.TokenBucket_tokens is: 5)

		// refill after enough time to earn, bucket should just remain full
		cl.TokenBucket_lastRefilled = _now.Minus(seconds: 30)
		cl.TokenBucket_refill()
		Assert(cl.TokenBucket_tokens is: 5)
		}
	}