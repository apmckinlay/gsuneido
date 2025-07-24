// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		cl_linux = HandleSystemWait
			{
			HandleSystemWait_linux?()
				{ return true }
			}
		cmd = cl_linux('test command', wait?: false)
		Assert(cmd is: 'test command &')
		cmd = cl_linux('test command', wait?:)
		Assert(cmd is: 'test command')

		cl_windows = HandleSystemWait
			{
			HandleSystemWait_linux?()
				{ return false }
			}
		cmd = cl_windows('test command', wait?: false)
		Assert(cmd is: 'start test command')
		cmd = cl_windows('test command', wait?:)
		Assert(cmd is: 'start /w test command')
		}
	}