// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidData?()
		{
		vd = RadioButtonsControl.ValidData?

		Assert({ vd() } throws: 'missing: 0')
		args = Object()
		Assert({ vd(@args) }throws: 'missing: 0')

		args = Object('test1', 'test1', 'test2')
		Assert(vd(@args))

		args = Object('test1', 'test2')
		Assert(vd(@args) is: false)

		args = Object('','test1', 'test2')
		Assert(vd(@args))

		args = Object('', 'test1', 'test2', mandatory:)
		Assert(vd(@args) is: false)
		}
	}