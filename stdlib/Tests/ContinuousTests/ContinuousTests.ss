// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(noData? = false, currency = 'USD', skipTags = #(win32, installedSystem),
		dept? = true, additionalServerTests? = false, bookCheckOnly = false,
		largeData? = false, sujs? = false, testType = '')
		{
		result = ""
		DoWithAlertToSuneidoLog()
			{
			result = .doContinuousTests(noData?, currency, skipTags, dept?,
				additionalServerTests?, bookCheckOnly, largeData?, sujs?, testType)
			}
		DeleteFile('TestStart.txt')
		return result
		}
	doContinuousTests(noData?, currency, skipTags, dept?, additionalServerTests?,
		bookCheckOnly = false, largeData? = false, sujs? = false, testType = '')
		{
		result = localtestresult = servertestresult = ""
		try
			{
			result = ServerEval('ContinuousTests.CreateTablesAndDemoData',
				noData?, currency, dept?, skipTags, largeData?, sujs?)

			result $= .getBuiltDateInfo() $ '\r\n'
			result $= 'Libraries in use: ' $ Libraries().Join(',') $ '\r\n'
			result $= SystemSummary() $ '\r\n'
			libs = Libraries().Difference(skipTags)
			result $= Opt('Skipped Libraries for System Tests: ',
				Libraries().Intersect(skipTags).Join(','), '\r\n')
			if result.Has?('ERROR')
				return result $ '\nERROR: ContinuousTests.CreateTablesAndDemoData ' $
					'failed, continuous test ABORTED\n'
			if bookCheckOnly
				return result
			if largeData?
				Suneido.ValidateQueryAny1? = true
			result $= .TryTypeCheckServerStart()
			localtestresult = TestRunner.Run(TestObserverStringLog(
				Opt(testType, '_') $ 'systems_test_log.txt', quiet:), :libs)
			result $= localtestresult
			if additionalServerTests?
				{
				result $= ServerEval('ContinuousTests.TryTypeCheckServerStart')
				servertestresult = 'Server Ran Tests\r\n' $ TestRunner.RunOnServer(:libs)
				result $= servertestresult
				}
			result $= .RunNightlyChecks(skipTags)

			if not Sys.Client?()
				{
				count = QueryCount('views')
				result $= "Check Views (" $ count $ ")"
				if "" is s = CheckViews()
					{
					result $= " - OKAY\n"
					}
				else
					result $= " - FAILURES:\n" $ s $ "\n"
				}
			}
		catch (x)
			result $= "FAILURES: " $ x $ '\n' $ FormatCallStack(x.Callstack()) $ '\n'
		.processTestResults(testType, localtestresult, servertestresult)
		return result
		}

	TryTypeCheckServerStart()
		{
		try
			{
			// avoid overhead and prevents slow test warning
			if TypeCheckHelper.BinaryExists?()
				TypeCheckHelper.Server()
			return ""
			}
		catch (e)
			return "ERROR: unable to start type checker: " $ e $ "\n"
		}

	getBuiltDateInfo()
		{
		if Sys.Client?()
			return 'Server Build date: ' $ ServerEval('Built') $ '\r\n' $
				'Client Build date: ' $ Built()
		else
			return 'Build date: ' $ Built()
		}

	CreateTablesAndDemoData(noData?, currency, dept?, skipTags = #(), largeData? = false,
		sujs? = false)
		{
		result = ''
		Plugins().ForeachContribution('ContinuousTests', 'demoData', showErrors:)
			{ |x|
			result $= (x.createDemoData)(
				noData?, currency, dept?, skipTags, largeData?, :sujs?)
			}
		// ensure the data is all written to disk, so the query optimization will work properly
		Database.Check()
		return result
		}

	RunNightlyChecks(skipTags = #())
		{
		result = ""
		Plugins().ForeachContribution('ContinuousTests', 'nightlyChecks', showErrors:)
			{ |x|
			result $= (x.nightlyCheck)(:skipTags)
			}
		return result
		}

	processTestResults(testtype, localtestresult, servertestresult)
		{
		Plugins().ForeachContribution('ContinuousTests', 'testResults', showErrors:)
			{ |x|
			(x.procesTestResults)(testtype, localtestresult, servertestresult)
			}
		}
	}
