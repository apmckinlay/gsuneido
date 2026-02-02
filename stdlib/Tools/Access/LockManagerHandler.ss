// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// anyplace that needs to maintain a lock should be using this class now.
// TODO: Change AccessControl to use this class
class
	{
	Lock(lock_key)
		{
		return .trylock(lock_key)
		}
	Unlock(key = false)
		{
		.unlock(key)
		}
	OwnLock?(lock_key)
		{
		return Object?(.lockedKeys) and .lockedKeys.Has?(lock_key)
		}
	LockedKeys()
		{
		return Object?(.lockedKeys) ? .lockedKeys : #()
		}
	trylock(lock_key)
		{
		if true isnt sid = .lock(lock_key)
			{
			user = LockManager.UserFromSessionId(sid)
			return "This record is already being edited by " $ user
			}
		return true
		}
	lockedKeys: false
	lock(lock_key)
		{
		if .OwnLock?(lock_key)
			return true
		if true isnt sid = LockManager.Lock(lock_key)
			return sid
		if .lockedKeys is false
			.lockedKeys = Object()
		if .lockedKeys.Empty?()
			.schedule_lock_renew()
		.lockedKeys.AddUnique(lock_key)
		return true
		}
	renew_lock()
		{
		if .noKeys?()
			return

		for lock_key in .lockedKeys
			LockManager.Renew(lock_key)
		.schedule_lock_renew()
		}
	noKeys?()
		{
		.lockedKeys is false or .lockedKeys.Empty?()
		}
	unlock(key)
		{
		if .noKeys?()
			return

		if key is false
			{
			for lock_key in .lockedKeys
				LockManager.Unlock(lock_key)
			.lockedKeys = #()
			}
		else
			{
			LockManager.Unlock(key)
			.lockedKeys.Remove(key)
			}

		if .lockedKeys.Empty?()
			{
			.kill_lock_timer()
			.lockedKeys = false
			}
		}
	lock_timer: false
	schedule_lock_renew()
		{
		.lock_timer = Delay(100000, .renew_lock)	/*= 100 seconds,
			based on LockManager timeout of 120 seconds */
		}
	kill_lock_timer()
		{
		if .lock_timer is false
			return
		.lock_timer.Kill()
		.lock_timer = false
		}
	}