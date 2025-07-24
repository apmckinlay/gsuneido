// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		proxy = ServerEvalProxy("ServerEvalProxy_Test") // self shunt
		Assert(proxy.A() is: 123)
		Assert(proxy.B(20, "th century ", #20000101) is: "20th century 2000-1-1")
		}
	A()
		{ return 123 }
	B(number, string, date)
		{ return number $ string $ date.Format('yyyy-M-d') }
	}