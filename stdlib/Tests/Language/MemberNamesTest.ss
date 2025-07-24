// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// NOTE: deliberately NOT using is: because that's what we're testing
Test
	{
	Test_symbol()
		{
		Assert(#is is 'is')
		Assert(#isnt is 'isnt')
		Assert(#and is 'and')
		Assert(#or is 'or')
		Assert(#not is 'not')
		Assert(#for is 'for')
		}
	Test_object_constant()
		{
		Assert(#(is: isnt) is #('is': 'isnt'))
		Assert(#(for: if) is #('for': 'if'))
		}
	Test_arguments()
		{
		Assert(Object(is: 123) is #('is': 123))
		Assert(Object(for: 99) is #('for': 99))
		}
	}