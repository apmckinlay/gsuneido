// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(LibCanonicalName('hello') is: 'hello')
		Assert(LibCanonicalName('Hello_Test') is: 'hellotest')
		}
	}