// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidateName()
		{
		m = MultiDrawControl.ValidateName
		Assert(m('', []) is: 'Name can not be blank or contain symbols')
		Assert(m('Test!', []) is: 'Name can not be blank or contain symbols')
		Assert(m('Test',  []) is: '')
		Assert(m('Test1', ['Test']) is: '')
		Assert(m('Test',  ['Test'])
			is: 'Test already exists.\r\nPlease enter another name.')
		}
	}