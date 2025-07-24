// Copyright (C) 2017 Axon Development Corporation All rights reserved worldwide.
Controller
	{
	output: false
	New()
		{
		.insertServerOutput()
		.output = .FindControl('Output')
		.FindControl('pwd').Set(ServerEval('GetCurrentDirectory'))
		.FindControl('port').Set(ServerPort())
		.tablecountwarning = .FindControl('tablecountwarning')
		.keepAlive()
		if 0 isnt cmdHwnd = FindWindow(NULL, 'Administrator:  Start Suneido Server')
			ShowWindow(cmdHwnd, SW.HIDE)
		Suneido.DevServerLog = .ServerLog
		ServerSuneido.Set('DevHwnd', .WindowHwnd())
		.startTimer()
		}

	insertServerOutput()
		{
		if false isnt output = .FindControl(#serverOutput)
			try
				output.Append(#(ScintillaAddons Addon_inspect:, Addon_show_references:,
					readonly:, height: 30, name: #Output))
			catch (err)
				output.AppendAll(CodeRecoveryControl.ErrorDisplay(err))
		}

	keepAlive()
		{
		SetTimer(NULL, 0, 5.MinutesInMs(), /*= interval < server timeout */
			{|@unused| ServerEval(#Date); 0 })
		}

	Commands()
		{
		return #(
			(Inspect,			"F4",		"Inspect"),
			(Go_To_Definition, 	"F12"),
			(Find_References, 	"F11"))
		}
	Controls()
		{
		.conts = Object()
		Contributions('CSDevServerWindowOptions').Each({ .conts.Add(new it) })
		return Object('Vert'
			Object('Horz' Object('Static' 'Suneido Server is Running in ' )
				#(Static '' name: 'pwd')
				#(Static ' on port ') #(Static '' name: 'port'))
			Object('Horz', .menuBar())
			#(Vert, name: 'serverOutput')
			#(Horz #(Static '' name: 'tablecountwarning')))
		}

	menuBar()
		{
		launchOptions = Object()
		.conts.Each({ it.LaunchOptions(launchOptions) })
		menuBar = Object('Horz',
			Object('MenuButton', 'Launch', launchOptions))

		tools = Object()
		.conts.Each({ it.MenuButtons(menuBar, tools) })
		menuBar.Add(
			Object('MenuButton' 'Tools', tools),
			#Fill, #(Button Clear),
			#Fill, #(MenuButton Close #(Close, 'Close and Compact',
				'Regenerate Dev Files'))
			)
		return menuBar
		}

	Recv(@args)
		{
		if args[0].Prefix?('On_')
			.conts.Each()
				{
				if it.Method?(args[0])
					it[args[0]](@+1args)
				}
		return 0
		}

	serverPostRunfile: 'serverPostRun.bat'
	On_Close(cmd, source /*unused*/)
		{
		switch cmd
			{
		case 'Close and Compact':
			.closeAndCompact()
		case 'Regenerate Dev Files':
			.closeAndRegen()
		default:
			.close()
			}
		}

	close()
		{
		if not .allowClose?()
			return
		super.On_Close()
		}

	serverPostRunExtra: false
	setupServerPostRun()
		{
		serverPostRun = Object('set errormsg=""')
		// sleep to ensure server is shut down
		serverPostRun.Add('sleep 5')
		if .serverPostRunExtra isnt false
			serverPostRun.Add(@.serverPostRunExtra)
		serverPostRun.Add(':final', 'exit',
			':error', 'echo %errormsg%')
		PutFile(.serverPostRunfile, serverPostRun.Join('\r\n'))
		SystemNoWait(.serverPostRunfile)
		}

	allowClose?()
		{
		if .hasConnections?()
			return YesNo('There are open connections. Are you sure you want to' $
				'Shutdown the server')
		return true
		}

	hasConnections?()
		{
		return Sys.Connections().Has?('IDE') or
			Sys.Connections().HasIf?({it.Has?('@') })
		}

	closeAndCompact()
		{
		if .serverPostRunExtra is false
			.serverPostRunExtra = Object()
		.serverPostRunExtra.Add('gsport.exe -compact',
			'if %ERRORLEVEL% neq 0 (
				set errormsg="COMPACT FAILED"
				goto error
				)')
		.close()
		}

	closeAndRegen()
		{
		GDevCreateFiles()
		.close()
		}

	timer: false
	startTimer()
		{
		.libs = Libraries()
		.timer = SetTimer(NULL, 0, 1.SecondsInMs(), .timerFunc)
		}

	timerFunc(@unused) // args from timer
		{
		if .libs isnt curLibs = Libraries()
			{
			diff = curLibs.Difference(.libs)
			.libs = curLibs
			Unload()
			ResetCaches()
			LibraryTags.Reset()
			ServerEval('CSDevServerWindow.DoExtraSetups', diff)
			}
		}
	killTimer()
		{
		if .timer isnt false
			{
			KillTimer(NULL, .timer)
			.timer = false
			ClearCallback(.timerFunc)
			}
		}

	Destroy()
		{
		.killTimer()
		if false is ServerSuneido.Get('DevExeHandled', false)
			.setupServerPostRun()
		Shutdown(alsoServer:)
		super.Destroy()
		}

	On_Clear()
		{
		if .output is false
			return
		.output.Set('')
		.output.Update()
		}

	ServerLog()
		{
		if .output is false
			return
		prints = ServerSuneido.Get('csDevPrint', Object())
		ServerSuneido.Set('csDevPrint', Object())
		for s in prints
			.output.AppendText(s)
		.output.Update()
		}

	ConfirmDestroy()
		{
		return .allowClose?()
		}
	}
