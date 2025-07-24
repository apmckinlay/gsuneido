// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_findPrevLines()
		{
		getChange = Addon_show_modified_lines.Addon_show_modified_lines_getChangedLines
		findPrevLines = Addon_show_modified_lines.Addon_show_modified_lines_findPrevLines

		addon = Mock(Addon_show_modified_lines)
		addon.When.GetSelect().Return(#(cpMin: 0, cpMax: 1))
		addon.When.LineFromPosition([anyArgs:]).Return(0)
		addon.When.PositionFromLine([anyArgs:]).Return(0)

		// testing modification
		addon.Eval(getChange, beforeText: 'a', text: 'b')
		prev = addon.Eval(findPrevLines)
		Assert(prev.prevLines is: 'a\n')
		Assert(prev.localLineHeight is: 1)

		// testing deletion
		addon.Eval(getChange, beforeText: 'a\nb\nc', text: 'a\nc')
		prev = addon.Eval(findPrevLines)
		Assert(prev.prevLines is: 'a\nb\n')
		Assert(prev.localLineHeight is: 1)

		// testing addition
		addon.Eval(getChange, beforeText: 'b\nc', text: 'a\nb\nc')
		prev = addon.Eval(findPrevLines)
		Assert(prev.prevLines is: '')
		Assert(prev.localLineHeight is: 1)
		}

	Test_getRGBString()
		{
		m = Addon_show_modified_lines.Addon_show_modified_lines_getRGBString
		Assert(m(0x12010f) is: '0f0112')
		Assert(m(0x000000) is: '000000')
		Assert(m(0x010203) is: '030201')
		Assert(m(0xabcdef) is: 'efcdab')
		}

	Test_validateParams()
		{
		mock = Mock(Addon_show_modified_lines)
		mock.When.validateParams().CallThrough()

		mock.When.Send([anyArgs:]).
			Throw('thread failed to connect to db server: socket connection timeout')
		Assert(mock.validateParams() is: false)

		mock.When.Send([anyArgs:]).Throw('uncaught error')
		Assert({ mock.validateParams() } throws: 'uncaught error')

		mock.When.Send([anyArgs:]).Return('name', 'library')
		Assert(mock.validateParams())
		Assert(mock.Addon_show_modified_lines_name is: 'name')
		Assert(mock.Addon_show_modified_lines_table is: 'library')
		}
	}