// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_differentMatchRecord?()
		{
		different? = IdFieldControl.IdFieldControl_differentMatchRecord?
		Assert(different?(false, [], []) is: false)

		ob = Object()
		ob.IdFieldControl_numField = 'f1'
		Assert(ob.Eval(different?, true, [], []) is: false)
		Assert(ob.Eval(different?, true, [f1: 'a'], [f1: 'a']) is: false)
		Assert(ob.Eval(different?, true, [f1: 'a'], [f1: 'b']))
		}
	}
