// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_savePosition()
		{
		addonClass = Addon_LibView_Cursor_History
			{
			Addon_LibView_Cursor_History_maxHistory: 3
			New()
				{
				super(false, #())
				.Init()
				}
			Init()
				{
				.Addon_LibView_Cursor_History_history = Object()
				.Addon_LibView_Cursor_History_historyPos = -1
				}
			}
		addon = addonClass()
		history = addon.Addon_LibView_Cursor_History_history

		save = addon.Addon_LibView_Cursor_History_savePosition

		Assert(history is: #())
		Assert(addon.Addon_LibView_Cursor_History_historyPos is: -1)

		save('test_lib', 'test_name', 10)
		Assert(history isSize: 1)
		Assert(addon.Addon_LibView_Cursor_History_historyPos is: 0)

		.testNoDuplicate(save, history, addon)

		save('test_lib', 'test_name', 12)
		Assert(history isSize: 2)
		Assert(addon.Addon_LibView_Cursor_History_historyPos is: 1)

		save('test_lib', 'test_name2', 18)
		Assert(history isSize: 3)
		Assert(addon.Addon_LibView_Cursor_History_historyPos is: 2)

		.testMaxHistory(save, history, addon)

		.testEraseForwadLocations(addon, save, history)
		}

	testNoDuplicate(save, history, addon)
		{
		save('test_lib', 'test_name', 10)
		Assert(history isSize: 1)
		Assert(addon.Addon_LibView_Cursor_History_historyPos is: 0)
		}

	testMaxHistory(save, history, addon)
		{
		save('test_lib2', 'test_name3', 20)
		Assert(history isSize: 3)
		Assert(addon.Addon_LibView_Cursor_History_historyPos is: 2)
		}

	testEraseForwadLocations(addon, save, history)
		{
		addon.Addon_LibView_Cursor_History_historyPos = 0
		save('test_lib3', 'test_name4', 55)
		Assert(history isSize: 2)
		Assert(addon.Addon_LibView_Cursor_History_historyPos is: 1)
		}
	}
