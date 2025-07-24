// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ValidData?()
		{
		args = Object('', list: Object('a', 'b'))
		Assert(ChooseManyControl.ValidData?(@args))

		args.saveNone = true
		Assert(ChooseManyControl.ValidData?(@args) is: false)

		args.mandatory = true
		Assert(ChooseManyControl.ValidData?(@args) is: false)

		args.Delete('saveNone')
		Assert(ChooseManyControl.ValidData?(@args) is: false)

		args[0] = 'None'
		Assert(ChooseManyControl.ValidData?(@args) is: false)

		args.saveNone = true
		Assert(ChooseManyControl.ValidData?(@args))

		args[0] = 'f'
		Assert(ChooseManyControl.ValidData?(@args) is: false)

		args[0] = 'a'
		Assert(ChooseManyControl.ValidData?(@args))

		args[0] = 'a,b'
		Assert(ChooseManyControl.ValidData?(@args))

		args[0] = 'a, b'
		Assert(ChooseManyControl.ValidData?(@args))

		args[0] = 'a, None'
		Assert(ChooseManyControl.ValidData?(@args) is: false)
		}
	}