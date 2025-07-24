// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		code = 'function ()
	{
	Timer()
		{
		for (i = 0; i < 10000; i++)
			Assert(i is: 1)
		i = 1
		}
	}'
		Assert(FindByExpression('test.js', code, 'e = e') is: #())
		Assert({ FindByExpression('test', code, 'if e = ') } throws: 'Parse Search text')
		result = FindByExpression('test', 'function(){1+}', 'a+b')
		Assert(result isSize: 1)
		Assert(result[0] hasPrefix: 'Parse error: ')

		Assert(FindByExpression('test', code, 'e = -1') is: #())
		Assert(FindByExpression('test', code, 'e = 0') is: #((4, 5)))
		Assert(FindByExpression('test', code, 'e = e') is: #((6), (4, 5)))
		Assert(FindByExpression('test', code, 'e < e') is: #((4, 5)))

		Assert(FindByExpression('test', code, 'Assert(e is: e)') is: #((5)))
		}
	}