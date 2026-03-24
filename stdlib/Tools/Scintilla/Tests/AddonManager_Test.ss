// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Name: 'Test'
	Addon_managertest_1:,
	Addon_managertest_2:,
	Test_Send_Collect()
		{
		mgr = AddonManager(this)
		Assert(mgr.Send(#Fake) is: false)
		Assert(mgr.Send(#Save, 'hello world'))
		Assert(.S is: 'hello world')

		results = mgr.Collect(#Get)
		Assert(results isSize: 2)
		Assert(results has: 'Parent: hello world') 	// stdlib:Addon_managertest_1
		Assert(results has: 'Addon 2: hello world')	// stdlib:Addon_managertest_2

		results = mgr.Collect(#Date)
		Assert(results isSize: 1)
		Assert(results[0] isDate:) 					// stdlib:Addon_managertest_2
		}

	Test_ConditionalSend()
		{
		mgr = AddonManager(this)

		mgr.Send(#Save, #Init)
		Assert(.S is: #Init)

		mgr.ConditionalSend({|unused| false }, #(#Save, #Save1))
		Assert(.S is: #Init)

		mgr.ConditionalSend({|addon| addon.For is #Test }, #(#Save, #Save2))
		Assert(.S is: #Save2)
		}

	Test_SendToOneAddon()
		{
		mgr = AddonManager(this)
		Assert(mgr.SendToOneAddon(#Fake) is: 0)
		Assert(mgr.SendToOneAddon(#Date) isDate:)
		Assert({ mgr.SendToOneAddon(#Save) }
			throws: 'Assert FAILED: AddonManager: multiple definitions for Save')
		}

	Test_AddonMethod?()
		{
		mgr = AddonManager(this)
		Assert(mgr.AddonMethod?(#Save))
		Assert(mgr.AddonMethod?(#Date))
		Assert(mgr.AddonMethod?(#Fake) is: false)
		}
	}