// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_runBlock()
		{
		runblockfn = LibViewCoreControl.LibViewCoreControl_runBlock

		result = runblockfn({ ['test0 passed'] })
		Assert(result.blockResult is: #("test0 passed"))
		Assert(result.err is: '')

		result = runblockfn({ throw 'test1 failed' })
		Assert(result.blockResult is: #(), msg: 'test1 result.blockResult')
		Assert(result.err is: 'test1 failed')

		add = 0
		result = runblockfn({ add += 25; [] })
		Assert(add is: 25)
		Assert(result.blockResult is: #())
		Assert(result.err is: '')
		}
	}
