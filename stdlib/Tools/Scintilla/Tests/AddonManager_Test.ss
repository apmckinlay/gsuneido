// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Name: 'Test'
	Test_main()
		{
		mgr = AddonManager(this)
		Assert(mgr.Send(#Hello) is: false)
		Assert(mgr.Send(#Save, 'hello world'))
		Assert(.S is: 'hello world')
		}
	Addon_managertest:
	}