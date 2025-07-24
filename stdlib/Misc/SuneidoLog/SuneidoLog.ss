// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(message, calls = "", params = "", switch_prefix_limit = 10, caughtMsg = '')
		{
		rec = Record()
		try
			rec = .run(message, calls, params, switch_prefix_limit, caughtMsg)
		catch (err)
			try ErrorLog('ERROR: ' $ err $' caused by:\n' $ message)
		return rec
		}

	run(message, calls, params, switch_prefix_limit, caughtMsg)
		{
		origMsg = message
		calls = .tryGetCalls(message, calls)
		rec = .buildRecord(message, calls, params, caughtMsg, switch_prefix_limit)
		if TestRunner.RunningTests?()
			{
			.handleTestRunnerOutput(rec)
			return rec
			}

		if not .repeatedMessage?(rec.sulog_message)
			{
			if .outputRecord(rec) is true // true means output to suneidolog
				.incrementMessageCount(rec.sulog_message)
			}
		if Suneido.User is 'default'
			{
			Print('SuneidoLog [' $ Display(rec.sulog_timestamp) $ ']: ' $ origMsg)
			if rec.sulog_message.Prefix?("ERROR") and rec.sulog_calls isnt ''
				Print(rec.sulog_calls)
			}
		return rec
		}
	handleTestRunnerOutput(rec)
		{
		pos = false
		expectedErrors = #()
		if rec.sulog_message.Prefix?("ERROR")
			{
			expectedErrors = ServerSuneido.Get("TestRunningExpectedErrors", Object())
			pos = expectedErrors.FindIf({ |x| rec.sulog_message.Prefix?(x) })
			}
		if pos isnt false
			ServerSuneido.Set('TestRunningExpectedErrors', expectedErrors.Delete(pos))
		else
			{
			.outputRecord(rec)
			logs = ServerSuneido.Get('TestRunningLogs', Object())
			ServerSuneido.Set('TestRunningLogs', logs.Add(rec.sulog_timestamp))
			}
		}

	max_same_message_per_day: 10
	repeatedMessage?(message)
		{
		// erratic and warning do not need to be included because both of these
		// will switch to errors after hitting switch_prefix_limit
		if message =~ `^WARNING[: ]` or message =~ `^ERRATIC[: ]`
			return false

		count = .tryGetMessageCount(message)
		if count is .max_same_message_per_day
			.log_repeated_message(message)
		return count > .max_same_message_per_day
		}

	switchPrefix(message, switch_prefix_limit)
		{
		if .checkPrefix(message, warn:)
			return message

		if .tryGetMessageCount(message) <= switch_prefix_limit
			return message

		type = message.BeforeFirst(':')
		return message.As(
			message.Replace('^' $ type $ ':',
				'ERROR: (Switched from ' $ type.Capitalize() $ ')'))
		}

	tryGetMessageCount(message)
		{ // if this is not authorized can not get the message count from the server
		// there will be no messsage counting if outputting to error log so return 0
		try
			count = .getMessageCount(message)
		catch (unused, 'not authorized')
			count = 0
		return count
		}

	checkPrefix(msg, warn = false)
		{
		return warn is true
			? msg !~ `^WARNING[: ]` and msg !~ `^ERRATIC[: ]`
			: msg !~ `^ERROR[: ]` and msg !~ `^INFO[: ]`
		}

	tryGetCalls(message, calls)
		{
		try
			calls = .getCalls(message, calls)
		catch (unused, 'not authorized')
			calls = ''
		return calls
		}
	getCalls(message, calls)
		{
		if Type(message) is 'Except'
			calls = message.Callstack()
		else if calls is true
			calls = GetCallStack(skip: 4, limit: 10)
		if Object?(calls)
			RemoveAssertsFromCallStack(calls)
		return calls
		}

	buildRecord(message, calls, params, caughtMsg, switch_prefix_limit)
		{
		if Suneido.GetDefault(#TestRunner, false)
			calls = ""
		if caughtMsg isnt ''
			{
			params = .ensureParamsIsObject(params).Copy()
			params.caughtMsg = caughtMsg
			}
		else if message.Has?("CAUGHT") and Suneido.User is 'default' and
			not TestRunner.RunningTests?()
			Print('SuneidoLog: CAUGHT suneidolog missing caughtMsg')

		oldMsg = message
		if oldMsg isnt message = .switchPrefix(message, switch_prefix_limit)
			{
			params = .ensureParamsIsObject(params).Copy()
			params.switchedToErrorAt = switch_prefix_limit
			}

		return Record(
			sulog_user: Suneido.User,
			sulog_message: message,
			sulog_calls: Object?(calls) ? FormatCallStack(calls, levels: 10) : calls,
			sulog_locals: .format_locals(calls),
			sulog_params: .format_params(params),
			sulog_option: Suneido.GetDefault(#CurrentBookOption, "")
			sulog_session_id: .getSessionId()
			)
		}
	ensureParamsIsObject(params)
		{
		if params is ''
			params = Object()
		return Object?(params) ? params : Object(:params)
		}
	getSessionId()
		{
		// we don't always need all of thread name and session id
		// but it's better to keep this simple
		return Thread.Name() $ ': ' $ Database.SessionId()
		}
	format_locals(calls, _secureLogging = false)
		{
		if secureLogging is true
			return #(msg: 'secure logging enabled: locals not logged.')

		return Object?(calls) and calls.Member?(0)
			? LogFormatEntry(calls[0].locals)
			: #()
		}
	format_params(params, _secureLogging = false)
		{
		if secureLogging is true
			return #(msg: 'secure logging enabled: params not logged.')

		return LogFormatEntry(params)
		}
	getMessageCount(message)
		{
		if Sys.Client?()
			return ServerEval('SuneidoLog.SuneidoLog_getMessageCount', message)
		.ensureMessageCount()
		return Suneido[.messageCount][.messageHash(message)]
		}

	messageDate: 	#SuneidoLog_MessageDate
	messageCount: 	#SuneidoLog_MessageCount
	ensureMessageCount()
		{ // only called on server
		countDate = Suneido.GetDefault(.messageDate, false)
		if not Suneido.Member?(.messageCount) or countDate isnt Date().NoTime()
			{
			Suneido[.messageCount] = Object().Set_default(0)
			Suneido[.messageDate] = Date().NoTime()
			}
		}

	messageHash(message)
		{
		return Adler32(message[:: 30 /*= max character length*/])
		}

	log_repeated_message(message)
		{
		SuneidoLog('INFO: Stopped logging repeated message.',
			params: Object(:message, max_messages: .max_same_message_per_day))
		}

	outputRecord(rec)
		{
		try
			{
			.outputToTable(rec)
			return true
			}
		catch (error)
			.outputToLog(rec, error)
		return false
		}

	outputToTable(rec)
		{
		// do timestamp/output on the server
		// so we don't get timestamps out of order
		if Client?()
			{
			rec.sulog_timestamp = ServerEval('SuneidoLog.SuneidoLog_outputToTable', rec)
			return
			}
		rec.sulog_timestamp = Timestamp()

		if .directToLogError?(rec)
			throw 'direct to log'

		try
			QueryOutput('suneidolog', rec)
		catch (unused, "*nonexistent table: suneidolog")
			{
			.Ensure()
			QueryOutput('suneidolog', rec)
			}
		return rec.sulog_timestamp
		}

	Ensure()
		{
		Database('ensure suneidolog
			(sulog_timestamp, sulog_user, sulog_message, sulog_calls,
				sulog_params, sulog_locals, sulog_option, sulog_session_id)
			key(sulog_timestamp)')
		}

	directToLogError?(rec)
		{
		return rec.sulog_message =~ 'too many (active|overlapping update) transactions'
		}

	outputToLog(rec, error)
		{
		try ErrorLog('ERROR encountered when writing to suneidolog: ' $ error $
			.FormatLog(rec))
		}

	incrementMessageCount(message)
		{
		if Sys.Client?()
			return ServerEval('SuneidoLog.SuneidoLog_incrementMessageCount', message)
		.ensureMessageCount()
		return ++Suneido[.messageCount][.messageHash(message)]
		}

	// factored out so we can test
	maxCount: 500
	last24HoursFilterTypes: "ERROR|WARNING|RESTART|AUTO-UPDATE|\<INFO\>"
	Last24Hours()
		{
		// WARNING: This is not precisely 24 hours.
		// Goes back an extra hour to handle cases where this is not being called exactly
		// 24 hours apart which resulted in some of the log entries being missed.
		msg = "Suneido Log Entries from Last 24 Hours:"
		count = 0
		timestamp = false
		yesterday = Date().Plus(hours: -25)
		.findSuneidoLogs(yesterday, .last24HoursFilterTypes)
			{ |x|
			++count
			if count > .maxCount
				{
				msg $= '\nMORE THAN ' $ count $ ' LOGGED RECORDS ...'
				timestamp = x.sulog_timestamp
				break
				}
			msg $= .FormatLog(x)
			}
		if timestamp isnt false
			{
			msg $= .errorsOrWarningsInTrimmedSection(timestamp)
			.findSuneidoLogs(timestamp, "RESTART|AUTO-UPDATE")
				{ |x|
				msg $= .FormatLog(x)
				}
			}
		.findSuneidoLogs(yesterday, "warning: \(CAUGHT\) Transaction:")
			{ |x|
			if x.sulog_message.Has?('from [] to [MAX]') or
				x.sulog_message.Has?('from [""] to [MAX]')
				msg $= .FormatLog(x)
			}
		return msg
		}

	Last24HoursAsObject()
		{
		// WARNING: This is not precisely 24 hours.
		// Goes back an extra hour to handle cases where this is not being called exactly
		// 24 hours apart which resulted in some of the log entries being missed.
		ob = Object()
		count = 0
		timestamp = false
		yesterday = Date().Plus(hours: -25)
		.findSuneidoLogs(yesterday, .last24HoursFilterTypes)
			{ |x|
			++count
			if count > .maxCount
				{
				ob.Add(Record(sulog_timestamp: Date(),
					sulog_message: 'MORE THAN ' $ count $ ' LOGGED RECORDS ...'))
				timestamp = x.sulog_timestamp
				break
				}
			ob.Add(x)
			}
		if timestamp isnt false
			{
			msg = .errorsOrWarningsInTrimmedSection(timestamp)
			if msg isnt ""
				ob.Add(Record(sulog_timestamp: Date(), sulog_message: msg.Trim()))
			.findSuneidoLogs(timestamp, "RESTART|AUTO-UPDATE")
				{ |x|
				ob.Add(x)
				}
			}
		.findSuneidoLogs(yesterday, "warning: \(CAUGHT\) Transaction:")
			{ |x|
			if x.sulog_message.Has?('from [] to [MAX]') or
				x.sulog_message.Has?('from [""] to [MAX]')
				ob.Add(x)
			}
		return ob
		}

	// Reduce duplicate code
	findSuneidoLogs(timestamp, types, block)
		{
		QueryApply('suneidolog where
			sulog_timestamp > ' $ Display(timestamp) $
			' and sulog_message =~ ' $ Display(types))
			{
			block(it)
			}
		}

	errorsOrWarningsInTrimmedSection(timestamp)
		{
		msg = ''
		.findSuneidoLogs(timestamp, 'ERROR|WARNING')
			{ |unused|
			msg $= ' ERROR/WARNING exists in trimmed section'
			break
			}
		return msg
		}

	FormatLog(x)
		{
		error_msg? = x.sulog_message.Prefix?('ERROR')
		msg = '\n' $
			.formatTimestamp(x.sulog_timestamp) $ ' ' $
			x.sulog_session_id.RightFill(minSize: 15) $ ' ' $
			.stripError(x.sulog_message, error_msg?) $ '    ' $
			x.sulog_option $ '\n'
		if x.sulog_params isnt '' and x.sulog_params isnt #()
			msg $= 'PARAMS: ' $ .stripError(String(x.sulog_params), error_msg?) $ '\n'
		if x.sulog_locals isnt '' and x.sulog_locals isnt #()
			msg $= 'LOCALS: ' $ .stripError(String(x.sulog_locals), error_msg?) $ '\n'
		if x.sulog_calls isnt ''
			msg $= 'CALLS: ' $
				.stripError(.addLibCommitted(x.sulog_calls), error_msg?) $ '\n'
		return msg
		}

	addLibCommitted(calls)
		{
		newCalls = ''
		for call in calls.Lines()
			{
			if call.Has?(':')
				{
				lib = call.BeforeFirst(':')
				if lib =~ '[_a-zA-Z0-9]+' and
					(false isnt recName = call.Extract('[A-Z][_a-zA-Z0-9]+')) and
					TableExists?(lib) and
					(false isnt libRec = Query1Cached(
						lib $ ' where group = -1', name: recName)) and
					Date?(libRec.lib_committed)
					newCalls $= libRec.lib_committed.Format('yyyyMMdd') $ ' '
				}
			newCalls $= call $ '\n'
			}
		return newCalls
		}

	formatTimestamp(sulog_timestamp)
		{
		return Date?(sulog_timestamp)
			? sulog_timestamp.StdShortDateTimeSec()
			: sulog_timestamp
		}

	// TODO handle marking as errors explicitly, not by text of message
	stripError(str, error_msg?)
		{
		if not error_msg?
			str = str.Replace('(?i)error', 'problem')
		return str
		}

	GetLast(prefix)
		{
		dateStr = Display(Date().Minus(hours: 12))
		return QueryLast("suneidolog
			where sulog_timestamp > " $ dateStr $ "
			where sulog_message =~ '^" $ prefix $ "'
			sort sulog_timestamp")
		}

	Once(message, calls = '', params = '', caughtMsg = '')
		{
		if .getMessageCount(message) is 0
			.CallClass(message, :calls, :params, :caughtMsg)
		}

	OnceByCallstack(msg)
		{
		calls = GetCallStack(skip: 5, limit: 3)
		callsCopy = calls.Copy()
		RemoveAssertsFromCallStack(callsCopy)
		trackMsg = FormatCallStack(callsCopy, levels: 3)
		if .incrementMessageCount(trackMsg) is 1
			SuneidoLog(msg, :calls)
		}
	}
