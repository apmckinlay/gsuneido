// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Setup()
		{
		JsSessionToken.Register('token', 'key')
		}

	Test_validate()
		{
		env = Object(queryvalues: Object(token: 'token'), cookie: 'token=key')
		Assert(JsSessionToken.Validate(env))

		env = Object(queryvalues: Object(), cookie: 'token=key')
		Assert(not JsSessionToken.Validate(env))

		env = Object(queryvalues: Object(token: 'token'))
		Assert(JsSessionToken.Validate(env) is: false)

		env = Object(queryvalues: Object(token: 'token_wrong'), cookie: 'token=key')
		Assert(JsSessionToken.Validate(env) is: false)

		env = Object(queryvalues: Object(token: 'token'), cookie: 'token=wrong_key')
		Assert(JsSessionToken.Validate(env) is: false)
		}

	Teardown()
		{
		JsSessionToken.Unregister('token')
		}
	}
