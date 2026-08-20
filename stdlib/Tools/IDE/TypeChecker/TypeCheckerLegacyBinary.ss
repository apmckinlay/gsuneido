// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// LEGACY-TYPECHECK-BINARY
class
	{
	binaryPath: "//server/d/Software/TypeChecker/suneidotypes.exe"

	Path()
		{
		return UserSettings.Get(#TypeChecker_BinaryPath, .binaryPath)
		}

	SetPath(path)
		{
		if path is .Path()
			return

		UserSettings.Put(#TypeChecker_BinaryPath, path)
		.Stop()
		(.checkAvailable).ResetCache()
		}

	checkAvailable: MemoizeSingle
		{
		Func()
			{
			path = TypeCheckerLegacyBinary.Path()
			if not FileExists?(path)
				{
				SuneidoLog.Once("suneidotypes binary does not exist at " $ path)
				return false
				}
			return true
			}
		}

	Available?()
		{
		return (.checkAvailable)()
		}

	Check(method, orderedSrc, refs, policy, restartOnError? = true)
		{
		request = Json.Encode(
			Object(:method, arguments: orderedSrc, references: refs, config: policy))
		return .send(request, :restartOnError?)
		}

	send(request, restartOnError? = true)
		{
		result = false
		for (i = 0; i < 3 and result is false; ++i)
			try
				result = .callCheck(.Start(), request)
			catch (e)
				{
				if not restartOnError?
					throw e
				// "actively refused" = server actually gone, respawn;
				// "did not properly respond" = AV dropped the dial but the
				// server is healthy - pause briefly and retry the same one
				if e.Has?("actively refused")
					.Stop()
				Thread.Sleep(300)
				}
		if result is false
			return .sendPiped(request) // no TCP, immune to AV filtering, just slower
		return Json.Decode(result.GetDefault(#content, ""))
		}

	sendPiped(request)
		{
		p = .Path()
		if not String?(p)
			return false
		rp = RunPiped(p)
		rp.Write(request)
		rp.CloseWrite()
		response = rp.Read()
		rp.Close()
		return Json.Decode(response is false ? "" : response)
		}

	// if no long running server then spawn else reuse existing server
	Start()
		{
		if .getProps("").Member?(#Server)
			return .getProps(#Server).port

		if not TypeCheckHelper.BinaryExists?()
			throw "Binary unavailable"

		if not Sys.Windows?()
			Spawn(P.WAIT, #chmod, "+x", .Path())
		rp = RunPiped(.Path() $ " -serve")
		port = .readReadyPort(rp)
		.warmUp(port)
		.setProps(#Server, Object(:rp, :port))
		return port
		}

	// AV drops the first loopback dials to a freshly rebuilt binary while it
	// re-evaluates the new hash; drops fail in ~1s, so burn through that
	// window at spawn instead of failing real requests
	warmUp(port)
		{
		for (i = 0; i < 10; ++i)
			try
				{
				SocketClient("127.0.0.1", port, timeoutConnect: 2).Close()
				return
				}
			catch (unused, "*connectex")
				Thread.Sleep(500)
		}

	// only throws on transport failure not on a non 2XX http code
	callCheck(port, request)
		{
		return Http(#POST, "http://127.0.0.1:" $ port $ "/check", request,
			header: Object("Content-Type": "application/json"), timeout: 120)
		}

	Stop()
		{
		if not .getProps("").Member?(#Server)
			return

		server = .getProps(#Server)
		Suneido.TypeCheckProperties.Delete(#Server)

		if Number?(server.GetDefault(#port, false))
			try
				Http.Post("http://127.0.0.1:" $ server.port $ "/shutdown", "", timeout: 5)

		try
			server.rp.Close()
		}

	readReadyPort(rp)
		{
		line = rp.Readline() // binary prints "READY port=NNNN" once listening
		if not String?(line) or not line.Has?("READY port=")
			throw "type checker did not report READY: " $ Display(line)

		return Number(line.AfterFirst('=').Trim())
		}

	// the spawned server is kept where TypeCheckHelper keeps its properties
	getProps(key, def = #())
		{
		if not Suneido.Member?(#TypeCheckProperties)
			Suneido.TypeCheckProperties = Object()

		if key is ""
			return Suneido.TypeCheckProperties

		return Suneido.TypeCheckProperties.GetDefault(key, def)
		}

	setProps(key, val)
		{
		if not Suneido.Member?(#TypeCheckProperties)
			Suneido.TypeCheckProperties = Object()

		Suneido.TypeCheckProperties[key] = val
		}
	}
