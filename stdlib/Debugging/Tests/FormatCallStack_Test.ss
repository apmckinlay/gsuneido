// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Format()
		{
		Assert(FormatCallStack(#()) is: "")
		Assert(FormatCallStack(#((fn: 123))) is: "123")
		Assert(FormatCallStack(#((fn: 123)(fn: 456))) is: "123\n456")
		Assert(FormatCallStack(#((fn: 123)(fn: 456)(fn: 789)), 2) is: "123\n456")
		Assert(FormatCallStack(#((fn: 123)(fn: 456)(fn: 789)), levels: 2) is: "123\n456")
		}
	}