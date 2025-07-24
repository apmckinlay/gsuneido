// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Type(ServerEval('Thread.List')) is: 'Object') // ServerEval does @+1 args
		Assert(Type(ServerEval(@#('Thread.List'))) is: 'Object')
		Assert(Type(ServerEval(@+1#(123, 'Thread.List'))) is: 'Object')
		}
	}