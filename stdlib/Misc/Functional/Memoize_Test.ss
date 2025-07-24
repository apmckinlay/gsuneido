// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()
		Suneido.Delete('Memoize_Memoize_Test.Test_Expire memoizeClass')
		Suneido.Delete('Memoize_Memoize_Test.Test_Expire memoizeClassExpired')
		}

	Test_Cache()
		{
		memoizeClass = Memoize
			{
			ExpirySeconds: 3600
			Func(callTracker, a /*unused*/)
				{ return ++callTracker.calls }
			}
		memoizeClass.ResetCache()
		callTracker = Object(calls: 0)
		numberOfCalls = memoizeClass(callTracker, 'first param set')
		Assert(numberOfCalls is: 1)
		numberOfCalls = memoizeClass(callTracker, 'first param set')
		Assert(numberOfCalls is: 1)
		numberOfCalls = memoizeClass(callTracker, 'second param set')
		Assert(numberOfCalls is: 2)
		memoizeClass.ResetCache()
		numberOfCalls = memoizeClass(callTracker, 'first param set')
		Assert(numberOfCalls is: 3)
		memoizeClass.ResetCache()
		}

	Test_Expiry()
		{
		memoizeClass =	Memoize
			{
			ExpirySeconds: 3600
			Expires?: true
			Func(callTracker, a /*unused*/)
				{ return ++callTracker.calls }
			}
		memoizeClass.ResetCache()
		callTracker = Object(calls: 0)
		numberOfCalls = memoizeClass(callTracker, 'first param set')
		Assert(numberOfCalls is: 1)
		numberOfCalls = memoizeClass(callTracker, 'first param set')
		Assert(numberOfCalls is: 1)
		numberOfCalls = memoizeClass(callTracker, 'second param set')
		Assert(numberOfCalls is: 2, msg: 'first')
		memoizeClass.ResetCache()
		numberOfCalls = memoizeClass(callTracker, 'first param set')
		Assert(numberOfCalls is: 3)
		}

	Test_Expired()
		{
		memoizeClassExpired = Memoize
			{
			ExpirySeconds: -3600
			Expires?: true
			Func(callTracker, a /*unused*/)
				{ return ++callTracker.calls }
			}
		Suneido.Delete('Memoize_Memoize_Test.Test_Expire$c')
		memoizeClassExpired.ResetCache()
		callTracker = Object(calls: 0)
		numberOfCalls = memoizeClassExpired(callTracker, 'first param set')
		Assert(numberOfCalls is: 2, msg: 'second') // two because it expired immediately
		numberOfCalls = memoizeClassExpired(callTracker, 'first param set')
		Assert(numberOfCalls is: 3)
		memoizeClassExpired.ResetCache()
		}
	}