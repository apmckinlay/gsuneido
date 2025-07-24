// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_already_authorized()
		{
		mock = Mock(AuthorizationHandler)
		mock.When.authorized?().Return(true)
		cmd = "some stuff"
		result = mock.Eval(AuthorizationHandler, cmd)
		Assert(result is: cmd)
		}

	Test_login() // no auth on command line
		{
		mock = Mock(AuthorizationHandler)
		mock.When.authorized?().Return(false)
		mock.When.cmdlineAuth([any:]).CallThrough()
		result = mock.Eval(AuthorizationHandler, "some stuff")
		Assert(result is:  'Login(origCmd:"some stuff")')
		}

	Test_bad_auth() // auth on command line, but not correct
		{
		ah = AuthorizationHandler
			{
			AuthorizationHandler_authorized?() { return false }
			AuthorizationHandler_authorize(token /*unused*/) { return false }
			AuthorizationHandler_fatal(msg) { throw msg }
			AuthorizationHandler_getFile(unused) { return "bad token" }
			}
		Assert({ ah("t:Zm9v some stuff") } throws: "failed")
		Assert({ ah("t:???? some stuff") } throws: "failed")
		Assert({ ah("t@nonexistent some stuff") } throws: "failed")
		Assert({ ah("t@file some stuff") } throws: "failed")
		}
	Test_good_auth() // correct auth on command line
		{
		ah = AuthorizationHandler
			{
			AuthorizationHandler_authorized?() { return false }
			AuthorizationHandler_authorize(token) { return token isnt "" }
			AuthorizationHandler_fatal(msg) { throw msg }
			AuthorizationHandler_getFile(file) { return file is "none" ? "" : "Zm9v" }
			}
		Assert(ah("t:Zm9v some stuff") is: "some stuff")
		Assert({ ah("t:not_Base64 some stuff") } throws: "failed")
		Assert({ ah("t@none some stuff") } throws: "failed")
		Assert(ah("t@file some stuff") is: "some stuff")
		}

	Test_unnecessary_auth()
		{
		ah = AuthorizationHandler
			{
			AuthorizationHandler_authorized?() { return true }
			AuthorizationHandler_alert(unused) { }
			}
		Assert(ah("t:token some stuff") is: "some stuff")
		Assert(ah("t@file some stuff") is: "some stuff")

		ah = AuthorizationHandler
			{
			AuthorizationHandler_authorized?() { return true }
			AuthorizationHandler_alert(msg) { throw msg }
			}
		Assert({ ah("t:token some stuff") } throws: "unnecessary")
		Assert({ ah("t@file some stuff") } throws: "unnecessary")
		}

	Test_decodeToken()
		{
		fn = AuthorizationHandler.AuthorizationHandler_decodeToken
		token = "aToken"
		encodedToken = Base64.Encode(token)
		Assert(fn(encodedToken) is: token)
		Assert(fn(encodedToken $ 'AFunction()') is: "")
		}

	Test_AddTokenToCmdLine()
		{
		mock = Mock()
		token = 'aToken'
		mock.When.AuthorizationHandler_token().Return(token)
		fn = AuthorizationHandler.AddTokenToCmdLine
		base64 = Base64.Encode(token)
		Assert(mock.Eval(fn, '') is: ' t:' $ base64)
		Assert(mock.Eval(fn, ' MyFunc()') is: ' t:' $ base64 $ ' MyFunc()')
		Assert(mock.Eval(fn, '"MyFunc()"') is: ' t:' $ base64 $ ' "MyFunc()"')
		}
	}