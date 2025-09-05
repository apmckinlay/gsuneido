// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// called at startup by exe
class
	{
	GsuneidoRequiredBuiltDate: #20250731

	CallClass()
		{
// vvvvvv WARNING: THE CODE BETWEEN THESE COMMENTS IS NOT COVERED BY CODE RECOVERY vvvvvv
		.setSuneido()
		.localeSettings()
		.checkBuiltDate()
		.allowServer()
		if Sys.Win32?()
			{
			InitWin32()
			}
		cmd = AuthorizationHandler(Cmdline())
// ^^^^^^^^^^^ ANY CHANGES MADE TO THE ABOVE CODE SHOULD BE DONE WITH CAUTION ^^^^^^^^^^^
		if not .startWithCodeRecovery(cmd)
			.start(cmd)
		.logServerInitTimes('end time')
		Suneido.init_end = Date()
		}

	recover: 		'\<rcvr\:\>'	// Commandline Example: "gsuneido.exe rcvr:"
	safeLaunch: 	'\<safe\:\>'	// Commandline Example: "gsuneido.exe safe:"
	startWithCodeRecovery(cmd)
		{
		if not String?(cmd)
			return false
		if false isnt flag = cmd.Match(.safeLaunch)
			cmd = cmd.Replace(.safeLaunch).Trim()
		if codeRecoveryEnabled? = cmd.Has?('CSDevServerWindow') or flag isnt false
			CodeRecoveryControl(exit?:) { .start(cmd) }
		else if codeRecoveryEnabled? = (cmd.Match(.recover) isnt false)
			CodeRecoveryControl(exit?:)
		return codeRecoveryEnabled?
		}

	start(cmd)
		{
		.logServerInitTimes('startup cmd')
		.startup(cmd)
		.logServerInitTimes('server startup')
		.serverStartup() // needs to run AFTER .startup
		.logServerInitTimes('after server startup')
		Suneido.initCmd = cmd // for logging in postinit
		try Query1('postinit').text.Eval() // needs Eval
		}

	Repl()
		{
		.setSuneido()
		Suneido.Print = Global('PrintStdout') // Global to pass CheckLibrary
		.localeSettings()
		cmd = AuthorizationHandler(Cmdline())
		cmd = .removeSurroundingQuotes(cmd)
		.startupFromCmdline(cmd)
		}

	startup(cmd)
		{
		cmd = .removeSurroundingQuotes(cmd)
		if cmd is ""
			{
			if Sys.Win32?() and not Server?()
				PersistentWindow.Load()
			}
		else if Sys.Win32?() and .isPersistentSet?(cmd)
			PersistentWindow.Load(cmd)
		else
			{
			if Sys.Win32Standalone?()
				.useInitLibs()
			.startupFromCmdline(cmd)
			}
		}

	useInitLibs()
		{
		try
			{
			Use('configlib')
			for lib in Global('InitLibs')
				Use(lib)
			Unuse('configlib')
			LibraryTags.Reset()
			}
		catch(unused, "Use: invalid library: configlib|can't find InitLibs")
			{
			}
		}

	startupFromCmdline(cmd)
		{
		if false is s = GetFile(cmd)
			s = cmd
		try
			s.Eval() // needs Eval
		catch (e)
			{
			prefix = e.Has?('socket') ? 'ERRATIC' : 'ERROR'
			try
				SuneidoLog(prefix $ ': in Init evaluating: ' $ s $ ' (' $ e $ ')')
			catch
				try ErrorLog(prefix $ ': in Init evaluating: ' $ s $ ' (' $ e $ ')')
			throw e $ " in: " $ cmd
			// exe will log in error.log
			// and do an alert only if not -unattended
			}
		}

	isPersistentSet?(cmd)
		{
		try
			{
			if not cmd.Identifier?() // where set is cmd line causes slow query
				return false
			return not QueryEmpty?('persistent', set: cmd)
			}
		catch (unused, '*not authorized')
			return false // not authorized
		}

	removeSurroundingQuotes(str)
		{
		return str[0] is '"' and str[-1] is '"'
			? str[1 .. -1]
			: str
		}

	checkBuiltDate()
		{
		if BuiltDate().NoTime() < .GsuneidoRequiredBuiltDate
			throw .builtDateErrMsg(.GsuneidoRequiredBuiltDate.ShortDate())
		// exe will log in error.log
		}

	builtDateErrMsg(date)
		{
		return 'ERROR: Please use a version of Suneido built on or after ' $ date
		}

	allowServer()
		{
		if not Server?()
			return

		for check in Contributions('InitAllowServer')
			if '' isnt msg = (check)()
				{
				SuneidoLog('ERROR: ' $ msg)
				Exit(-1)
				}
		}

	serverStartup()
		{
		.InitialStartup()
		}

	// factored out so standalone systems can call this from the go file
	InitialStartup()
		{
		if not Suneido.GetDefault('asService', false)
			return

		for serverStartup in Contributions('InitServerStartup')
			{
			.logServerInitTimes(String(serverStartup).BeforeFirst(" "))
			serverStartup()
			}
		}

	logServerInitTimes(title)
		{
		if Server?() or Sys.Linux?()
			Suneido.init_times.Add(Object(title, Date()))
		}

	ServerInitTimes()
		{
		return 'Server Init times:\r\n' $
			'start time:' $ Display(ServerSuneido.Get('start_time', 'MISSING')) $ '\r\n' $
			ServerSuneido.Get('init_times', #()).
				Map({ it[0] $ ': ' $ Display(it[1]) }).Join('\r\n') $ '\r\n'
		}

	setSuneido()
		{
		Sys.Init()
		Suneido.start_time = Date()
		Suneido.init_times = Object()
		Suneido.Language = #(name: "english", charset: "DEFAULT", dict: "en_US")

		// for TranslateLanguage Cache
		Suneido.CacheLanguage = ""

		Suneido.User = Suneido.User_Loaded = "default"
		Suneido.user_roles = #('admin')

		if not Sys.Win32?()
			{
			Suneido.Print = Global('PrintStdout') // Global to pass CheckLibrary
			Suneido.Alert = AlertToSuneidoLog
			}

		if Client?()
			{
			if Sys.Win32?()
				{
				// Set a dummy tag to avoid defaulting to the server's tags
				// and having tag 'webgui'
				// Don't use tag 'win32gui' because it is not in use on the server
				LibraryTags.AddMode('win32gui', onlyClient?:)
				}
			else
				LibraryTags.Reset(onlyClient?:)
			}
		}

	localeSettings()
		{
		if Sys.Win32?()
			{
			fmt = GetLocaleInfo(LOCALE_USER_DEFAULT, LOCALE.SSHORTDATE)
			Settings.Set('ShortDateFormat', fmt)
			Settings.Set('SystemShortDateFormat', fmt)
			fmt = GetLocaleInfo(LOCALE_USER_DEFAULT, LOCALE.SLONGDATE)
			Settings.Set('LongDateFormat', fmt)
			fmt = GetLocaleInfo(LOCALE_USER_DEFAULT, LOCALE.STIMEFORMAT)
			Settings.Set('TimeFormat', fmt.Replace(':ss', ''))
			fmt = GetLocaleInfo(LOCALE_USER_DEFAULT, LOCALE.STHOUSAND)
			Settings.Set('ThousandsSeparator', fmt)
			}
		else
			{
			Settings.Set('ShortDateFormat', "yyyy-MM-dd")
			Settings.Set('SystemShortDateFormat', "yyyy-MM-dd")
			Settings.Set('LongDateFormat', "dddd, MMMM dd, yyyy")
			Settings.Set('TimeFormat', "h:mm tt")
			Settings.Set('ThousandsSeparator', ",")
			}
		}
	}
