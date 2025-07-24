// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		if Thread.Name() isnt "main"
			return
		save = Database.SessionId()
		.AddTeardown({ Database.SessionId(save is "" ? " " : save) })
		threadName = Thread.Name()
		Database.SessionId("SessionId_Test")
		Assert(Database.SessionId() is: "SessionId_Test")
		// setting session id should not affect thread name
		Assert(Thread.Name() is: threadName)
		ob = Object()
		.threadSetSessionId(ob)
		if ob.Empty?() // Allow retrying incase thread could not get data set fast enough
			.threadSetSessionId(ob)
		// thread should not affect main
		Assert(Database.SessionId() is: "SessionId_Test")
		Assert(ob.s1 matches: `^SessionId_Test:Thread-\d+$`)
		Assert(ob.s2 is: "thread")
		}

	threadSetSessionId(ob)
		{
		Thread()
			{
			Thread()
				{
				ob.s1 = Database.SessionId()
				Database.SessionId("thread")
				ob.s2 = Database.SessionId()
				}
			}
		Thread.Sleep(25) // want this short, but test may fail if it's too short
		}
	}
