// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Addon
	{
	Trylock()
		{
		if true isnt sid = .lock()
			{
			user = LockManager.UserFromSessionId(sid)
			.AlertInfo(.GetTitle(), "This record is already being edited by " $ user)
			return false
			}
		return true
		}
	locked: false
	lock()
		{
		lock_key = .GetLockKey()
		if true isnt sid = LockManager.Lock(lock_key)
			return sid
		.locked = lock_key
		.schedule_lock_renew()
		// Print(locked: .lock_key())
		return true
		}
	renew_lock()
		{
		Assert(.locked isnt false, "renew_lock - lock should not be false")
		lock_result = LockManager.Renew(.locked)

		// locked successfully. Just renew the lock
		if lock_result is true
			.schedule_lock_renew() // case 1

		// if the lock was renewed, but LockManager did not find the previous lock
		// (i.e. "this" session went to sleep and another user	went in and out of edit
		// mode for this record). In this case, want to see if that user made changes
		else if lock_result is 'lock_expired_but_renewed'
			{
			if 'deleted' is x = .CheckDeleted(quiet:)
				.logAndExit('the record was deleted by another user')

			if .RecordConflict?(x, quiet?:)
				{
				user = x.bizuser_modified isnt '' and x.bizuser_modified isnt 'default'
					? x.bizuser_modified
					: 'someone'
				.logAndExit('the record was modified by ' $ user)
				}
			.schedule_lock_renew()
			return
			}
		// lock renew failed. Another user took the lock and is currently editing
		// (i.e. "this" session went to sleep and another user took it)
		else
			{
			user = LockManager.UserFromSessionId(lock_result)
			.logAndExit(user $ ' is modifying this record')
			}
		}

	logAndExit(extraInfo)
		{
		SuneidoLog('INFO: Forced exit due to lock taken - ' $ extraInfo)
		.AlertError('Fatal Error', 'Lost Connection')
		.Unlock()
		ExitClient(true) // Exit(true) forces immediate Exit. No Destroy or Save.
		}

	Unlock()
		{
		if .locked is false
			return
		.kill_lock_timer()
		LockManager.Unlock(.locked)
		.locked = false
		// Print(unlocked: .locked)
		}
	lock_timer: false
	schedule_lock_renew()
		{

		.lock_timer = Delay(100.SecondsInMs(), .renew_lock) /*= 100 seconds
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