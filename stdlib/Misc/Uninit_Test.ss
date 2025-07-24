// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_global()
		{
		Assert(Uninit?('NonExistenT') is: true)
		Assert(Uninit?('Uninit?') is: false)
		}
	Test_dynamic()
		{
		Assert(Uninit?('_nonExistenT') is: true)
		_nonExistenT = 123
		Assert(Uninit?('_nonExistenT') is: false)
		}
	}