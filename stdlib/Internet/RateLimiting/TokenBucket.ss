// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// Usage: call the Consume method with the number of tokens to use, if the
	// tokens are available it will return true, otherwise false will be returned
	// and the API should return error code 429 (too many requests)
	// capacity - The maximum burst size allowed
	// refillRatePerSec - How many tokens are added per second (can be a fraction)
	New(.capacity, .refillRatePerSec)
		{
		// Start with a full bucket
		.tokens = capacity
		.lastRefilled = Date()
		.mutex = Mutex()
		}

	Consume(tokensRequested = 1)
		{
		result = false
		.mutex.Do()
			{
			.refill()
			if .tokens >= tokensRequested
				{
				.tokens -= tokensRequested
				result = true
				}
			else
				result = false
			}
		return result
		}

	refill()
		{
		// Lazily updates the token count based on time passed (called by Consume method)
		now = .now()
		timePassedInSeconds = now.MinusSeconds(.lastRefilled).Floor()
		// Calculate how many tokens were earned in that time span
		tokensEarned = timePassedInSeconds * .refillRatePerSec

		if tokensEarned > 0
			{
			// Add the new tokens, but cap it at the bucket's max capacity
			.tokens = Min(.capacity, .tokens + tokensEarned)
			.lastRefilled = now
			}
		}
	now()
		{
		return Date()
		}
	}