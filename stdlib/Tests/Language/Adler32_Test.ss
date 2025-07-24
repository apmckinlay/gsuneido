// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(Adler32("") is: 1)
		Assert(Adler32("\x00\xff") is: 16843008)

		cksum = 389415997
		Assert(Adler32("helloworld") is: cksum)
		Assert(Adler32("hello", "world") is: cksum)
		Assert(Adler32().Update("hello").Update("world").Value() is: cksum)
		}
	}