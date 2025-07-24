// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	tests: #((
		src: 'dll long User32:GetSystemMetrics(long nIndex)',
		expected: 'dll long User32:GetSystemMetrics(long nIndex)'), (
		src: '#(test: 1)',
		expected: '#(test: 1)'), (
		src: 'function (a, b = 1)
	{
	c = a + b
	return c
	}',
		expected: 'function (a, b = 1)
	{
	SteppingDebugger(0);
	c = a + b
	SteppingDebugger(1);
	return c
	}'), (
		src: 'class
	{
	A: #(fn: function (a, b) { return a + b })
	Test(a, b)
		{
		c = (.A.fn)(a, b)
		if a > 0
			if b > 0
				{
				c++
				}
		c = c +
			2
		try
			c = Test(c)
		catch (e)
			{
			c = 0
			for i in .. 10
				c += i
			}
		return c > 50
			? 1
			: 0
		}
	}',
		expected: 'class
	{
	A: #(fn: function (a, b) { SteppingDebugger(0);
		return a + b })
	Test(a, b)
		{
		SteppingDebugger(1);
		c = (.A.fn)(a, b)
		SteppingDebugger(2);
		if a > 0
			{ SteppingDebugger(6); if b > 0
				{
				SteppingDebugger(7);
				c++
				} }
		SteppingDebugger(3);
		c = c +
			2
		SteppingDebugger(4);
		try
			{ SteppingDebugger(8); c = Test(c) }
		catch (e)
			{
			SteppingDebugger(9);
			c = 0
			SteppingDebugger(10);
			for i in .. 10
				{ SteppingDebugger(11); c += i }
			}
		SteppingDebugger(5);
		return c > 50
			? 1
			: 0
		}
	}'), (
		src: 'function ()
	{
	block = { it + 1 }
	return block
	}',	expected: 'function ()
	{
	SteppingDebugger(0);
	block = { SteppingDebugger(2);
		it + 1 }
	SteppingDebugger(1);
	return block
	}'), (
		src: 'function (a)
	{
	switch (a)
		{
	case 1:
		return a + 1
	case 2:
		a = a + 2
		return a
	default:
		return a
		}
	}', expected: 'function (a)
	{
	SteppingDebugger(0);
	switch (a)
		{
	case 1:
		SteppingDebugger(1);
		return a + 1
	case 2:
		SteppingDebugger(2);
		a = a + 2
		SteppingDebugger(3);
		return a
	default:
		SteppingDebugger(4);
		return a
		}
	}'), (
		src: 'class
	{
	New()
		{
		super()
		.test()
		}
	test()
		{
		return 1
		}
	}', expected: 'class
	{
	New()
		{
		super()
		SteppingDebugger(0);
		.test()
		}
	test()
		{
		SteppingDebugger(1);
		return 1
		}
	}'))

	Test_main()
		{
		num = 0
		registerDebuggerFn = { |unused| num++ }
		for test in .tests
			{
			num = 0
			root = Tdop(test.src)
			writeMgr = AstWriteManager(test.src, root)
			Assert(AddDebugger(test.src, root, writeMgr, registerDebuggerFn)
				like: test.expected)
			}
		}
	}
