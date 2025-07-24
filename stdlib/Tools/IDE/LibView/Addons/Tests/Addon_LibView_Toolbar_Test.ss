// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_buildToolBar()
		{
		m = Addon_LibView_Toolbar.Addon_LibView_Toolbar_buildToolBar
		baseToolbar = Addon_LibView_Toolbar.Addon_LibView_Toolbar_baseToolbar.Copy()

		Assert(m(#()) is: baseToolbar.Add(''))

		baseToolbar = Addon_LibView_Toolbar.Addon_LibView_Toolbar_baseToolbar.Copy()
		cmds = Object(
			#(Find, 	seq: 0)
			#(Add, 		seq: 1.0)
			#(Remove, 	seq: 1.1)
			#(Clear, 	seq: 2.1)
			)
		expected = baseToolbar.Add('', #(Find), '', #(Add), #(Remove), '', #(Clear), '')
		Assert(m(cmds) is: expected)
		}
	}