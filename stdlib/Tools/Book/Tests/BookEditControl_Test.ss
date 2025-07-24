// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Goto()
		{
		c = BookEditControl
			{
			BookEditControl_table: 'TestHelp'
			Getter_Editor()
				{
				return class
					{
					GetSelText()
						{
						return _testCases[0]
						}
					}
				}
			Getter_Explorer()
				{
				return class
					{
					GotoPath(text)
						{
						_testCases[1] = text
						}
					}
				}
			}
		.testGoto(c, 'Shared', 'TestHelp/res/Shared')
		.testGoto(c, 'abc.png', 'TestHelp/res/abc.png')
		.testGoto(c, '/Setup/Page', 'TestHelp/Setup/Page')
		.testGoto(c, '/TestHelp/Setup/Page', 'TestHelp/Setup/Page')
		.testGoto(c, 'TestHelp/Setup/Page', 'TestHelp/Setup/Page')
		}

	testGoto(c, sel, goto)
		{
		_testCases = Object(sel)
		c.On_Goto()
		Assert(_testCases[1] is: goto)
		}

	Test_FindLinkedHelpPage()
		{
		c = Mock(BookEditControl)
		c.Editor = Mock()
		c.Editor.When.GetCurrentPos().Return(
			0
			1
			2
			38  // href
			42  // Help
			64	// General
			65	// Test
			71	// Test2
			96 	// tab
			111 // Screenshot
			120 // help.png
			)
		c.Editor.When.GetSelect().Return(
			Object(cpMin: 0, cpMax: 0),
			Object(cpMin: 1, cpMax: 1),
			Object(cpMin: 1, cpMax: 3),
			Object(cpMin: 34, cpMax: 38) // href
			Object(cpMin: 41, cpMax: 45) // Help
			Object(cpMin: 48, cpMax: 48) // General
			Object(cpMin: 64, cpMax: 68) // Test
			Object(cpMin: 70, cpMax: 75) // Test2
			Object(cpMin: 96, cpMax: 99) // tab
			Object(cpMin: 108, cpMax: 118) // Screenshot
			Object(cpMin: 120, cpMax: 128) // help.png
			)
		c.Editor.When.GetAt([anyArgs:]).Do({ |call| _text[call[1]] })
		c.Editor.When.GetRange([anyArgs:]).Do({ |call| _text[call[1] .. call[2]] })
		c.When.findLink([anyArgs:]).CallThrough()

		_text = `<h3><a onClick="showhide('flip')" href="/Help/General/Reference/Test">` $
			"Test2</a> &raquo; General tab</h3>\n<$ Screenshot('help.png') $>"

		Assert(c.Eval(BookEditControl.FindLinkedHelpPage) is: '')
		Assert(c.Eval(BookEditControl.FindLinkedHelpPage) is: '')
		Assert(c.Eval(BookEditControl.FindLinkedHelpPage) is: '')
		Assert(c.Eval(BookEditControl.FindLinkedHelpPage) is: '')  			// href

		Assert(c.Eval(BookEditControl.FindLinkedHelpPage)
			is: '/Help/General/Reference/Test')								// Help
		Assert(c.Eval(BookEditControl.FindLinkedHelpPage)
			is: '/Help/General/Reference/Test')								// General
		Assert(c.Eval(BookEditControl.FindLinkedHelpPage)
			is: '/Help/General/Reference/Test')								// Test
		Assert(c.Eval(BookEditControl.FindLinkedHelpPage) is: '') 			// Test2
		Assert(c.Eval(BookEditControl.FindLinkedHelpPage) is: '') 			// tab
		Assert(c.Eval(BookEditControl.FindLinkedHelpPage) is: '')  			// Screenshot
		Assert(c.Eval(BookEditControl.FindLinkedHelpPage) is: 'help.png') 	// help.png
		}
	}