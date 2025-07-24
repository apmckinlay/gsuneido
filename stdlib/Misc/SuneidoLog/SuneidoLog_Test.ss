// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_switchPrefix()
		{
		switchPrefix = SuneidoLog.SuneidoLog_switchPrefix
		s = "ERROR: error"
		Assert(switchPrefix(s, 10) is: s)

		mock = Mock(SuneidoLog)
		mock.When.SuneidoLog_tryGetMessageCount([anyArgs:]).CallThrough()
		mock.When.SuneidoLog_checkPrefix([anyArgs:]).Return(false)
		mock.When.SuneidoLog_getMessageCount([anyString:]).Return(0)
		Assert(mock.Eval(switchPrefix, "WARNING: warning", 10) is: "WARNING: warning")

		mock = Mock(SuneidoLog)
		mock.When.SuneidoLog_tryGetMessageCount([anyArgs:]).CallThrough()
		mock.When.SuneidoLog_getMessageCount([anyString:]).Return(99)
		mock.When.SuneidoLog_checkPrefix([anyArgs:]).Return(false)
		Assert(mock.Eval(switchPrefix, "WARNING: warning", 10)
			is: "ERROR: (Switched from Warning) warning")

		mock = Mock(SuneidoLog)
		mock.When.SuneidoLog_tryGetMessageCount([anyArgs:]).CallThrough()
		mock.When.SuneidoLog_getMessageCount([anyString:]).Return(0)
		mock.When.SuneidoLog_checkPrefix([anyArgs:]).Return(false)
		Assert(mock.Eval(switchPrefix, "ERRATIC: msg", 3) is: "ERRATIC: msg")

		mock = Mock(SuneidoLog)
		mock.When.SuneidoLog_tryGetMessageCount([anyArgs:]).CallThrough()
		mock.When.SuneidoLog_getMessageCount([anyString:]).Return(4)
		mock.When.SuneidoLog_checkPrefix([anyArgs:]).Return(false)
		Assert(mock.Eval(switchPrefix, "ERRATIC: msg", 3)
			is: "ERROR: (Switched from Erratic) msg")
		}

	Test_switchedFromErratic()
		{
		mock = Mock(SuneidoLog)
		mock.When.SuneidoLog_buildRecord([anyArgs:]).CallThrough()
		mock.When.SuneidoLog_switchPrefix([anyArgs:]).Return('ERROR: msg')
		mock.When.SuneidoLog_ensureParams([anyArgs:]).Return(Object())
		rec = mock.SuneidoLog_buildRecord("ERRATIC: msg", '', '', '', 3)
		Assert(rec.sulog_params is: #(switchedToErrorAt: 3))

		rec = mock.SuneidoLog_buildRecord("ERRATIC: msg", 'calls', 'params', '', 3)
		Assert(rec.sulog_params is: #(switchedToErrorAt: 3, params: 'params'))

		rec = mock.SuneidoLog_buildRecord("ERRATIC: msg", 'calls', 'params', 'caught', 3)
		Assert(rec.sulog_params is: #(switchedToErrorAt: 3, params: 'params',
			caughtMsg: 'caught'))
		}

	Test_repeated_message?()
		{
		repeated_message? = SuneidoLog.SuneidoLog_repeatedMessage?
		Assert(repeated_message?('hello world') is: false)

		mock = Mock(SuneidoLog)
		mock.SuneidoLog_max_same_message_per_day = 10
		mock.When.SuneidoLog_tryGetMessageCount([anyArgs:]).CallThrough()
		mock.When.SuneidoLog_checkPrefix([anyArgs:]).Return(false)
		Suneido.Delete(#SuneidoLog_MessageCount)
		mock.When.SuneidoLog_getMessageCount([anyArgs:]).Return(9)
		msg = 'ERROR: error'
		9.Times
			{ Assert(mock.Eval(repeated_message?, msg) is: false) }
		mock.Verify.Never().SuneidoLog_log_repeated_message(msg)

		mock.When.SuneidoLog_getMessageCount([anyArgs:]).Return(10)
		Assert(mock.Eval(repeated_message?, msg) is: false)
		mock.Verify.Times(1).SuneidoLog_log_repeated_message(msg)

		mock.When.SuneidoLog_getMessageCount([anyArgs:]).Return(11)
		Assert(mock.Eval(repeated_message?, msg))

		mock.When.SuneidoLog_getMessageCount([anyArgs:]).Return(0)
		Assert(mock.Eval(repeated_message?, 'ERROR: different') is: false)

		// make sure prefix in included in checking (required for .switchPrefix() to work)
		Assert(mock.Eval(repeated_message?, 'INFO: error') is: false)
		}

	cl: SuneidoLog
		{
		SuneidoLog_maxCount: 20
		New(records)
			{
			.records = records
			}
		SuneidoLog_findSuneidoLogs(timestamp /*unused*/, types, block)
			{
			if types is .SuneidoLog_last24HoursFilterTypes
				{
				try
					{
					for record in .records
						if record.sulog_message =~ types // same as what query would do
							block(record)
					}
				catch (ex, "block:")
					{
					if ex is "block:break"
						return
					}
				return
				}
			if types is 'ERROR|WARNING'
				{
				try
					block(Record(sulog_message: "ERROR: error #1",
						sulog_session_id: "127.0.0.1",
						sulog_locals: #(), sulog_timestamp: #20000102.012934,
						sulog_user: "default"))
				catch(ex, 'block:')
					{
					if ex is 'block:break'
						return
					}
				}
			if types is 'RESTART|AUTO-UPDATE'
				{
				block(Record(sulog_message: 'RESTART: post maxCount',
					sulog_session_id: '127.0.0.1', sulog_locals: #(),
					sulog_timestamp: #20000102.012944,
					sulog_user: 'default'))
				}
			if types is "warning: \(CAUGHT\) Transaction:"
				{
				block(Record(sulog_message: 'warning: \(CAUGHT\) Transaction: ' $
					'from [] to [MAX], shuting down.',
					sulog_session_id: '127.0.0.1', sulog_locals: #(),
					sulog_timestamp: #20000102.010100,
					sulog_user: 'default'))
				block(Record(sulog_message: 'warning: \(CAUGHT\) Transaction: ' $
					'is NOT "" to max ',
					sulog_session_id: '127.0.0.1', sulog_locals: #(),
					sulog_timestamp: #20000102.010100,
					sulog_user: 'default'))
				}
			}
		}

	Test_Last24Hours()
		{
		today = #20000102.013000000
		records = [
			[sulog_message: "ERROR: error #1", sulog_session_id: "127.0.0.1",
				sulog_locals: #(), sulog_timestamp: today.Plus(seconds: -10),
				sulog_user: "default"],
			[sulog_message: "WARNING: warning #2", sulog_session_id: "127.0.0.1",
				sulog_locals: #(), sulog_timestamp: today.Plus(seconds: -20),
				sulog_user: "default"],
			[sulog_message: "something from CASSINFO: x", sulog_session_id: "127.0.0.1",
				sulog_locals: #(), sulog_timestamp: today.Plus(seconds: -30),
				sulog_user: "default"],
			[sulog_message: "INFO: info #3", sulog_session_id: "127.0.0.1",
				sulog_locals: #(), sulog_timestamp: today.Plus(seconds: -30),
				sulog_user: "default"]]
		sLog = new .cl(records)
		Assert(sLog.Last24Hours() like: 'Suneido Log Entries from Last 24 Hours:\r\n' $
			"2000-01-02 01:29:50 127.0.0.1      " $
				"ERROR: error #1\r\n\r\n" $
			"2000-01-02 01:29:40 127.0.0.1      " $
				"WARNING: warning #2\r\n\r\n" $
			"2000-01-02 01:29:30 127.0.0.1      " $
				"INFO: info #3\r\n\r\n" $
			"2000-01-02 01:01:00 127.0.0.1      " $
				"warning: \(CAUGHT\) Transaction: from [] to [MAX], shuting down.\r\n"
			)

		records = Object()
		for(i = 0; i <= 25; i++)
			{
			record = Record(sulog_message: 'INFO: info #' $ i,
				sulog_session_id: '127.0.0.1', sulog_user: "default", sulog_locals: #(),
				sulog_timestamp: today.Plus(seconds: -i))
			records.Add(record)
			}

		sLog = new .cl(records)
		log = sLog.Last24Hours().Lines().Remove('')
		Assert(log isSize: 24)
		Assert(log[0] is: "Suneido Log Entries from Last 24 Hours:")
		Assert(log[21] like: "MORE THAN 21 LOGGED RECORDS ... " $
			"ERROR/WARNING exists in trimmed section")
		Assert(log[22] like: "2000-01-02 01:29:44 127.0.0.1      RESTART: post maxCount")
		Assert(log[23] like: '2000-01-02 01:01:00 127.0.0.1      ' $
			'warning: \(CAUGHT\) Transaction: from [] to [MAX], shuting down.\r\n')
		}
	Test_Last24HoursAsObject()
		{
		today = #20000102.013000000
		records = [
			[sulog_message: "ERROR: error #1", sulog_session_id: "127.0.0.1",
				sulog_locals: #(), sulog_timestamp: today.Plus(seconds: -10),
				sulog_user: "default"],
			[sulog_message: "WARNING: warning #2", sulog_session_id: "127.0.0.1",
				sulog_locals: #(), sulog_timestamp: today.Plus(seconds: -20),
				sulog_user: "default"],
			[sulog_message: "something from CASSINFO: x", sulog_session_id: "127.0.0.1",
				sulog_locals: #(), sulog_timestamp: today.Plus(seconds: -30),
				sulog_user: "default"],
			[sulog_message: "INFO: info #3", sulog_session_id: "127.0.0.1",
				sulog_locals: #(), sulog_timestamp: today.Plus(seconds: -30),
				sulog_user: "default"]]

		sLog = new .cl(records)
		result = sLog.Last24HoursAsObject()

		Assert(result has: records[0], msg: 'as object error 1')
		Assert(result has: records[1], msg: 'as object warning 2')
		Assert(result hasnt: records[2], msg: 'as object cassinfo')
		Assert(result has: records[3], msg: 'as object info3')
		Assert(result.Any?({ it.sulog_message is
			"warning: \(CAUGHT\) Transaction: from [] to [MAX], shuting down." }),
			msg: 'as object transaction conflict')

		records = Object()
		for(i = 0; i <= 25; i++)
			{
			record = Record(sulog_message: 'INFO: info #' $ i,
				sulog_session_id: '127.0.0.1', sulog_user: "default", sulog_locals: #(),
				sulog_timestamp: today.Plus(seconds: -i))
			records.Add(record)
			}

		sLog = new .cl(records)
		log = sLog.Last24HoursAsObject()
		Assert(log isSize: 24)
		Assert(log[20] hasSubset: #(sulog_message: "MORE THAN 21 LOGGED RECORDS ..."),
			msg: 'as object over limit')
		Assert(log[21]
			hasSubset: #(sulog_message: "ERROR/WARNING exists in trimmed section"),
				msg: 'as object trimmed')
		Assert(log[22] hasSubset: #(sulog_timestamp: #20000102.012944,
			sulog_message: "RESTART: post maxCount", sulog_session_id: '127.0.0.1'),
			msg: 'as object restart')
		Assert(log[23] hasSubset: #(sulog_timestamp: #20000102.010100,
			sulog_session_id: '127.0.0.1',
			sulog_message: 'warning: \(CAUGHT\) Transaction: from [] to [MAX], ' $
				'shuting down.'), msg: 'as object tran conflict 2')
		}

	Test_secureLogging()
		{
		mock = Mock(SuneidoLog)
		mock.When.SuneidoLog_buildRecord([anyArgs:]).CallThrough()

		rec = mock.SuneidoLog_buildRecord("ERROR: Has secure data", '',
			#(params: 'secure params'), '', 10)
		Assert(rec.sulog_params is: #(params: 'secure params'))
		Assert(rec.sulog_locals isnt:
			#(msg: 'secure logging enabled: locals not logged.'))

		_secureLogging = true
		rec = mock.SuneidoLog_buildRecord("ERROR: Has secure data", '',
			#(params: 'secure params'), '', 10)
		Assert(rec.sulog_params is:
			#(msg: 'secure logging enabled: params not logged.'))
		Assert(rec.sulog_locals is:
			#(msg: 'secure logging enabled: locals not logged.'))
		}
	Test_checkPrefix()
		{
		fn = SuneidoLog.SuneidoLog_checkPrefix

		Assert(fn(''), msg: 'empty')
		Assert(fn('', warn:), msg: 'empty warn')
		Assert(fn('ERROR: testing') is: false, msg: 'error:')
		Assert(fn('INFO: testing') is: false, msg: 'info:')
		Assert(fn('WARNING: testing'), msg: 'no warning:')
		Assert(fn('WARNING: testing', warn:) is: false, msg: 'yes warning:')
		Assert(fn('ERRATIC: testing'), msg: 'no erratic:')
		Assert(fn('ERRATIC: testing', warn:) is: false, msg: 'yes erratic:')

		Assert(fn('ERROR (CAUGHT): testing') is: false, msg: 'error')
		Assert(fn('INFO testing') is: false, msg: 'info')
		Assert(fn('WARNING testing'), msg: 'no warning')
		Assert(fn('WARNING testing', warn:) is: false, msg: 'yes warning')
		Assert(fn('ERRATIC testing'), msg: 'no erratic')
		Assert(fn('ERRATIC testing', warn:) is: false, msg: 'yes erratic')

		Assert(fn('Not an actual ERROR log'), msg: 'not an error log')
		}

	Test_addLibCommitted()
		{
		fn = SuneidoLog.SuneidoLog_addLibCommitted
		Assert(fn('') is: '')

		calls = 'stdlib:PdfDriver.Finish:652
stdlib:Report.Close:531
stdlib:Report.Report_print:229
stdlib:Report.PrintPDF:203
stdlib:Params.PrintPDF:949
stdlib:Params.Params_pdf:894
eval /* function */
stdlib:CatchFileAccessErrors.CallClass:18
stdlib:Params.Params_pdf:892
stdlib:Working.Working_runBlock:49
stdlib:Dialog.ActivateDialog:55
Invalid[]NotFound:Hello.World:33
libraryNotFound:Hello[]Invalid.World:33
libraryNotFound:Hello.World:33
stdlib:StdlibRecordNotFound.World:33'
		newcalls = fn(calls).Lines()
		preCalls = calls.Lines()
		Assert(newcalls.Size() is: preCalls.Size())
		for m, line in newcalls
			if line.Prefix?('eval') or line.Has?('NotFound')
				Assert(line is: preCalls[m])
			else
				{
				Assert(line[..9] matches: '[0-9]+ ')
				Assert(line[9..] is preCalls[m])
				}
		}
	}