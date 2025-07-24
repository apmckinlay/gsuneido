// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(conns = false)
		{
		if conns is false
			conns = .userConnections()
		.skipSessionsToKill(conns)
		users = conns.Map({
			it.BeforeLast('@').RemovePrefix(Login.PreLogin)
			}).UniqueValues().Size()
		ips = conns.Map({ it.AfterLast('@').BeforeFirst('<') }).UniqueValues().Size()
		return Max(users, ips)
		}

	userConnections()
		{
		return UserConnections(ignoreDefault:)
		}

	skipSessionsToKill(conns)
		{
		Sys.SkipSessionsToKill(conns)
		}
	}