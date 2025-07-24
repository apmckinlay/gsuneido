// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.file = .MakeFile('one\ntwo\r\nthree\r\r\n' $
			'helloworld'.Repeat(500) $ '\r\n' $ 'bye')
		File(.file)
			{
			.test(it)
			}
		}
	test(f)
		{
		Assert(f.Readline() is: 'one')
		Assert(f.Readline() is: 'two')
		Assert(f.Readline() is: 'three')
		s = f.Readline()
		Assert(s isSize: 4000)
		Assert(s is: 'helloworld'.Repeat(400))
		// got 4000, remainder of line discarded
		Assert(f.Readline() is: 'bye')
		Assert(f.Readline() is: false)
		}
	}
