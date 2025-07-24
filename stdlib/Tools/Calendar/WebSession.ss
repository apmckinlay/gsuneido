// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
class
	{
	sessionMember: 'Sessions'
	// running on server
	Register(user)
		{
		if not Suneido.Member?(.sessionMember)
			Suneido[.sessionMember] = Record()

		uuid = UuidString()

		Suneido[.sessionMember][uuid] = Object(User: user,
			Expire: Date().Plus(hours: 24))

		return uuid
		}

	// running on http server
	GetValidUser(args)
		{
		if not args.Member?('sessionid')
			return false
		return ServerEval('WebSession.Authenticate', args.sessionid)
		}

	// running on server; return false or user
	Authenticate(uuid)
		{
		if not Suneido.Member?(.sessionMember)
			return false

		session = Suneido[.sessionMember][uuid]
		if session isnt '' and session.Member?('Expire') and session.Member?('User')
			{
			if session.Expire < Date()
				.Close(uuid)
			else
				return session.User
			}
		return false
		}

	// running on server
	Close(uuid)
		{
		// if server restarted while page is open, Sessions could be uninitialized
		if Suneido.Member?(.sessionMember)
			Suneido[.sessionMember].Delete(uuid)
		}
	}