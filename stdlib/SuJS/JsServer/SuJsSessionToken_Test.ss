// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Setup()
		{
		SuJsSessionToken.Register('token', 'key')
		}

	Test_validate()
		{
		env = Object(queryvalues: Object(token: 'token'), cookie: 'token=key')
		Assert(SuJsSessionToken.Validate(env))

		env = Object(queryvalues: Object(), cookie: 'token=key')
		Assert(not SuJsSessionToken.Validate(env))

		env = Object(queryvalues: Object(token: 'token'))
		Assert(SuJsSessionToken.Validate(env) is: false)

		env = Object(queryvalues: Object(token: 'token_wrong'), cookie: 'token=key')
		Assert(SuJsSessionToken.Validate(env) is: false)

		env = Object(queryvalues: Object(token: 'token'), cookie: 'token=wrong_key')
		Assert(SuJsSessionToken.Validate(env) is: false)
		}

	Teardown()
		{
		SuJsSessionToken.Unregister('token')
		}
	}
