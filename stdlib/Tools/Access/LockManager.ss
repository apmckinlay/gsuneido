// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
ServerEvalProxy
	{
	Remote: "LockManagerImpl()"

	DefaultKey(rec, query)
		{
		return ShortestKey(query).Split(',').Map!({|k| rec[k] }).Join('\x01')
		}

	DoWithLock(rec, query, block, exceptionIfLocked = false)
		{
		lock_key = .DefaultKey(rec, query)
		if true isnt sid = .Lock(lock_key)
			{
			.recordLocked(sid, exceptionIfLocked)
			return
			}
		Finally({ block(rec) }, { .Unlock(lock_key) })
		}

	recordLocked(sid, exceptionIfLocked = false)
		{
		user = .UserFromSessionId(sid)
		msg = 'This record is already being edited by ' $ user
		if exceptionIfLocked
			throw msg
		else
			Alert(msg, title: 'Record Locked', flags: MB.ICONINFORMATION)
		}
	}