// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
SocketServer
	{
	Name: "WebSocket Server"
	Port: 8888
	New(.app, with = false)
		{
		if with isnt false
			.app = Compose(@with)(app) // combine app and middleware
		}

	Run()
		{
		return WebSocketHandler(.readHeader(), this, .app)
		}

	readHeader()
		{
		env = Object()
		request = .Readline()
		env.request = request
		// add method, path, query, version
		env.Merge(HttpRequest.SplitRequestLine(request))
		env.queryvalues = Url.SplitQuery(env.query)
		env.remote_user = .RemoteUser()
		header = InetMesg.ReadHeader(this)
		InetMesg.HeaderValues(header, env, translateHeaderNames:)
		return env
		}
	}