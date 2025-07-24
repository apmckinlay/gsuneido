// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Lock()
		{
		lmhClass = LockManagerHandler
			{
			LockManagerHandler_schedule_lock_renew()
				{/* do nothing */} // Prevent the renew timer from starting
			}
		lmh = new lmhClass()
		Assert(lmh.LockedKeys() is: #())
		Assert(lmh.OwnLock?('1') is: false)
		Assert(lmh.Lock('1'))
		Assert(lmh.OwnLock?('1'))
		Assert(lmh.OwnLock?('1'))
		Assert(lmh.OwnLock?('1'))
		Assert(lmh.LockedKeys() is: #('1'))
		LockManager.Lock('2')  // different handler
		Assert(lmh.LockedKeys() isSize: 1)
		Assert(lmh.Lock('2') has: "This record is already being edited by ")
		Assert(lmh.Lock('3'))
		Assert(lmh.LockedKeys() isSize: 2)
		lmh.Unlock()
		Assert(LockManager.Lock('1'))
		LockManager.Unlock('2')
		LockManager.Unlock('1')
		}
	}