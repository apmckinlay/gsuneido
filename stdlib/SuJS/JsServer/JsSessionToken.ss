// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
class
	{
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
			if ServerSuneido.Get(#SuSessionTokens, Object()).Member?(token)
				return false

			now = Date()
			if ServerSuneido.Get(#SuSessionTokenCandidates, Object()).
				GetDefault(token, false) > now
				return false

			ServerSuneido.Add(#SuSessionTokenCandidates, now.Plus(minutes: 5), token)
			return true
			}
		}

	Register(token, key)
		{
		.Synchronized()
			{
			ServerSuneido.Add(#SuSessionTokens, key, token)
			ServerSuneido.DeleteAt(#SuSessionTokenCandidates, token)
			}
		}

	Unregister(token)
		{
		.Synchronized()
			{
			ServerSuneido.DeleteAt(#SuSessionTokens, token)
			}
		}

	Validate(env)
		{
		if not env.queryvalues.Member?(#token)
			return false

		token = env.queryvalues.token
		if false is key = ServerSuneido.GetAt(#SuSessionTokens, token, false)
			return false

		if env.GetDefault('cookie', '').Has?(token $ '=' $ key)
			return true

		return false
		}
	}
