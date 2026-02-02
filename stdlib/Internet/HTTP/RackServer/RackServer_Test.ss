// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_readRequestLine()
		{
		rs = Mock()
		rs.When.Readline().Return(s = "GET url?query=abc HTTP/1.1")
		rs.Eval(RackServer.RackServer_readRequestLine, env = RackEnv(rs))
		Assert(env.Eq?(request: s, method: 'GET', path: 'url',
				query: 'query=abc',	queryvalues: #(query: 'abc'), version: 1.1))

		rs = Mock()
		rs.When.Readline().Return(s = "GET / HTTP/1.1")
		rs.Eval(RackServer.RackServer_readRequestLine, env = RackEnv(rs))
		Assert(env.Eq?(request: s, method: 'GET', path: '/', version: 1.1,
			query: "", queryvalues: #()))

		rs = Mock()
		rs.When.Readline().Return(s = "GET / HTTP/1.1^M")
		err = Catch()
			{
			rs.Eval(RackServer.RackServer_readRequestLine, env = RackEnv(rs))
			}
		Assert(err is: HttpRequest.BadRequest)

		rs = Mock()
		rs.When.Readline().Return(s = "HEAD / HTTP/1.1^M")
		err = Catch()
			{
			rs.Eval(RackServer.RackServer_readRequestLine, env = RackEnv(rs))
			}
		Assert(err is: HttpRequest.BadRequest)
		}

	Test_readHeaders()
		{
		rs = Mock()
		rs.When.Readline().Return("")
		rs.Eval(RackServer.RackServer_readHeaders, env = RackEnv(rs))
		Assert(env.Eq?())

		rs = Mock()
		rs.When.Readline().Return("Name: value", "")
		rs.Eval(RackServer.RackServer_readHeaders, env = RackEnv(rs))
		Assert(env.Eq?(name: 'value'))

		rs = Mock()
		rs.When.Readline().Return("Name: value", "Name-2:  value2  ", "")
		rs.Eval(RackServer.RackServer_readHeaders, env = RackEnv(rs))
		Assert(env.Eq?(name: 'value', name_2: 'value2'))
		}

	Test_checkBody()
		{
		Assert(RackServer.RackServer_checkBody(env = RackEnv.Build()))
		Assert(env.body is: '')

		rs = Mock()
		rs.When.Read(10).Return("a".Repeat(10))
		Assert(rs.Eval(RackServer.RackServer_checkBody,
			env = RackEnv.Build(content_length: '10', socket: rs)))
		Assert(env.content_length is: 10)
		Assert(env.body is: 'a'.Repeat(10))

		env = RackEnv.Build(content_length: String(40.Mb()))
		Assert(RackServer.RackServer_checkBody(env),
			is: RackServer.RackServer_bodyTooLarge)
		Assert(env.body is: '')

		env = RackEnv.Build(content_length: 'not a number')
		Assert(RackServer.RackServer_checkBody(env),
			is: RackServer.RackServer_badRequest)
		Assert(env.body is: '')

		env = RackEnv.Build(content_length: '10\n10')
		Assert(RackServer.RackServer_checkBody(env),
			is: RackServer.RackServer_badRequest)
		Assert(env.body is: '')
		}

	Test_ResultOb()
		{
		fn = RackServer.ResultOb
		Assert(fn('Test') is: #('200 OK', #(), 'Test'))

		Assert(fn(#(200, "Test")) is: #('200 OK', #(), 'Test'))
		Assert(fn(#("NotFound", "Test")) is: #('404 Not Found', #(), 'Test'))

		Assert(fn(#("OK", "Test")) is: #('200 OK', #(), 'Test'))
		Assert(fn(#("200 OK", "Test")) is: #('200 OK', #(), 'Test'))
		Assert(fn(#("NotModified", #(hdr), "Test"))
			is: #('304 Not Modified', #(hdr), 'Test'))

		Assert(fn(#(200, "Test")) is: #('200 OK', #(), 'Test'))
		Assert(fn(#(304, "Test")) is: #('304 Not Modified', #(), 'Test'))
		Assert(fn(#(304, #(hdr), "Test")) is: #('304 Not Modified', #(hdr), 'Test'))
		}

	Test_responseCode()
		{
		rc = RackServer.RackServer_responseCode
		Assert(rc("123 custom") is: "123 custom")
		Assert(rc(123) is: "123")
		Assert(rc("123 my error") is: "123 my error")
		Assert(rc(200) is: "200 OK")
		Assert(rc(304) is: "304 Not Modified")
		Assert(rc("NotModified") is: "304 Not Modified")
		Assert({ rc(false) } throws:)
		Assert({ rc("typo") } throws:)
		Assert({ rc(12) } throws:)
		Assert({ rc(1234) } throws:)
		Assert({ rc("1234") } throws:)
		Assert({ rc("foo 123") } throws:)
		}

	Test_error()
		{
		fn = RackServer.RackServer_error

		stage = ''
		env = #()
		e = 'socket client error'
		fn(stage, e, env)
		Assert(.GetSuneidoLog() is: #(), msg: 'socket error no log')

		e = 'ERROR: CopyTo: readfrom tcp ' $
			'192.168.172.15:8080->192.168.172.109:60616: write tcp ' $
			'192.168.172.15:8080->192.168.172.109:60616: wsasend: An established ' $
			'connection was aborted by the software in your host machine.'
		fn(stage, e, env)
		Assert(.GetSuneidoLog() is: #(), msg: 'copy to readfrom no log')

		// also test with CopyTo errors without ERROR prefix as it may not be present
		e = 'CopyTo: readfrom tcp ' $
			'192.168.172.15:8080->192.168.172.109:60616: write tcp ' $
			'192.168.172.15:8080->192.168.172.109:60616: wsasend: An established ' $
			'connection was aborted by the software in your host machine.'
		fn(stage, e, env)
		Assert(.GetSuneidoLog() is: #(), msg: 'copy to readfrom no log')

		e = 'ERROR: CopyTo: read 20241210233059: sendfile: broken pipe'
		fn(stage, e, env)
		Assert(.GetSuneidoLog() is: #(), msg: 'copy to sendfile no log')

		e = HttpRequest.BadRequest
		fn(stage, e, env)
		Assert(.GetSuneidoLog() is: #(), msg: 'bad request no log')

		stage = 'in App'
		e = 'a different error'
		fn(stage, e, env)
		recs = .GetSuneidoLog()
		Assert(recs[0].sulog_message is: 'ERROR: (CAUGHT) RackServer: in App: ' $ e)
		}

	Test_logError?()
		{
		fn = RackServer.RackServer_logError?

		Assert(fn(''))
		Assert(fn('some other error'))
		Assert(fn('ERROR: some other error'))
		Assert(fn('ERROR: CopyTo: some other error'))
		Assert(fn('CopyTo: some other error'))
		Assert(fn('CopyTo: wsasend: an established connection was aborted') is: false)
		Assert(fn('ERROR: CopyTo: wsasend: an established connection was aborted')
			is: false)
		Assert(fn('CopyTo: broken pipe') is: false)
		Assert(fn('ERROR: CopyTo: broken pipe') is: false)
		Assert(fn('ERROR: wsasend: an established connection was aborted'))
		Assert(fn('something happened! broken pipe'))
		Assert(fn('socket: something happened!') is: false)
		Assert(fn(HttpRequest.BadRequest) is: false)
		Assert(fn('Error! - ' $ HttpRequest.BadRequest))
		}
	}
