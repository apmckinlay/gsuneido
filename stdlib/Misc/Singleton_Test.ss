// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		c = Singleton_TestClass
		Assert(c() is: Suneido.Singleton_TestClass)
		Assert(c().N is: 0)
		c().Inc()
		Assert(c().N is: 1)
		c.Reset()
		Assert(Suneido.Member?('Singleton_TestClass') is: false)
		Assert(c().N is: 0)
		}
	Teardown()
		{
		Suneido.Delete('Singleton_TestClass')
		}
	}