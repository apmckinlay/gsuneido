// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// REFERENCE: http://martinfowler.com/bliki/CircuitBreaker.html
// STATES: [closed] -> (fail with threshold reached) -> 	[open]
// STATES:		-> timout ->	[half open] -> (success) -> [closed]
// STATES:						[half open] -> (fail) -> 	[open]

class
	{
	memberPrefix: 'CircuitBreaker_'
	CallClass(@serviceAndRunnable) // service, runnable, arguments
		{
		service = serviceAndRunnable[0]
		name = .memberPrefix $ service
		if false is ServerSuneido.Get(name, false)
			{
			config = .getServiceConfig(service).Copy()
			config.service = service
			ServerSuneido.Set(name, config)
			}
		config = ServerSuneido.Get(name, false).Copy()
		cb = .getCbInstance(service)
		Finally(
			{
			retval = cb.Run(@+1 serviceAndRunnable)
			},
			{
			.setConfig(name, cb)
			})
		return retval
		}

	setConfig(name, cb)
		{
		ServerSuneido.Set(name, cb.GetConfig())
		}

	getServiceConfig(service)
		{
		if service.Has?('@')
			service = service.BeforeFirst('@')
		allCircuitBreaks = GetContributions(#CircuitBreakerConfig)
		config = allCircuitBreaks.GetDefault(service, #())
		if config.Empty?()
			{
			for fn in Contributions('CircuitBreakerConfigOverride')
				{
				config = fn(service)
				if Object?(config) and not config.Empty?()
					break
				}
			if not Object?(config)
				config = #()
			}
		return config.Copy()
		}

	GetConfig()
		{
		return Object(
			service: .service,
			threshold: .threshold,
			timeout: .timeout,
			origTimeout: .origTimeout
			timeoutIncrement: .timeoutIncrement,
			failureExpiry: .failureExpiry,
			lastFailureTime: .lastFailureTime,
			failures: .failures,
			state: .state,
			prefix: .prefix)
		}

	Reset(service = false)
		{
		if Sys.Client?()
			return ServerEval('CircuitBreaker.Reset', service)
		Suneido.DeleteIf({
			service is false
				? it.Prefix?(.memberPrefix)
				: it is .memberPrefix $ service })
		return true
		}

	GetStatus(service)
		{
		if String?(cb = .getCbInstance(service))
			return cb
		statusStr = ''
		for m in cb.Members()
			statusStr $= m.AfterFirst(.memberPrefix) $ ': ' $ String(cb[m]) $ '\r\n'
		return statusStr
		}

	getCbInstance(service)
		{
		name = .memberPrefix $ service
		if false is config = ServerSuneido.Get(name)
			return 'uninitialized circuit breaker: ' $ name

		return new this(@config)
		}

	MinutesUntilRunnable(service)
		{
		if Class?(this)
			return .getCbInstance(service).MinutesUntilRunnable(false)

		if .state is 'closed' or .halfOpen?()
			return 0

		return ((.timeout - Date().MinusSeconds(.lastFailureTime)) / 60/*=min*/).Ceiling()
		}

	AttemptsRemaining(service)
		{
		if Class?(this)
			{
			cb = .getCbInstance(service)
			if String?(cb)
				return .getServiceConfig(service).GetDefault(
					'threshold', .defaultThreshold)
			else
				return cb.AttemptsRemaining(false)
			}

		switch(.state)
			{
		case 'open':
			return .halfOpen?() ? 1 : 0
		case 'closed':
			.removeExpiryFailures(Date())
			return .threshold - .failures.Size()
			}
		}

	Run(@args) // function, and arguments, should not be called separately
		{
		fn = args[0]
		if .state is 'closed'
			return .runClosed(fn, args)

		if not .halfOpen?()
			return false

		// try once if half open
		result = false
		try
			result = fn(@+1 args)
		catch (err)
			{
			.failedToRunWhenHalfOpen()
			throw err
			}

		if result isnt false
			.closeCircuit()
		else
			.failedToRunWhenHalfOpen()
		return result
		}

	runClosed(fn, args)
		{
		result = false
		try
			result = fn(@+1 args)
		catch (err)
			{
			if .addAndCheckThreshold()
				.openCircuit(err)
			throw err
			}

		if result isnt false
			.failures = Object()
		else if .addAndCheckThreshold()
			.openCircuit(false)

		return result
		}

	closeCircuit()
		{
		.state = 'closed'
		.lastFailureTime = Date.Begin()
		.failures = Object()
		.timeout = .origTimeout
		SuneidoLog('INFO: CircuitBreaker for ' $ .service $ ' is closed (normal state).')
		}

	addAndCheckThreshold()
		{
		curTime = Date()
		.removeExpiryFailures(curTime)
		.failures.Add(curTime)
		return .thresholdExceeded?()
		}

	removeExpiryFailures(curTime)
		{
		if .failureExpiry isnt false
			.failures.RemoveIf({ curTime.MinusMinutes(it) >= .failureExpiry })
		}

	thresholdExceeded?()
		{
		return .failures.Size() >= .threshold
		}

	prefix: 'ERROR'
	openCircuit(result)
		{
		.state = 'open'
		.lastFailureTime = Date()
		.log(.prefix, .service, .timeout, result, .threshold, .failures)
		}
	log(prefix, service, timeout, result, threshold, failures)
		{
		SuneidoLog(prefix $ ': CircuitBreaker for ' $ service $ ' is open, ' $
			'service calls will be suspended for ' $ timeout $ ' seconds.' $
			' last result: ' $ result, params: [:threshold, :failures], calls:)
		}

	failedToRunWhenHalfOpen()
		{
		maxTimeout = 7200
		.lastFailureTime = Date()
		.timeout = Min(.timeout + .timeoutIncrement, maxTimeout)
		}

	/******** INTERNAL ********/
	state: 'closed' // 'open'
	defaultThreshold: 10
	New(.service, .threshold = false, .timeout = 600 /*seconds*/,
		.timeoutIncrement = 0 /*seconds*/, .failureExpiry = false, .prefix = 'ERROR',
		.lastFailureTime = false, .state = 'closed', .failures = false,
		.origTimeout = false)
		{
		if .origTimeout is false
			.origTimeout = .timeout
		if .lastFailureTime is false
			.lastFailureTime = Date.Begin()
		if .failures is false
			.failures = Object()
		if .threshold  is false
			.threshold = .defaultThreshold
		}

	halfOpen?()
		{
		Assert(.state is 'open', 'invalid state')
		// if the state is open, lastFailureTime should be recent
		Date().MinusSeconds(.lastFailureTime) > .timeout
		}
	}
