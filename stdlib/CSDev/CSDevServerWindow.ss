// Copyright (C) 2016 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass()
		{
		.initDev()
		if FileExists?('serverPostRun.bat')
			DeleteFile('serverPostRun.bat')
		ServerEval('CSDevServerWindow.SetServerPrint')
		ServerSuneido.Set('SchedulerStopEmail', true)
		ServerSuneido.Set('CSDev?', true)
		.DoExtraSetups()
		if not Thread.List().Any?({ it =~ 'Rack Server' })
			Thread(RunHttpServer)
		// start client at end, so server has enough time to start serving
		.StartClient('CSDevServerWindow.OpenWindow()')
		}

	initDev()
		{
		Suneido.NoCredentials = true
		rec = false
		QueryApply('persistent where classname.Prefix?("WorkSpaceControl") and
			user is "default" and set is "IDE"')
			{
			rec = it
			if rec.pos.master isnt false
				break
			}
		libs = rec.pos.master isnt false ? rec.pos.master : #("stdlib")
		loadliberrors = Object()
		for lib in libs
			try Use(lib)
			catch (e)
				loadliberrors.Add(e)
		LibraryTags.Reset()
		Suneido.User = Suneido.User_Loaded = 'none'
		Suneido.user_roles = #('none')
		if not loadliberrors.Empty?()
			CSDevServerPrint("ERROR: Problems loading libraries " $
				loadliberrors.Join(","))
		}

	DoExtraSetups(libs = 'all')
		{
		try
			for extraSetup in Contributions('CSDevExtraSetups')
				{
				skip? = false
				if libs isnt 'all'
					{
					lib = Name(extraSetup).BeforeFirst('_').Lower()
					skip? = not libs.Has?(lib)
					}
				if not skip?
					extraSetup()
				}
		catch (e)
			CSDevServerPrint('EXTRA SETUP ERROR: ' $ e)
		}

	SetServerPrint()
		{
		Suneido.Print = CSDevServerPrint
		return
		}

	OpenWindow()
		{
		try
			{
			Database.SessionId('CSDevServerWindow')
			Thread(.thread)
			.startupSuneidoLog()
			Window(CSDevServerWindowControl, title: 'Suneido Server IDE',
				w: 480, h: 640)
			ServerSuneido.Set('CSDevServerWindowProc', GetCurrentProcessId())
			.StartClient('CSDevServerWindow.StartIDE()')
			}
		catch(err)
			SuneidoLog('SERVER START ERROR: ' $ err)
		}

	StartClient(cmd)
		{
		if "" is serverIP = ServerIP()
			serverIP = '127.0.0.1'
		args = Object(P.NOWAIT, .getExe(), '-c', serverIP, '-p', String(ServerPort()))
		if ClientTokenRequired?()
			args.Add('t:' $ Base64.Encode(Database.Token()))
		args.Add(cmd)
		Spawn(@args)
		}

	defaultExe: `gsuneido.exe`
	getExe()
		{
		return FileExists?(.defaultExe)
			? .defaultExe
			: 'axon.exe'
		}

	StartUnauthorizedClient()
		{
		Spawn(P.NOWAIT, 'gsuneido.exe', '-c', '-iv', '-p', String(ServerPort()))
		}

	StartIDE()
		{
		Database.SessionId('IDE')
		PersistentWindow.Load()
		}

	threadErrorDelay: 5
	threadLoopDelay: 100
	thread()
		{
		Thread.Name("ServerWindowUpdate-thread")
		path = Paths.ToLocal(Sys.ServerDir() $ '/')
		initialSize = FileExists?(path $ 'golang.log')
			? Dir1(path $ 'golang.log', details:).size
			: 0
		forever
			{
			Thread.Sleep(.threadLoopDelay)
			try
				{
				if Suneido.Member?('DevServerLog')
					{
					// can't use last modified date. It stays set to the file
					// creation date until the server closes the file
					// need to use file size instead
					newSize = FileExists?(path $ 'golang.log')
						? Dir1(path $ 'golang.log', details:).size
						: 0
					.readNewLog(initialSize, newSize, path)
					if not ServerSuneido.Get('csDevPrint', Object()).Empty?()
						Defer(Suneido.DevServerLog)
					}
				if ServerSuneido.Get('SafeShutdown', false)
					Defer(Exit)
				}
			catch(err)
				{
				AddFile(path $ 'devError.log', Display(err) $ '\r\n')
				SuneidoLog('Dev ERROR: ' $ err)
				Thread.Sleep(.threadErrorDelay.SecondsInMs())
				}
			}
		}

	// handle suneidologs that occur durring startup
	startupSuneidoLog()
		{
		if false is start = ServerSuneido.Get('start_time', false)
			return
		QueryApply('suneidolog where sulog_timestamp >= ' $
			Display(start.Plus(seconds: -1)))
			{
			ServerPrint('SuneidoLog [' $ Display(it.sulog_timestamp) $ ']: ' $
				it.sulog_message)
			}
		}

	readNewLog(initialSize, newSize, path)
		{
		if initialSize is false or  newSize > initialSize
			{
			fileLines = Object()
			File(path $ 'golang.log')
				{
				it.Seek(initialSize)
				while false isnt s = it.Readline()
					fileLines.Add(s $ '\r\n')
				}
			initialSize = newSize
			ServerEval('CSDevServerWindow.UpdateDevPrint', fileLines)
			}
		}

	UpdateDevPrint(ob)
		{
		csDevPrint = Suneido.GetInit('csDevPrint', Object())
		for ib in ob
			csDevPrint.Add(ib)
		}

	ShutdownAll()
		{
		System('taskkill /PID ' $ ServerSuneido.Get('CSDevServerWindowProc') $ ' /F')
		Shutdown(alsoServer:)
		}
	}
