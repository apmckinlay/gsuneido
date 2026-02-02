// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
class
	{
	App(req)
		{
		wsHandler = req.wsHandler

		SuSessionManager.CheckDBConnection()

		if req.type is #OPEN
			return SuSessionManager.Open(req.queryvalues.connectid, req.queryvalues.token,
				wsHandler, req.queryvalues.GetDefault(#reconnect, false), req)

		if req.type is #CLOSE
			{
			.cleanup()
			return ''
			}

		try
			{
			.process(req, wsHandler)
			}
		catch (e)
			{
			.handleError(e, wsHandler)
			}
		return ''
		}

	handleError(e, wsHandler)
		{
if e is WebSocketHandler.QUITLOOP
	{
	SuRenderBackend().AddLog('QuitLoop\r\n' $
		FormatCallStack(e.Callstack(), levels: 20))

	if 1 >= level = wsHandler.GetLevel()
		e = 'QuitLoop when level is ' $ level $ ' (Suggestion 33827) - ' $ e
	}

		if WebSocketHandler.InternalError?(e)
			throw e
		try
			{
			SuRenderBackend().CancelAllReserved()
			SuRenderBackend().DumpStatus(e)
			Handler(e, NULL, e.Callstack())
			actions = SuRenderBackend().Actions
			.send(wsHandler, Pack(FlatObject(actions, maxLevel: 15)))
			}
		catch (err)
			{
			if WebSocketHandler.InternalError?(err)
				throw err
			reason = err $ ' (e: ' $ e $ ')'
			SuRenderBackend().Terminate(e: 'JsWebSocketServer FATAL - ' $ reason, :reason)
			}
		}

	process(req, wsHandler)
		{
		originalAlert = Suneido.GetDefault(#Alert, false)
		Suneido.Delete(#Alert)
		Finally(
			{
			event = Unpack(req.body)
			if false isnt actions = SuRenderBackend().EventHandler(@event)
				.send(wsHandler, Pack(FlatObject(actions, maxLevel: 15)))
			},
			{
			if originalAlert isnt false
				Suneido.Alert = originalAlert
			})
		}

	MessageLoop()
		{
SuRenderBackend().AddLog('MessageLoop\r\n' $
	FormatCallStack(GetCallStack(), levels: 20))
		actions = SuRenderBackend().Actions
		wsHandler = SuRenderBackend().WSHandler
		// use try here to avoid socket errors from skipping .Loop call
		// message control will handle the failure and resend the content later
		try .send(wsHandler, Pack(FlatObject(actions, maxLevel: 15)))
		wsHandler.Loop()
		}

	send(wsHandler, packed)
		{
		if packed.Size() > 100 /*=compress threshold*/
			{
			compressed = Zlib.Compress(packed)
			if packed.Size() > 32_000_000 /*=StringLimit*/
				{
				SuneidoLog('ERROR: (CAUGHT) packed size over limit',
					params: [packed: packed.Size(), compressed: compressed.Size()],
					caughtMsg: 'Need attention. Furture Suneido will not allow this')
				}
			wsHandler.Send(#BINARY, 0xff.Chr()/*=compress flag*/, compressed)
			}
		else
			wsHandler.Send(#BINARY, packed)
		}

check(ob)
	{
	if Object?(ob)
		for item in ob
			if .check(item) is true
				{
				SuneidoLog('item', params: item)
				return true
				}
	return Class?(ob) or Instance?(ob) or Function?(ob)
	}

	cleanup()
		{
		SuSessionManager.Close()
		}

	OnConnectionError(e, wsHandler)
		{
		try
			return SuSessionManager.OnConnectionError(e, wsHandler)
		catch (e)
			{
			SuneidoLog('ERRATIC: (CAUGHT) - JsWebSocketServer.OnConnectionError - ' $ e,
				caughtMsg: 'Session will be terminated')
			return true // stop
			}
		}

	BeforeDisconnect()
		{
		SuSessionManager.BeforeDisconnect()
		}
	}
