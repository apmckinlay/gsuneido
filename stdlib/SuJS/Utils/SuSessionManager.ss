// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	CreateLoginToken(user, remote, host, userAgent, tfaEmail = false, run = false)
		{
		tokens = Suneido.GetInit(#SuLoginTokens, Object())
		ob = JsSessionToken.CreateToken(user)
		tokens[ob.token] = Object(:user, :remote, :host, expire: Date().Plus(minutes: 5),
			:userAgent, token: ob.token, key: ob.key, :tfaEmail, :run)
		return ob
		}

	Open(connectid, token, wsHandler, reconnect = false)
		{
		waitResumeTimeout = .setTimeout(wsHandler)
		if reconnect is 'true'
			return .handlerReconnect(connectid, token, wsHandler, waitResumeTimeout)

		tokens = Suneido.GetDefault(#SuLoginTokens, #())
		info = tokens.Extract(token, false)
		if info is false or info.expire < Date()
			return false

		ServerSuneido.DeleteAt(#WebSocketExpiredConnection, token)
		Thread.NewSuneidoGlobal()
		SuInit.FromServer(info)
		SuSessionLog().Login(info.user, info.remote, info.userAgent)
		SuRenderBackend.Init(wsHandler, token, info.key)
		wsHandler.Send(#BINARY, Pack(SuMessageFormatter.FormatResponse(
			SuMessageFormatter.Type.CONNECTED, arg1: connectid)))
		return true
		}

	setTimeout(wsHandler)
		{
		// set the timeout to make the server notice the disconnection earlier
		try
			{
			wsHandler.GetSocket().SetTimeout(75/*=timeout*/)
			return 85/*=75(socket timeout) + 10(extra)*/
			}
		return 180/*=3 minutes (a longer period so that the socket can timeout)*/
		}

	handlerReconnect(connectid, token, wsHandler, waitResumeTimeout)
		{
		time = Date().StdShortDateTimeSec()
		id = '(reconnect ' $ time $ ')'
		.setThreadName(id)
		Database.SessionId(id)

		.Synchronized()
			{
			expired = ServerSuneido.GetAt(#WebSocketExpiredConnection, token, false)
			}
		if expired is true
			return 'Server has closed the connection'

		socket = wsHandler.GetSocket()
		.Synchronized()
			{
			ServerSuneido.Add(#WebSocketReconnect, [:socket, :connectid], token)
			}

		if .waitConnectionResume(token, waitResumeTimeout) is false
			{
			.Synchronized()
				{
				ServerSuneido.DeleteAt(#WebSocketReconnect, token)
				}
			return 'Reconnect expired'
			}
		socket.ManualClose()
		return 'close'
		}

	waitConnectionResume(token, waitResumeTimeout)
		{
		count = 0
		do
			{
			Thread.Sleep(1.SecondsInMs())
			.Synchronized()
				{
				reconnect = ServerSuneido.GetAt(#WebSocketReconnect, token, false)
				}
			if reconnect is false // the reconnectSession has been retrieved
				break
			count++
			}
		while count < waitResumeTimeout
		return count < waitResumeTimeout
		}

	Close()
		{
		SuRenderBackend().Close()
		}

	CheckDBConnection()
		{
		Sys.KillSessionIfNeeded()
			{
			SuRenderBackend().Terminate('DB session closed')
			}
		}

	BeforeDisconnect()
		{
		SuSessionLog().Logout()
		if false isnt backend = SuRenderBackend(noThrow:)
			backend.BeforeDisconnect()
		}

	// return true to quit the connection
	OnConnectionError(e, wsHandler)
		{
		.setThreadName(Thread.Name() $ '(waiting for reconnect)')
		socket = wsHandler.GetSocket()
		socket.ManualClose()
		socket.Close()
		// to avoid closing the reconnect socket twice
		SuRenderBackend().SetReconnectSocket(false)

		token = SuRenderBackend().Token
		ServerSuneido.DeleteAt(#WebSocketExpiredConnection, token)
		count = 0
		forever
			{
			.Synchronized()
				{
				reconnect = ServerSuneido.GetAt(#WebSocketReconnect, token, false)
				ServerSuneido.DeleteAt(#WebSocketReconnect, token)
				}
			if false isnt reconnect
				{
				SuRenderBackend().SetReconnectSocket(reconnect.socket)
				wsHandler.SetSocket(reconnect.socket)
				wsHandler.Send(#BINARY, Pack(SuMessageFormatter.FormatResponse(
					SuMessageFormatter.Type.CONNECTED, arg1: reconnect.connectid)))
				.setThreadName(Thread.Name().RemoveSuffix('(waiting for reconnect)'))
				SuSessionLog().Error(e, reconnected?:)
				return false
				}
			if count++ > 300/*=wait 5 mins*/
				{
				ServerSuneido.Add(#WebSocketExpiredConnection, true, token)
				SuSessionLog().Error(e)
				return true
				}
			Thread.Sleep(1.SecondsInMs())
			}
		}

	setThreadName(name)
		{
		if name.Prefix?('Thread-')
			name = name.AfterFirst(' ')
		Thread.Name(name)
		}
	}
