// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_timeoutChecker()
		{
		types = #(
			#(set: "ETA", user: "orders", name: "Orders", screen: "/ETA/Trucking/Orders"),
			#(set: "ETA", user: "jobs", name: "Jobs", screen: "/ETA/Jobs/Jobs"))
		minutes = 25
		startTime = Timestamp()
		exeType = 'gsuneido'

		cl = .clFailed
		result = #()
		try
			cl.TimeoutTester_timeoutChecker(types, minutes, startTime, exeType)
		catch (err)
			result = err.SafeEval()
		Assert(result.status is: 'FAILED')
		Assert(result.expectedConnStatus startsWith: 'FAILED\r\n' $
			'\tExpected Clients not started: #("IM")\r\n' $
			'\t\tmissingClients: #("IM")\r\n' $
			'\t\texpectedClients: #("IM", "Orders", "Dispatching", "Tickets", "Jobs")' $
				'\r\n' $
			'\t\texistingClients: #("Orders", "Dispatching", "Tickets", "Jobs")')
		Assert(result.connsWthStatus is: 'FAILED #("timeoutjobs_gsuneido@127.0.0.1", ' $
			'"timeoutorders_gsuneidoNew@127.0.0.1", ' $
			'"timeoutorders_gsuneidoNew@127.0.0.1")')
		Assert(result.startTime is: startTime)

		cl = .clSucceeded
		result = #()
		try
			cl.TimeoutTester_timeoutChecker(types, minutes, startTime, exeType)
		catch (err)
			result = err.SafeEval()
		Assert(result.status is: 'SUCCEEDED')
		Assert(result.expectedConnStatus is: 'SUCCEEDED')
		Assert(result.connsWthStatus is: 'SUCCEEDED #()')
		Assert(result.startTime is: startTime)
		}
	clFailed: TimeoutTester
		{
		TimeoutTester_testerConnections()
			{
			return Object(":main",
				"timeoutorders_gsuneido@127.0.0.1 - heartbeat",
				"timeoutjobs_gsuneido@127.0.0.1 - heartbeat",
				"timeoutjobs_gsuneido@127.0.0.1",
				"timeoutorders_gsuneidoNew@127.0.0.1",
				"timeoutorders_gsuneidoNew@127.0.0.1").Filter(
				{ not it.Has?(' - heartbeat') and not it.Has?(':main') })
			}
		ReportStatus(status, expectedConnStatus, connsWthStatus, startTime, unused)
			{
			throw Display(Object(:status, :expectedConnStatus, :connsWthStatus,
				:startTime))
			}
		TimeoutTester_dumpDb() { }
		TimeoutTester_addToLog(msg /*unused*/) { }
		TimeoutTester_killConnections(types /*unused*/, exeType /*unused*/) { }
		TimeoutTester_copyErrorLogs() { }
		TimeoutTester_getExpectedClients()
			{ return #('IM', 'Orders', 'Dispatching', 'Tickets', 'Jobs') }
		TimeoutTester_getExistingClients()
			{ return #('Orders', 'Dispatching', 'Tickets', 'Jobs') }
		}

	clSucceeded: TimeoutTester
		{
		TimeoutTester_testerConnections()
			{
			return Object(":main",
				"timeoutorders_gsuneido@127.0.0.1 - heartbeat",
				"timeoutjobs_gsuneido@127.0.0.1 - heartbeat").Filter(
				{ not it.Has?(' - heartbeat') and not it.Has?(':main') })
			}
		TimeoutTester_checkExpectedConnections()
			{
			return 'SUCCEEDED'
			}
		ReportStatus(status, expectedConnStatus, connsWthStatus, startTime, unused)
			{
			throw Display(Object(:status, :expectedConnStatus, :connsWthStatus,
				:startTime))
			}
		TimeoutTester_dumpDb() { }
		TimeoutTester_addToLog(msg /*unused*/) { }
		TimeoutTester_killConnections(types /*unused*/, exeType /*unused*/) { }
		TimeoutTester_copyErrorLogs() { }
		}
	}
