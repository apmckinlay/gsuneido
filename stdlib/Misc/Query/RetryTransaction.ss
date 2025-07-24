// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
class
	{
	maxRetries: 10
	CallClass(block)
		{
		e = ""
		for (i = 0; i < .maxRetries; ++i)
			try
				{
				Transaction(update:)
					{|t|
					block(t)
					// forceRetryTransaction option is used to test code to make sure
					// it is safe to run more than once when RetryTransaction is used.
					// Simulates a transaction failure.
					if Suneido.GetDefault("forceTooManyRetryTransaction", false) or
						(i is 0 and Suneido.GetDefault("forceRetryTransaction", false))
						{
						t.Rollback()
						throw 'Transaction: block commit failed'
						}
					}
				return // succeeded
				}
			catch (e)
				{
				if not .IsTransactionConflict?(e)
					throw e
				// else try again after delay, unless it is the last retry
				if i isnt .maxRetries - 1
					RetrySleep(i, min: 5) /*= so max delay is long enough
						to hopefully handle conflicting with a long transaction */
				}
		throw "RetryTransaction: too many retries, last error: " $ e
		}

	IsTransactionConflict?(e)
		{
		return e.Prefix?('Transaction: block commit failed') or
			e.Has?('transaction conflict') or
			e.Has?('preempted by exclusive') or
			e.Has?('conflict with exclusive')
		}
	}