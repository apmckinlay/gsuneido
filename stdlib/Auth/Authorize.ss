// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function (user, password)
	{
	authStr = user $ '\x00' $ Sha1(Database.Nonce() $ PassHash(user, password))
	return Database.Auth(authStr)
	}
