// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Singleton
	{
	timeout: 120 /*= seconds = 2 minutes */

	New()
		{
		.locks = Object()
		}
	Lock(key)
		{
		.Synchronized()
			{
			return .lock(key)
			}
		}
	lock(key)
		{
		if .locks.Member?(key)
			{
			lock = .locks[key]
			if not .expired?(lock)
				return lock.sid
			.unlock(key)
			}
		.add(key)
		return true
		}

	Unlock(key)
		{
		.Synchronized()
			{
			.unlock(key)
			}
		}
	unlock(key)
		{
		.locks.Delete(key)
		return // don't return anything
		}

	Renew(key)
		{
		.Synchronized()
			{
			if .locks.Member?(key) and .locks[key].sid is .sessionId()
				{
				.renew(key)
				return true
				}
			if true is sid = .lock(key)
				return 'lock_expired_but_renewed'
			return sid
			}
		}
	Reset()
		{
		// disable Reset because we don't want ResetCaches / Singleton.ResetAll
		// to wipe out the locks
		}
	UserFromSessionId(sid)
		{
		return sid.Has?('@') ? sid.BeforeFirst('@') : 'someone'
		}
	AllCurrentLockedKeys()
		{
		curLocks = Object()
		.Synchronized()
			{
			for key in .locks.Members()
				if not .expired?(.locks[key])
					curLocks.Add(key)
			}
		return curLocks
		}

	add(key)
		{
		.locks[key] = Object(
			sid: .sessionId(),
			expiry: .expiry())
		}
	sessionId()
		{
		return Database.SessionId()
		}
	renew(key)
		{
		.locks[key].expiry = .expiry()
		}
	expiry()
		{
		Date().Plus(seconds: .timeout)
		}
	expired?(lock)
		{
		return lock.expiry < Date()
		}

	/* READ ME FIRST:
		- IF you are using this with the intention of updating data, you are WRONG
		- THIS function is only for checking if a record is in edit mode for:
			- Adjusting controls (IE: color, readonly, etc)
			- Informing users (IE: This trip has an order in edit mode, etc.)
		- Returns:
			- If not Locked:	false
			- If Locked: 		the user id
	*/
	CheckIfLockedForReadOnlyPurposes(key)
		{
		.Synchronized()
			{
			if not .locks.Member?(key)
				return false
			lock = .locks[key]
			return .expired?(lock) ? false : .UserFromSessionId(lock.sid)
			}
		}
	}
