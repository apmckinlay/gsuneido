// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_concatenate_should_preserve_exception()
		{
		try
			throw 'foo'
		catch (e)
			{
			Assert(e $ 'bar' is: 'foobar')
			Assert('kung' $ e is: 'kungfoo')
			Assert(Type(e $ 'oops') is 'Except')
			Assert(Type('oops' $ e) is 'Except')
			(e $ 'oops').Callstack()
			('oops' $ e).Callstack()
			}
		}
	}