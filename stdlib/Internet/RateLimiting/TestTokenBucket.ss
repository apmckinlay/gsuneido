// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		/* PURPOSE: to test full functionality of TokenBucket
			- this requires letting time pass so cannot run in test framework
			- attempts to ensure thread safe
			- feel free to make improvements or additions
		*/

		tb = TokenBucket(capacity: 5, refillRatePerSec: 1)
		for unused in ..30 /* = number of concurrent calls */
			{
			Thread({
				Thread.Sleep(100) /* = tenth second delay to simulate network request */
				tb.Consume()
			})
			}
		Thread.Sleep(1000) /* = small delay to ensure all threads complete */

		// all tokens should be used by now
		Assert(tb.TokenBucket_tokens is: 0)

		Thread.Sleep(6000) /*= if time passes and no Consume calls, tokens not adjusted */
		Assert(tb.TokenBucket_tokens is: 0)

		// at this point consume should refill full bucket to 5, then use 1
		tb.Consume()
		Assert(tb.TokenBucket_tokens is: 4)

		Assert(tb.Consume())
		Assert(tb.Consume())

		Assert(tb.TokenBucket_tokens is: 2)

		Assert(tb.Consume())
		Assert(tb.Consume())
		Assert(tb.Consume() is: false) // out of tokens
		Assert(tb.TokenBucket_tokens is: 0)

		// another burst of requests with no tokens available, all should fail
		for unused in ..30 /* = number of concurrent calls */
			{
			// shorter sleep this time to ensure we don't "earn" any tokens
			Thread({
				Thread.Sleep(20) /* = tenth second delay to simulate network request */
				Assert(tb.Consume() is: false)
			})
			}
		Thread.Sleep(1000) /* = small delay to ensure all threads complete */
		Assert(tb.TokenBucket_tokens is: 0)

		Print("TestTokenBucket completed successfully")
		}

	TestConcurrency()
		{
		// WARNING: not guarunteed to fail each time when Mutex is disabled
		// 100 tokens maximum, 0 refill rate.
		// This guarantees the bucket should ONLY ever hand out 100 tokens.
		maxTokens = 100
		numThreads = 10
		Suneido.TestTokenBucketData = Object(bucket: TokenBucket(maxTokens, 0)
			results: Object(), threadsCompleted: 0, testMutex: Mutex())

		for (i = 0; i < numThreads; ++i)
			{
			Thread(function()
				{
				attemptsPerThread = 20 // 10 threads * 20 attempts = 200 total requests
				successes = 0

				// Hammer the bucket
				for (j = 0; j < attemptsPerThread; ++j)
					{
					if Suneido.TestTokenBucketData.bucket.Consume(1)
						successes++
					}
				Suneido.TestTokenBucketData.testMutex.Do()
					{
					Suneido.TestTokenBucketData.results.Add(successes)
					Suneido.TestTokenBucketData.threadsCompleted++
					}
				})
			}

		// Wait for all threads to finish their loops
		cycles = 0
		while Suneido.TestTokenBucketData.threadsCompleted < numThreads
			{
			Thread.Sleep(10) /* = Wait 10ms before checking again */
			if cycles++ > 500 /* = max time of 5 seconds */
				{
				Print("TOO MANY CYCLES - THREADS NOT FINISHING!")
				return
				}
			}

		// Tally the total tokens successfully handed out
		totalConsumed = 0
		for res in Suneido.TestTokenBucketData.results
			totalConsumed += res

		if totalConsumed isnt maxTokens
			Print("TEST FAILED: " $
				"Consumed " $ totalConsumed $ " tokens - should be " $ maxTokens)
		else
			Print("TEST PASSED")
		}
	}