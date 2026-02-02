// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Setup()
		{
		SuJsTimeoutTester.RunSuJsHttpServer()
		}

	CallClass(minutes = 25, exeType = 'gsuneido')
		{
		ServerSuneido.Set('TimeoutTester_exeType', exeType)
		startTime = Timestamp()
		types = GetContributions('TimeoutTesterTypes').Shuffle!()
		.startClients(types, exeType)

		// Minimize all windows to avoid restarting the AccessLoopAddonManager timeout.
		// When an active window is closed, Windows can sometimes activate other windows.
		// This can cause AccessLoopAddonManager to restart its timeout timer
		SuneidoLog('INFO: before minimize windows')
		if '' isnt res = RunPipedOutput('powershell ' $
			'(New-Object -ComObject Shell.Application).MinimizeAll()')
			SuneidoLog('INFO: minimize all windows failed - ' $ res)
		SuneidoLog('INFO: after minimize windows')

		// ^--- Before Timeout ---^
		.sleepForMinutes(minutes)
		// v --- AFTER Timeout ---v
		.timeoutChecker(types, minutes, startTime, exeType)
		}

	startClients(types, exeType)
		{
		// ensure the log file is fresh
		PutFile(Rlog.CurrentLog('timeout'), '')
		expectedClients = Object()
		expectedClientPorts = Object()
		ServerSuneido.Set('TimeoutTester_ExistingClients', Object())
		ServerSuneido.Set('TimeoutTester_ExistingClientPorts', Object())
		i = 0
		.forTimeoutTestTypes(types, exeType)
			{ |name, user, web?|
			try DeleteFile('./' $ name $ '/trace.log')
			catch (err)
				.log('cannot delete trace log for ' $ Display(name) $ ': ' $ Display(err))
			expectedClients.Add(name)
			clientPort = .clientPort(i++)
			expectedClientPorts.Add(clientPort)
			.startClient(name, clientPort, user, web?)
			}
		ServerSuneido.Set('TimeoutTester_ExpectedClients', expectedClients)
		ServerSuneido.Set('TimeoutTester_ExpectedClientPorts', expectedClientPorts)
		}

	clientPort(i = false)
		{
		if i is false
			i = 20 + Random(80) /*= random in 100*/
		return 10000 + i + ServerPort() /*= port range */
		}

	startClient(name, clientPort, user = false, web? = false)
		{
		exeType = ServerSuneido.Get('TimeoutTester_exeType', 'gsuneido')
		exeName = exeType $ (name.Has?('New') ? 'New' : '') $ '.exe'

		startCode = AuthorizationHandler.AddTokenToCmdLine(
			web? is false
				? ' TimeoutClient(`' $ name $ '`)'
				: ' SuJsTimeoutTester(name: `' $ name $ '`, user: `' $ user $ '`, ' $
					'title: `' $ name $ '`)')

		EnsureDir(name)
		cmd1 = 'cd ' $ name
		cmd2 =  `start ..\` $ exeName
		if clientPort isnt false
			cmd2 $= ' -w=' $ clientPort
		cmd2 $= ' -u -c -p ' $ String(ServerPort()) $ startCode
		.log('command ' $ cmd2)
		result = System(cmd1 $ " && " $ cmd2)
		Rlog('timeout',
			Display(Date()) $ ' - Start Timeout Client ' $ name $
				', result: ' $ String(result) $ ')\r\n')
		}

	StartClient(name)
		{
		if Sys.Client?()
			{
			ServerEval('TimeoutTester.StartClient', name)
			return
			}

		Thread()
			{
			Thread.Sleep(5.SecondsInMs()) /*= in case fighting same trace.log */
			.log('re-attempting to start client for ' $ name)
			clients = ServerSuneido.Get('TimeoutTester_ExpectedClients', #())
			i = clients.Find(name)
			ports = ServerSuneido.Get('TimeoutTester_ExpectedClientPorts', #())
			ports[i] = .clientPort()
			ServerSuneido.Set('TimeoutTester_ExpectedClientPorts', ports)
			.startClient(name, ports[i])
			}
		}

	AddStartedClient(timeoutTester)
		{
		.log("Adding Client: " $ timeoutTester)
		Suneido.TimeoutTester_ExistingClients.Add(timeoutTester)
		// ensure Suneido variable does NOT get returned to the client
		// causing Pack/Unpack errors
		return
		}

	log(s)
		{
		Rlog('timeout', s $ '\r\n')
		}

	minuteCount: 12 //(5sec * 12 = 1 minute)
	sleepTime: 5
	sleepForMinutes(minutes)
		{
		// cannot sleep for > 10 mins or timeout will kill this client
		// using variation of KeepAlive
		elapsed = 0
		elapsedMinute = 0
		while elapsed < minutes*60 /* = convert to seconds*/
			{
			// only want this once a minute, not every 5 seconds
			if elapsedMinute > .minuteCount
				{
				elapsedMinute = 0
				.logMinute()
				}
			seconds = Timer()
				{
				Thread.Sleep(.sleepTime.SecondsInMs())
				try
					Timestamp()
				catch
					Exit()
				.dumpGoRoutines()
				}
			elapsedMinute++
			elapsed += seconds.Int()
			}
		}

	dumpGoRoutines()
		{
		ports = ServerSuneido.Get('TimeoutTester_ExpectedClientPorts', #())
		clients = ServerSuneido.Get('TimeoutTester_ExpectedClients', #())
		for i in ports.Members()
			{
			if ports[i] isnt false
				{
				try
					{
					s = Http.Get('http://127.0.0.1:' $ ports[i] $
						'/debug/pprof/goroutine?debug=2',
						timeout: 5, timeoutConnect: 5)
					PutFile(clients[i] $ '/debug_goroutine.txt', s)
					}
				}
			}
		}

	logMinute()
		{
		ts = Timestamp()
		Rlog('timeout', Display(ts) $ ':' $
			Display(.testerConnections()) $ '\r\n')
		ServerSuneido.Set('SchedulerHeartBeat', Timestamp())
		}

	testerConnections()
		{
		Sys.Connections().Filter(
			{ not it.Has?(' - heartbeat') and not it.Has?(':main') })
		}

	timeoutChecker(types, minutes, startTime, exeType)
		{
		cutoff = Timestamp()
		conns = .testerConnections()
		.addToLog(Display(Date()) $
			' - Checking for "timeout client" after ' $ minutes $ ' minutes: ' $
			Display(conns) $ '\r\n')
		status = .checkConnections(types, conns, exeType)
		connsWthStatus = status $ ' ' $ Display(conns)

		expectedConnStatus = .checkExpectedConnections()
		if expectedConnStatus isnt "SUCCEEDED"
			status = 'FAILED'

		.dumpDb()

		.copyErrorLogs()

		.ReportStatus(status, expectedConnStatus, connsWthStatus, startTime, cutoff)

		.killConnections(types, exeType)
		}
	dumpDb()
		{
		Database.Dump()
		}
	addToLog(msg)
		{
		Rlog('timeout', msg)
		}

	checkConnections(types, conns, exeType)
		{
		.forTimeoutTestTypes(types, exeType)
			{ |user, web?|
			if .hasConnection?(conns, user, web?)
				return 'FAILED'
			}
		return 'SUCCEEDED'
		}
	hasConnection?(conns, user, web?)
		{
		return web?
			? conns.Has?('timeout' $ user $ '@local(jsS)')
			: conns.Has?("timeout" $ user $ '@127.0.0.1')
		}

	checkExpectedConnections()
		{
		status = 'SUCCEEDED'
		expectedClients = .getExpectedClients()
		existingClients = .getExistingClients()

		if expectedClients.Empty?()
			{
			status = 'FAILED: NO CLIENTS WERE SETUP'
			}
		missingClients = expectedClients.Difference(existingClients)
		if not (expectedClients.Difference(existingClients)).Empty?()
			{
			status = 'FAILED\r\n\tExpected Clients not started: ' $
				Display(missingClients) $
				'\r\n\t\tmissingClients: ' $ Display(missingClients) $
				'\r\n\t\texpectedClients: ' $ Display(expectedClients) $
				'\r\n\t\texistingClients: ' $ Display(existingClients)

			for client in missingClients
				status $= '\r\n\r\nGo Routines - ' $ client $ ':\r\n' $
					GetFile(client $ '/debug_goroutine.txt')
			}
		return status
		}
	getExpectedClients() // for tests
		{
		return ServerSuneido.Get('TimeoutTester_ExpectedClients', Object())
		}
	getExistingClients() // for tests
		{
		return ServerSuneido.Get('TimeoutTester_ExistingClients', Object())
		}

	copyErrorLogs()
		{
		file = Paths.Combine(Getenv("APPDATA"), 'suneido' $ ServerPort() $ '.err')
		if FileExists?(file)
			{
			CopyFile(file, Paths.Combine(ExeDir(), 'clienterror.log'), false)
			try DeleteFile(file)
			catch (err)
				SuneidoLog("ERROR: Delete client err failed: " $ Display(file),
					params: err)
			}
		}

	ReportStatus(status, expectedConnStatus, conns,
		startTime /*unused*/, cutoff /*unused*/)
		{
		SuneidoLog('INFO: Timeout Tester - ' $
			"Connection Timeout Tester " $ status $ ": Connections: " $ Display(conns) $
			"; Clients Started: " $ expectedConnStatus)
		}

	GetLogFiles(startTime, exeType)
		{
		.suneidoLog(startTime)
		logs = Object('error.log', 'suneidolog.log')
		logs.Add(Rlog.CurrentLog('timeout'))
		types = GetContributions('TimeoutTesterTypes')
		.forTimeoutTestTypes(types, exeType)
			{ |name, web?|
			if not web?
				logs.Add('./' $ name $ '/trace.log')
			}
		return logs
		}

	suneidoLog(startTime)
		{
		PutFile('suneidolog.log', 'SuneidoLog:\r\n')
		QueryApply('suneidolog where sulog_timestamp >= ' $ Display(startTime))
			{
			AddFile('suneidolog.log', Display(it) $ '\r\n')
			}
		}

	killConnections(types, exeType)
		{
		.forTimeoutTestTypes(types, exeType)
			{ |user|
			Sys.Kill("timeout" $ user.Lower())
			}
		.killProcesses()
		}

	forTimeoutTestTypes(types, exeType, block)
		{
		exeTypes = Object(exeType)
		for type in types
			for exe in exeTypes
				{
				block(name: type.name $ '_' $ exe,
					user: type.user $ '_' $ exe,
					:type,
					web?: false)
				block(name: type.name $ '_' $ exe $ '_sujsweb',
					user: type.user $ '_' $ exe $ '_sujsweb',
					:type,
					web?:)
				}
		return true
		}

	killProcesses()
		{
		RunPiped(PowerShell() $ ' -Command "Get-WmiObject -Class win32_Process ' $
			'-Property Processid,Commandline | Select-Object ' $
			'-Property Processid,Commandline"')
			{ |rp|
			while false isnt s = rp.Readline()
				{
				if not s.Blank?() and s.Has?('TimeoutClient')
					{
					pid = s.Trim().BeforeFirst(' ').Trim()
					LocalCmds.Taskkill(pid)
					}
				}
			}
		}
	}
