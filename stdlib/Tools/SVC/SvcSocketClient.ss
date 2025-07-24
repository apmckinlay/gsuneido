// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
/*
This class manages a pool of SocketClients for SvcClient. This class is designed to:
-> Open / Reserve SocketClients
	-> Run blocks with the reserved SocketClient
		-> Free / Re-add the SocketClient to the pool.

To test connections manually run: SvcSocketClient().TestConnect( ... )
- A "server" and valid Svc credentials are required in order to carry out the connect.
- NOTE: This will change the settings the singleton is using
- NOTE: If the connection fails, all open clients will be closed

For complex processing, call: SvcSocketClient().Run(){ |sc| ... }
- This will handle getting a valid / free SocketClient and run the block accordingly
- If settings change while a SocketClient is used, it will be closed upon completion
	as it's credentials may no longer be valid.

For simple calls/returns, call: SvcSocketClient.Send(['CMD', arg1: val1 .. ])
- This will do everything from reserving a SocketClient, to reading and unpacking
	the final result

NOTES:
1. This class is designed to run using the table "svc_settings".
2. This class is designed to NOT attempt multiple connects if a failure has occurred
	and the settings haven't changed.
3. This class is designed to run with SvcSettingsIcons. SvcSettingsIcons is designed
	to reflect the current status of SvcSocketClients.
	- While this class is designed to run with SvcSettingsIcons, it is not required.
4. The getters for Settings and Connections are designed for debugging.
*/
Singleton
	{
	Error:			''
	state:			'closed'
	server:			''
	userId:			''
	passhash:		''
	guiEnv?:		false
	New()
		{
		.guiEnv? = Sys.GUI?()
		.connections = Object()
		.connectionsInUse = Object() // Only used for informational purposes
		.getSvcSettings()
		PubSub.Subscribe('SvcSettings_ConnectionModified', .updateSettings)
		// never unsubscribed
		}

	getSvcSettings()
		{
		if false is settings = .svcSettings()
			settings = []
		if settings.svc_local? is true or settings.svc_server is ''
			.server = .userId = .passhash = ''
		else
			{
			.server = settings.svc_server
			.userId = settings.svc_userId
			.passhash = settings.svc_passhash
			}
		}

	svcSettings()
		{ return SvcSettings.Get() }

	updateSettings()
		{
		.state = 'updating' // Force publish event
		.getSvcSettings()
		.Close()
		}

	Close(.Error = '')
		{
		if .Error isnt ''
			.logError(.Error)
		.setState('closed')
		.connections.Each({ .close(it.sc) })
		.connections = Object()
		}

	logError(error)
		{
		if .fatalError?(error)
			SuneidoLog('ERROR: (CAUGHT) SvcSocketClient: ' $ error)
		}

	fatalError?(error)
		{
		excludeErrors = Object(.InvalidCredentials, .InvalidServer, .InvalidKey)
		return not .connectionError?(error) and not excludeErrors.Any?({ error.Has?(it) })
		}

	connectionError?(error)
		{
		return error.Prefix?('socket') or error.Suffix?('SocketClient') or
			error.Has?('lost connection') or error.Has?('timeout')
		}

	close(sc)
		{
		try .send(sc, [#QUIT])
		try sc.Close()
		}

	setState(state)
		{
		publish? = .state isnt state and .guiEnv?
		.state = state
		if publish?
			PubSub.Publish('SvcSocketClient_StateChanged')
		}

	Send(args, result = false)
		{
		.Run({|sc| result = .send(sc, args) })
		return result
		}

	send(sc, args)
		{
		.Write(sc, args)
		return .Read(sc, args)
		}

	Write(sc, args)
		{
		s = Pack(args)
		sc.Writeline(s.Size())
		sc.Write(s)
		}

	Read(sc, params = '')
		{
		.verifyResult(result = Unpack(sc.Read(Number(sc.Readline()))), params)
		return result
		}

	InfoPrefix: 'INFO '
	ErrorPrefix: 'ERR '
	verifyResult(result, params)
		{
		if not String?(result)
			return
		if result.Prefix?(.ErrorPrefix)
			throw result.AfterFirst(.ErrorPrefix)
		if result.Prefix?(.InfoPrefix) and .Verbose()
			SuneidoLog('INFO SvcSocketClient: ' $ result.AfterFirst(.InfoPrefix), :params)
		}

	// Set to true to log server 'INFO ' (InfoPrefix) messages
	verbose?: false
	Verbose(verbose? = '')
		{
		if verbose? isnt ''
			.verbose? = verbose?
		return .verbose?
		}

	Run(block)
		{
		if .Error is '' and false isnt scOb = .reserveSocketClient()
			{
			.connectionsInUse.Add(scOb, at: key = Timestamp())
			try
				block(scOb.sc)
			catch (err)
				.Close(err)
			.freeSocketClient(scOb)
			.connectionsInUse.Delete(key)
			}
		}

	reserveSocketClient()
		{
		scOb = .connections.PopLast()
		if not scOb.Member?('sc') // got the base object, not a SocketClient Object
			scOb = .connect()
		return scOb
		}

	connect()
		{
		scOb = false
		if .openNewSocket?() and .validSocketClient?(sc = .openSocketClient())
			{
			scOb = [:sc, server: .server, passhash: .passhash]
			.setState('open')
			}
		return scOb
		}

	openNewSocket?()
		{ return .Error is '' and not .socketClientAvailable?() }

	socketClientAvailable?()
		{ return not .connections.Empty?() }

	invalidConnection: 'Socket Client could not connect'
	openSocketClient()
		{
		sc = false
		try
			if .server is ''
				throw 'Server is not set'
			else
				sc = .SC()
		catch (err)
			{
			if sc isnt false
				try sc.Close()
			.Close(.invalidConnection $ Opt(':\r\n\r\n', err))
			sc = false
			}
		return sc
		}

	SC()
		{ return SocketClient(.server, 2222 /*= default port */, timeoutConnect: 1) }

	validSocketClient?(sc)
		{
		if sc is false
			return false
		validServer? = validCredentials? = false
		try
			if validServer? = .svcServer?(sc)
				validCredentials? = .validCredentials?(sc)
		catch (err)
			{
			if validServer?
				.close(sc)
			else
				try sc.Close()
			.Close(err)
			}
		return validServer? and validCredentials?
		}

	InvalidServer: 'Server is not a Suneido Version Control Server'
	ValidServer: 'Suneido Version Control Server'
	svcServer?(sc)
		{
		svcServer? = sc.Readline() is .ValidServer
		if not svcServer?
			throw .InvalidServer
		return svcServer?
		}

	InvalidCredentials: 'Invalid Credentials'
	validCredentials?(sc)
		{
		if not valid? = .verifyCredentials(sc)
			throw .InvalidCredentials
		return valid?
		}

	InvalidKey: 'Unable to Retrieve Key'
	verifyCredentials(sc)
		{
		if .userId is '' or .passhash is ''
			throw .InvalidCredentials $ '\r\n\r\nLogin or Password is not set'
		if false is key = .send(sc, [#NONCE])
			throw .InvalidKey
		args = [#LOGIN, svcuser_id: .userId, svcuser_passhash: Sha256(key $ .passhash)]
		return .send(sc, args)
		}

	freeSocketClient(scOb)
		{
		if .state is 'closed' or not .settingsCurrent?(scOb.server, scOb.passhash)
			.close(scOb.sc)
		else
			.connections.Add(scOb)
		}

	settingsCurrent?(server, passhash)
		{ return server is .server and passhash is .passhash }

	Getter_Connected?()
		{ return .state is 'open' }

	Getter_Connections()
		{
		return [
			free: .connections.DeepCopy().Map!({ it.Delete(#passhash) }),
			used: .connectionsInUse.DeepCopy().Values().Map!({ it.Delete(#passhash)})
			]
		}

	Getter_Settings()
		{ return [server: .server, state: .state] }

	// For checking connection, use .Run() if requiring SvcSocketClient process
	TestConnect(server, print? = false)
		{
		.Close()
		serverBefore = .server
		.server = server
		scOb = .connect()
		connected? = scOb isnt false
		if connected?
			.close(scOb.sc)
		.server = serverBefore
		if print? and .Error isnt ''
			Print('TestConnect Error: ' $ .Error)
		return connected?
		}

	RetryState()
		{
		if .Error isnt '' and .connectionError?(.Error)
			.Close('')
		}

	// No need to Reset this, SocketClients close / open themselves as required
	Reset()	{ }
	}