// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.origSession = Database.SessionId()
		}

	Test_LockAndUnlock()
		{
		key1 = 'key1'
		key2 = 'key2'
		sessionId = 'sessionId'

		lm = new LockManagerImpl

		Database.SessionId(sessionId)
		Assert(lm.Lock(key1))
		Assert(lm.Lock(key1) is: sessionId)
		Assert(lm.Lock(key2))
		Assert(lm.Lock(key1) is: sessionId)
		.force_expired(lm, key1)
		Assert(lm.Lock(key1))
		lm.Unlock(key2)
		Assert(lm.Lock(key2))
		// should be able to unlock twice with no errors
		lm.Unlock(key1)
		lm.Unlock(key1)
		Assert(lm.Lock(key1))
		}

	force_expired(lm, key)
		{
		locks = lm[lm.Members()[0]]
		locks[key].expiry = Date().Plus(seconds: -1)
		}

	Test_Renew()
		{
		key = 'renew_key'
		sessionId1 = 'sessionId1'
		sessionId2 = 'sessionId2'

		lm = new LockManagerImpl

		// user should be able to renew his own key
		Database.SessionId(sessionId1)
		Assert(lm.Lock(key))
		Assert(lm.Renew(key))

		// sessionId1 goes to sleep, causing lock to expire.
		// other user locks the access.
		.force_expired(lm, key)
		Database.SessionId(sessionId2)
		Assert(lm.Lock(key))
		Assert(lm.Renew(key))

		// sessionId1 comes back. Should not be able to renew the other sessionId2's key
		Database.SessionId(sessionId1)
		Assert(lm.Renew(key) is: sessionId2)

		// same as above case, but sessionId2 left the Access before sessionId1 renews
		lm.Unlock(key)
		Assert(lm.Renew(key) is: 'lock_expired_but_renewed')
		}

	Test_UserFromSessionId()
		{
		Assert(LockManager.UserFromSessionId('test') is: 'someone')
		Assert(LockManager.UserFromSessionId('fred@192.168.1.55') is: 'fred')
		Assert(LockManager.UserFromSessionId('wilma@wts15') is: 'wilma')
		}

	Teardown()
		{
		Database.SessionId(.origSession)
		}
	}