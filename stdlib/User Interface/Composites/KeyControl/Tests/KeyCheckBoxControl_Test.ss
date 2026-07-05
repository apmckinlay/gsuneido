// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_setStr()
		{
		m = KeyCheckBoxControl.KeyCheckBoxControl_setStr

		Assert(m('') is: 'None')
		Assert(m(#()) is: 'None')
		Assert(m(#(1)) is: 'Selected: 1')
		Assert(m(#(1, 2, 3, 4)) is: 'Selected: 4')
		}

	Test_Valid?()
		{
		mock = Mock(KeyCheckBoxControl)
		mock.When.Valid?().CallThrough()

		mock.KeyCheckBoxControl_checked = ''
		Assert(mock.Valid?())

		mock.KeyCheckBoxControl_checked = #()
		Assert(mock.Valid?())

		mock.KeyCheckBoxControl_checked = 'filled'
		Assert(mock.Valid?() is: false)
		}
	}