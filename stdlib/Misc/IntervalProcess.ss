// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// Use shortSleepInMS when you want to run again quickly if not all processed
// - return false from callable to indicate not complete and short sleep
// - return true from callable to indicate complete and regular sleep
class
	{
	CallClass(callable, sleepInMS, shortSleepInMS = false,
		skipFunc = function() { false }, circuitBreakerName = false)
		{
		.assertOptions(circuitBreakerName, shortSleepInMS)
		nextSleepMS = sleepInMS
		callable = .wrapCallableIfOnCircuitBreaker(callable, circuitBreakerName)
		extraProcessPause = OptContribution('ExtraProcessPause',
			function () { return false })
		forever
			{
			Thread.Sleep(nextSleepMS)
			skip = skipFunc()
			if skip is true or extraProcessPause()
				continue
			else if skip is 'quit'
				return

			if shortSleepInMS is false
				callable()
			else
				nextSleepMS = callable() ? sleepInMS : shortSleepInMS
			}
		}

	wrapCallableIfOnCircuitBreaker(callable, circuitBreakerName)
		{
		if circuitBreakerName is false
			return callable
		return {
			try
				CircuitBreaker(circuitBreakerName, callable)
			catch (err)
				SuneidoLog(circuitBreakerName $ ' - ' $ err)
			}
		}

	assertOptions(circuitBreakerName, shortSleepInMS)
		{
		if circuitBreakerName isnt false and shortSleepInMS isnt false
			throw "circuitBreakerName and shortSleepInMS options do not work together"
		}
	}