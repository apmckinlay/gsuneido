// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// see: http://rack.rubyforge.org/doc/SPEC.html
// see: http://www.python.org/dev/peps/pep-3333/
// TODO handle chunked encoding of request
// TODO allow streaming input and/or output
SocketServer
	{
	Name: "Rack Server"
	Port: 80
	bodyTooLarge: ('413 Request Entity Too Large', (), '')
	badRequest: ('400 Bad Request', (), 'Error 400 (Bad Request)!!!')
	errorResponse: ('500 Internal Server Error', (), 'internal server error')
	New(.app, with = false)
		{
		if with isnt false
			.app = Compose(@with)(app) // combine app and middleware
		}
	Run()
		{
		// forever loop is to allow "persistent" connections
		// i.e. more than one request with the same connection
		// returning from Run will close the connection
		forever
			{
			stage = ''
			env = RackEnv(this)
			try
				{
				stage = 'reading request'
				.readRequestHeader(env)
				if env.request is ''
					return

				stage = 'reading body'
				result = .checkBody(env)

				headRequestOnly = false
				if result is true
					{
					stage = 'in App'
					if env.method is 'HEAD'
						{
						headRequestOnly = true
						env.method = 'GET'
						}
					result = (.app)(:env)
					}

				stage = 'writing response'
				// .app should return -1 if it has handled the response
				// e.g. WebSocketHandler
				if result is -1
					return

				.writeResponse(result, headRequestOnly)
				if env.version is 1.0/*=http version*/ or
					env.GetDefault('connection', '') is 'close'
					return
				}
			catch (e)
				{
				if e.Has?('lost connection')
					return
				.error(stage, e, env)
				if stage is 'writing response'
					return // close connection
				try
					{
					response = e is HttpRequest.BadRequest ? .badRequest : .errorResponse
					.writeResponse(response)
					}
				catch
					return
				}
			}
		}

	readRequestHeader(env)
		{
		.readRequestLine(env)
		env.remote_user = .RemoteUser()
		.readHeaders(env)
		if env.version >= 1.1 and /*= HTTP version 1.1 */
			env.GetDefault('expect', '').Lower() is '100-continue'
			{
			.Writeline('HTTP/1.1 100 Continue')
			.Writeline('')
			}
		}
	readRequestLine(env)
		{
		request = .Readline()
		if request is false
			request = ''
		env.request = request
		// add method, path, query, version
		env.Merge(HttpRequest.SplitRequestLine(request))
		env.queryvalues = Url.SplitQuery(env.query)
		}
	readHeaders(env)
		{
		header = InetMesg.ReadHeader(this)
		InetMesg.HeaderValues(header, env, translateHeaderNames:)
		}

	checkBody(env)
		{
		if false isnt len = env.Extract('content_length', false)
			{
			if not len.Number?()
				return .badRequest
			toRead = Number(len)
			if toRead > 32.Mb() /* = max read size for content body */
				return .bodyTooLarge
			else if toRead < 0
				return .badRequest
			env.content_length = toRead
			}
		return true
		}

	writeResponse(result, headRequestOnly = false)
		{
		result = .ResultOb(result)
		headers = result[1].Copy()
		headers.Server = 'Suneido RackServer'
		content = result[2]
		if headRequestOnly
			{
			headers.Content_Length = result[2].Size()
			content = ''
			}
		HttpSend(this, "HTTP/1.1 " $ result[0], headers, content)
		}
	ResultOb(result) //TODO use multiple return values after BuiltDate 0227
		{
		if String?(result)
			return ['200 OK', #(), result]

		code = .responseCode(result[0])
		if result.Size() is 2
			return [code, #(), result[1]]

		Assert(result.Size() is 3) /*= normal result size is 3 values */
		return code is result[0]
			? result // avoid creating a new object if not necessary
			: [code, result[1], result[2]]
		}
	responseCode(code)
		{
		switch
			{
		case String?(code):
			if HttpResponseCodes.Member?(code)
				{
				num = HttpResponseCodes[code]
				return num $ ' ' $ HttpResponseCodes[num]
				}
		case Number?(code):
			if HttpResponseCodes.Member?(code)
				return code $ ' ' $ HttpResponseCodes[code]
			code = String(code)
			}
		Assert(code =~ `\A\d\d\d( .+)?\Z`)
		return code
		}

	error(stage, e, env)
		{
		// put env in params since it just shows as <object> in locals
		if .logError?(e)
			SuneidoLog("ERROR: (CAUGHT) RackServer: " $ stage $ ': ' $ e,
				params: env,
				caughtMsg: '400 response returned to client if bad request, ' $
					'else returns 500 internal server error to client')
		}

	logError?(e)
		{
		if e.Prefix?('socket')
			return false

		if e.RemovePrefix('ERROR: ').Prefix?('CopyTo') and
			(e.Has?('wsasend') or e.Has?('broken pipe'))
			return false

		return e isnt HttpRequest.BadRequest
		}
	}
