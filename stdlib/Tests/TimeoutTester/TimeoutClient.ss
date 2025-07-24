// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(name)
		{
		type = name.BeforeFirst('_')
		.log('attempting to start client for ' $ type $ '-' $ ExeName())
		types = GetContributions('TimeoutTesterTypes')
		if false is timeoutTester = types.FindOne({ it.name is type })
			{
			.log('Time out type ' $ type $ ' not found')
			return
			}
		Trace(TRACE.CLIENTSERVER | TRACE.QUERIES | TRACE.LOGFILE)
		sid = Database.SessionId()
		exe = ExeName().RemoveSuffix('.exe')
		timeoutTesterName = timeoutTester.name $ '_' $ exe
		Suneido.User = 'timeout' $ timeoutTester.user $ '_' $ exe
		Database.SessionId(Suneido.User $ '@' $ sid)

		msg = 'Client Built: ' $ Built()
		SuneidoLog.Once(msg)

		msg = 'Server Built: ' $ ServerEval('Built')
		SuneidoLog.Once(msg)

		.log("Client started for " $ type $ ' as ' $ Suneido.User)
		if false is .startPersistentWindow(timeoutTester, name)
			return
		.OpenScreen(type, timeoutTester)
		.wait(type, timeoutTesterName, sid)
		.log("attempting to add client for " $ type $ ' as ' $ Suneido.User)
		ServerEval('TimeoutTester.AddStartedClient', timeoutTesterName)
		}

	startPersistentWindow(timeoutTester, name)
		{
		try
			{
			PersistentWindow.Load(timeoutTester.set)
			}
		catch(err)
			{
			.persistentWindowFailed(name, timeoutTester, err)
			return false
			}
		return true
		}

	persistentWindowFailed(name, timeoutTester, err)
		{
		.log('Persistent Window Load Failed ' $ name $' - ' $ err)
		PutFile('persistent.log'
			QueryAll('persistent where set is ' $ Display(timeoutTester.set)).
				Map(Display).Join('\r\n'))

		if err.Has?('CreateWindow failed')
			{
			if ServerSuneido.Get('TimeoutTester_' $ name, false) is true
				{
				ServerEval('ErrorLog', 'Timeout tester already retried ' $ name)
				return
				}

			ServerSuneido.Set('TimeoutTester_' $ name, true)
			.log('CreateWindow failed - retry with a new client - ' $ name)
			ServerEval('AddFile', 'error.log'
				String(Timestamp()) $
				' INFO: Timeout tester cannot create window - retry ' $ name $ '\r\n')
			TimeoutTester.StartClient(name)
			}
		}

	OpenScreen(type, timeoutTester)
		{
		.log("Book Opened for " $ type $ ' as ' $ Suneido.User)
		try
			{
			DoWithAlertToSuneidoLog()
				{
				if timeoutTester.Member?('screen')
					Defer()
						{
						url = 'suneido:' $ timeoutTester.screen
						Suneido.OpenBooks[Suneido.CurrentBook].Browser.Goto(url)
						}
				else if timeoutTester.Member?('control')
					Global(timeoutTester.control)()
				}
			}
		catch(err)
			{
			.log('Control Construct failed - ' $
				(timeoutTester.Member?('screen')
					? timeoutTester.screen
					: timeoutTester.control) $ ' - ' $ err)
			}
		}

	wait(type, timeoutTesterName, sid)
		{
		.log("Screen Loaded for " $ type $ ' as ' $ Suneido.User)
		// Defer so screen can finish constructing.
		Defer(
			{
			Trace("")
			Trace("=========Starting " $ timeoutTesterName $
				' at: ' $ Display(Date()) $ "==============\r\n")
			Trace("")
			// Start Heartbeat thread
			// Spits out Date Time every minute to trace log.
			Thread(
				{
				Thread.Name(Suneido.User $ '@' $ sid $ ' - heartbeat')
				sleepMs = 100
				minute = 600
				forever
					{
					Trace("")
					Trace("=========" $ Display(Date()) $ "=========")
					Trace("")
					for .. minute
						Thread.Sleep(sleepMs) // wait 1 minute
					}
				})
			})
		}

	log(s)
		{
		Rlog('timeout', s $ '\r\n')
		}
	}
