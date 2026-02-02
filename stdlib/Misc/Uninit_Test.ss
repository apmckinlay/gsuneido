// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_global()
		{
		Assert(Uninit?('NonExistenT'))
		Assert(Uninit?('Uninit?') is: false)
		}
	Test_dynamic()
		{
		Assert(Uninit?('_nonExistenT'))
		_nonExistenT = 123
		Assert(Uninit?('_nonExistenT') is: false)
		}
	}