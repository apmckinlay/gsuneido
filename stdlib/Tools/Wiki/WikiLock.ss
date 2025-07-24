// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (page, remote_user = '')
	{
	if not Suneido.Member?('WikiLock')
		Suneido.WikiLock = Object()
	if Suneido.WikiLock.Member?(page)
		return Suneido.WikiLock[page]
	Suneido.WikiLock[page] = Object(user: remote_user, date: Date())
	return true
	}