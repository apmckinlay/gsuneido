// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		// shouldn't give program error, see Close method for explanation
		WebSession.Close('uuid_biz_tests')

		// test normal case
		uuid = WebSession.Register('biz_tests')
		Assert(WebSession.Authenticate(uuid) is: 'biz_tests')
		WebSession.Close(uuid)
		Assert(WebSession.Authenticate(uuid) is: false)

		// additional close shouldn't give program error
		WebSession.Close(uuid)
		}

	Test_Authenticate()
		{
		session = new WebSession()
		session.WebSession_sessionMember = 'biztest_session'
		Assert(session.Authenticate('uuid_biz_tests') is: false)
		}
	}