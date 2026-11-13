// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
class
	{
	// Run in login request - using server Suneido
	CreateToken(user)
		{
		for i in ..20 /*=try*/
			{
			token = Base64.Encode(user $ '+' $ i)
			if .lockToken?(token)
				{
				key = Md5(token $ Display(Timestamp()) $ EncryptControlKey() $ Built()).
					ToHex()
				return Object(:token, :key)
				}
			}
		// fallback
		token = Base64.Encode(Display(Timestamp()))
		key = Md5(token $ EncryptControlKey() $ Built()).ToHex()
		return Object(:token, :key)
		}

	lockToken?(token)
		{
		.Synchronized()
			{
			if Suneido.GetDefault(#SuSessionTokens, #()).Member?(token)
				return false

			now = Date()
			if Suneido.
				GetInit(#SuSessionTokenCandidates, Object).
				GetDefault(token, false) > now
				return false

			Suneido.SuSessionTokenCandidates[token] = now.Plus(minutes: 5)
			return true
			}
		}

	// Run in client sessions - using client Suneido
	Register(token, key)
		{
		if Sys.Client?()
			return ServerEval("SuJsSessionToken.Register", token, key)

		.Synchronized()
			{
			Suneido.GetInit(#SuSessionTokens, Object)[token] = key
			Suneido.SuSessionTokenCandidates.Delete(token)
			}
		}

	// Run in client sessions  - using client Suneido
	Unregister(token)
		{
		if Sys.Client?()
			return ServerEval("SuJsSessionToken.Unregister", token)

		.Synchronized()
			{
			Suneido.SuSessionTokens.Delete(token)
			}
		}

	// Run in Http request handler - using server Suneido
	Validate(env)
		{
		if not env.queryvalues.Member?(#token)
			return false

		token = env.queryvalues.token
		if false is
			key = Suneido.GetInit(#SuSessionTokens, Object).GetDefault(token, false)
			return false

		if env.GetDefault(#cookie, "").Has?(token $ '=' $ key)
			return true

		return false
		}
	}
