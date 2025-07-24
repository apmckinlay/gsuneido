// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
TestObserver
	{
	New(.quiet = false)
		{
		.succeeded = .failed = 0
		.total_slowqueries = ""
		}
	BeforeTest(name)
		{
		if not .quiet
			.Output(name)
		.ok = true
		.Name = name
		.slowqueries = ""
		}
	BeforeMethod(method)
		{
		if .quiet isnt true
			.Output('    ' $ method)
		.method = method
		}
	Error(method, error)
		{
		if method is '' // not related to a specific test
			{
			.Output('ERROR: ' $ error)
			++.failed
			return
			}

		if .quiet
			.Output(Opt(.Name, ".") $ .method)
		pre = '        ERROR: '
		.Output(pre $ error.Replace('\n', '\n' $ ' '.Repeat(pre.Size())))
		.ok = false
		}
	Warning(method/*unused*/, warning)
		{
		if not .quiet
			.Output('        WARNING: ' $ warning)
		if warning.Has?('SLOWQUERY')
			++.slowqueries
		}
	AfterMethod(method/*unused*/, time = false)
		{
		if not .quiet and time isnt false
			.Output('        time: ' $ time)
		}
	AfterTest(name/*unused*/, time, dbgrowth, memory)
		{
		if not .quiet
			{
			.Output("    time: " $ time.RoundToPrecision(2) $ " sec")
			if memory isnt 0
				.Output("    memory: " $ ReadableSize(memory))
			if dbgrowth isnt 0
				.Output("    db growth: " $ ReadableSize(dbgrowth))
			if .slowqueries isnt ''
				.Output("    slow queries: " $ .slowqueries)
			}
		if .ok
			++.succeeded
		else
			++.failed
		.total_slowqueries += .slowqueries
		}
	After(time, dbgrowth, memory)
		{
		if (not .quiet and .succeeded + .failed isnt 1)
			{
			if memory isnt 0
				.Output("memory: " $ ReadableSize(memory))
			if dbgrowth isnt 0
				.Output("db growth: " $ ReadableSize(dbgrowth))
			if .slowqueries isnt ''
				.Output("slow queries: " $ .total_slowqueries)
			}
		if .failed > 0
			.Output(.failed $ " ERRORS, all tests took " $
				time.RoundToPrecision(2) $ " sec")
		else
			.Output(Plural(.succeeded, "test") $ " SUCCEEDED " $
				time.RoundToPrecision(2) $ " sec")
		return .failed is 0
		}
	HasError?()
		{
		return .failed > 0
		}
	Output()
		{
		throw 'MUST IMPLEMENT IN DERIVED CLASS'
		}
	ClearFailed() // used for running multiple tests from run libview associated tests
		{
		.failed = 0
		}
	}