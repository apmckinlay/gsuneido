// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
// Sys is designed so it doesn't require an instance (like Singleton does)
// This has the benefit that autocompletion works on "Sys."
// The functions like .Client?() should not do any work,
// they should just return Suneido.sys members
// NOTE: Suneido.sys members are used for //TAGS:
class
	{
	Init(suneidojs = false)
		{
		sys = Object()
		sys.suneidojs = suneidojs

		sys.win32 = .win32?()
		os = OSName()
		sys.windows = os.Has?('windows')
		sys.linux = os.Has?("linux")
		sys.macos = os.Has?("mac")

		sys.client = .client?(suneidojs)
		sys.server = .server?(suneidojs)

		try
			.init2(sys)
		catch (unused, "not authorized")
			; // call later from init3

		sys.Set_readonly()
		Suneido.sys = sys // after readonly so it doesn't need to be concurrent
		}

	init2(sys) // if client-server this must be done after authorization
		{
		sos = ServerEval(#OSName).Lower()
		sys.linuxserver = sos.Has?('linux')
		sys.serverDir = ServerEval(#ExeDir)
		}

	init3()
		{
		if Suneido.sys.Member?(#serverDir) // if authorized
			return
		sys = Suneido.sys.Copy() // can't modify existing because it's read-only
		.init2(sys)
		sys.Set_readonly()
		Suneido.sys = sys
		}

	SetServer()
		{
		if .Server?()
			return true
		sys = Suneido.sys.Copy() // can't modify existing because it's read-only
		sys.server = true
		sys.Set_readonly()
		Suneido.sys = sys
		}

	win32?()
		{
		try
			{
			// Cannot use SetTimer because SetTimer is defined sujswebgui
			Type(CreateWindowEx)
			return true
			}
		catch
			return false
		}

	client?(suneidojs)
		{
		return suneidojs is true
			? true
			: Client?() // builtin
		}

	server?(suneidojs)
		{
		return suneidojs is true
			? false
			: Server?() // builtin
		}

	SuneidoJs?()
		{
		return Suneido.sys.suneidojs
		}

	Browser?()
		{
		return Suneido.Member?('SuRender')
		}

	Win32?()
		{
		return Suneido.sys.win32
		}

	Windows?()
		{
		return Suneido.sys.windows
		}
	Linux?()
		{
		return Suneido.sys.linux
		}
	MacOS?()
		{
		return Suneido.sys.macos
		}

	LinuxServer?()
		{
		.init3()
		return Suneido.sys.linuxserver
		}

	Client?()
		{
		return Suneido.sys.client
		}

	Server?()
		{
		return Suneido.sys.server
		}

	GUI?()
		{
		return Suneido.sys.suneidojs or Suneido.sys.win32
		}

	Win32Standalone?()
		{
		return .Win32?() and .Standalone?()
		}

	Standalone?()
		{
		return not .Client?() and not .Server?()
		}

	ServerDir()
		{
		.init3()
		return Suneido.sys.serverDir
		}

	Systags()
		{
		return Suneido.sys.MembersIf({ |m| Suneido.sys[m] })
		}

	Connections()
		{
		if .Client?()
			return ServerEval('Sys.Connections')

		connections = Database.Connections()
		connections.Append(
			Thread.List().Filter({ it.Suffix?('(jsS)') }).Map({ it.AfterFirst(' ') }))
		return connections
		}

	Kill(session_id)
		{
		if .Client?()
			return ServerEval('Sys.Kill', session_id)

		.kill(session_id)
		return true
		}

	kill(session_id)
		{
		.killDbSession(session_id)
		toKill = .toKill()
		userRemote = session_id.BeforeLast('<')
		prefix = userRemote.Blank?() ? session_id : userRemote
		for conn in .threadList()
			if conn.Suffix?('(jsS)') and conn.Prefix?(prefix.RemoveSuffix('(jsS)'))
				toKill[conn] = true
		}

	killDbSession(session_id)
		{
		Database.Kill(session_id)
		}

	toKill()
		{
		return Suneido.GetInit('SessionsToKill', { Object() })
		}

	threadList()
		{
		return Thread.List().Map({ it.AfterFirst(' ') })
		}

	KillSessionIfNeeded(block)
		{
		sessionId = Database.SessionId()
		if false is ServerSuneido.GetAt('SessionsToKill', sessionId)
			return
		ServerSuneido.DeleteAt('SessionsToKill', sessionId)
		block()
		}

	SkipSessionsToKill(conns)
		{
		ServerSuneido.Get('SessionsToKill', #()).Members().Each({ conns.Remove(it) })
		}

	MainThread?()
		{
		return .SuneidoJs?()
			? Thread.Name().Suffix?('(jsS)')
			: Thread.Main?()
		}
	}
