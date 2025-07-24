// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.reset()
		}

	reset()
		{
		CircuitBreaker.Reset('test')
		for(i = 1; i <= 6; i++)
			CircuitBreaker.Reset('test' $ i)
		CircuitBreaker.Reset('test@1')
		CircuitBreaker.Reset('test@2')
		}

	getCircuitBreaker(suffix)
		{
		return ServerSuneido.Get('CircuitBreaker_' $ suffix)
		}
	Test_states()
		{
		cl = CircuitBreaker
			{
			HalfOpen?: false
			CircuitBreaker_halfOpen?()
				{
				return .HalfOpen?
				}
			CircuitBreaker_log(@unused) { }
			}
		result = cl('test', { true })
		Assert(result)
		Assert(.getCircuitBreaker('test').state is: 'closed')
		result = cl('test', { true })
		Assert(result)
		Assert(.getCircuitBreaker('test').state is: 'closed')
		result = cl('test', { false })
		Assert(result is: false)
		Assert(.getCircuitBreaker('test').state is: 'closed')
		for .. 8
			cl('test', { false })
		Assert(.getCircuitBreaker('test').state is: 'closed')
		result = cl('test', { false })
		Assert(result is: false)
		Assert(.getCircuitBreaker('test').state is: 'open')


		result = cl('test', { throw "should not be called" })
		Assert(result is: false)
		cb = .getCircuitBreaker('test')
		Assert(cb.state is: 'open')

		cl = CircuitBreaker
			{
			HalfOpen?: true
			CircuitBreaker_halfOpen?()
				{
				return .HalfOpen?
				}
			CircuitBreaker_log(@unused) { }
			}
		result = cl('test', { false })
		Assert(result is: false)
		Assert(.getCircuitBreaker('test').state is: 'open')

		result = cl('test', { true })
		Assert(result)
		Assert(.getCircuitBreaker('test').state is: 'closed')
		}

	Test_delimiter()
		{
		cl = CircuitBreaker { CircuitBreaker_log(@unused) { } }
		for .. 10
			cl('test@1', { false })
		cb = ServerSuneido.Get('CircuitBreaker_test@1')
		Assert(cb isnt: false)
		Assert(cb.state is: 'open')
		cl('test@2', { false })
		cb = ServerSuneido.Get('CircuitBreaker_test@2')
		Assert(cb.state is: 'closed')
		}


	Test_threshold_and_second_instance()
		{
		cl = CircuitBreaker { CircuitBreaker_log(@unused) { } }
		for .. 5
			cl('test2', { false })
		Assert(.getCircuitBreaker('test2').state is: 'closed')
		for .. 5
			cl('test2', { false })
		Assert(.getCircuitBreaker('test2').state is: 'open')
		}

	Test_program_errors_in_block()
		{
		cl = CircuitBreaker { CircuitBreaker_log(@unused) { } }
		for .. 9
			try
				cl('test4', { throw "program error in CircuitBreaker" })
		cb = .getCircuitBreaker('test4')
		Assert(cb.state is: 'closed')
		Assert(cb.failures.Size() is: 9)

		try
			cl('test4', { throw "program error in CircuitBreaker" })
		Assert(.getCircuitBreaker('test4').failures.Size() is: 10)
		Assert(.getCircuitBreaker('test4').state is: 'open')
		}

	Test_timeout_increment()
		{
		cl = CircuitBreaker
			{
			HalfOpen?: false
			CircuitBreaker_getServiceConfig(service /*unused*/)
				{
				return Object('test5', threshold: 1, timeout: 10, timeoutIncrement: 50)
				}
			CircuitBreaker_halfOpen?()
				{
				return .HalfOpen?
				}
			CircuitBreaker_log(@unused) { }
			}
		cl('test5', { false })
		Assert(.getCircuitBreaker('test5').state is: 'open')

		cl('test5', { false })
		// still 10 because it won't increment when circuit breaker is open
		cb = .getCircuitBreaker('test5')
		Assert(cb.timeout is: 10)
		Assert(cb.state is: 'open')


		cl = CircuitBreaker
			{
			HalfOpen?: true
			CircuitBreaker_getServiceConfig(service /*unused*/)
				{
				return Object('test5', threshold: 1, timeout: 10, timeoutIncrement: 50)
				}
			CircuitBreaker_halfOpen?()
				{
				return .HalfOpen?
				}
			CircuitBreaker_log(@unused) { }
			}
		// re-close circuit breaker
		cl('test5', { true })
		cb = .getCircuitBreaker('test5')
		Assert(cb.timeout is: 10)
		Assert(cb.state is: 'closed')

		// back to open
		cl('test5', { false })
		cb = .getCircuitBreaker('test5')
		Assert(cb.timeout is: 10)
		Assert(cb.state is: 'open')

		// half-open with failure.  Time should now increment
		cl('test5', { false })
		cb = .getCircuitBreaker('test5')
		Assert(cb.timeout is: 60)
		Assert(cb.state is: 'open')

		// half-open with another failure.  Should increment again
		cl('test5', { false })
		cb = .getCircuitBreaker('test5')
		Assert(cb.timeout is: 110)
		Assert(cb.state is: 'open')

		// back to open.  timeout should reset
		cl('test5', { true })
		cb = .getCircuitBreaker('test5')
		Assert(cb.timeout is: 10)
		Assert(cb.state is: 'closed')
		}

	Test_failureExpiry()
		{
		cl = CircuitBreaker
			{
			CircuitBreaker_getServiceConfig(service /*unused*/)
				{
				return Object('test6', threshold: 3, timeout: 600, failureExpiry: 50)
				}
			CircuitBreaker_log(@unused) { }
			}
		Assert(CircuitBreaker.AttemptsRemaining('test6') is: 10)
		cl('test6', { true })
		for i in ..2
			{
			cl('test6', { false })
			Assert(cl.AttemptsRemaining('test6') is: 2-i)
			cb = .getCircuitBreaker('test6')
			Assert(cb.failures isSize: i+1)
			Assert(cb.failures[i] isDate:)
			cc = cl.CircuitBreaker_getCbInstance('test6')
			Assert(cc.CircuitBreaker_thresholdExceeded?() is: false)
			}
		cb = .getCircuitBreaker('test6')

		// set first failure to be expired
		cb.failures[0] = cb.failures[0].Minus(minutes: 50)
		ServerSuneido.Set('CircuitBreaker_test6', cb)
		cl('test6', { false })
		cb = .getCircuitBreaker('test6')
		Assert(cb.failures isSize: 2)
		Assert(cl.MinutesUntilRunnable('test6') is: 0)
		Assert(cb.state is: 'closed')
		Assert(cl.AttemptsRemaining('test6') is: 1)
		cc = cl.CircuitBreaker_getCbInstance('test6')
		Assert(cc.CircuitBreaker_thresholdExceeded?() is: false)

		cl('test6', { false })
		cb = .getCircuitBreaker('test6')
		Assert(cb.failures isSize: 3)
		Assert(cl.MinutesUntilRunnable('test6') is: 10)
		Assert(cb.state is: 'open')
		Assert(cl.AttemptsRemaining('test6') is: 0)
		cc = cl.CircuitBreaker_getCbInstance('test6')
		Assert(cc.CircuitBreaker_thresholdExceeded?())

		called = false
		cl('test6', { called = true })
		cb = .getCircuitBreaker('test6')
		Assert(called is: false)
		Assert(cb.failures isSize: 3)
		Assert(cl.MinutesUntilRunnable('test6') is: 10)
		Assert(cb.state is: 'open')
		Assert(cl.AttemptsRemaining('test6') is: 0)
		cc = cl.CircuitBreaker_getCbInstance('test6')
		Assert(cc.CircuitBreaker_thresholdExceeded?())

		cb.lastFailureTime = Date().Minus(minutes: 11)
		ServerSuneido.Set('CircuitBreaker_test6', cb)
		//TODO - test HalfOpen
		Assert(cl.AttemptsRemaining('test6') is: 1)
		cl('test6', { false })
		cb = .getCircuitBreaker('test6')
		Assert(cb.failures isSize: 3)
		Assert(cl.MinutesUntilRunnable('test6') is: 10)
		Assert(cb.state is: 'open')
		Assert(cl.AttemptsRemaining('test6') is: 0)
		cc = cl.CircuitBreaker_getCbInstance('test6')
		Assert(cc.CircuitBreaker_thresholdExceeded?())

		cb.lastFailureTime = Date().Minus(minutes: 11)
		ServerSuneido.Set('CircuitBreaker_test6', cb)
		//TODO - test HalfOpen
		cl('test6', { true })
		cb = .getCircuitBreaker('test6')
		Assert(cb.failures isSize: 0)
		Assert(cl.MinutesUntilRunnable('test6') is: 0)
		Assert(cl.AttemptsRemaining('test6') is: 3)
		cc = cl.CircuitBreaker_getCbInstance('test6')
		Assert(cc.CircuitBreaker_thresholdExceeded?() is: false)
		}

	Teardown()
		{
		.reset()
		ServerSuneido.DeleteMember('CircuitBreaker_test')
		ServerSuneido.DeleteMember('CircuitBreaker_test2')
		ServerSuneido.DeleteMember('CircuitBreaker_test4')
		ServerSuneido.DeleteMember('CircuitBreaker_test5')
		ServerSuneido.DeleteMember('CircuitBreaker_test6')
		super.Teardown()
		}
	}
