// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		.lib = .TestLibName()
		.name = .TempName()
		s =  "function (a, b)
	{
	c = a + b
	c *= 2
	return c
	}"
		.MakeLibraryRecord([name: .name, text: s])

		spy = .SpyOn(SteppingDebugger).Return('')
		callLogs = spy.CallLogs()

		// source should not exist
		Assert(SteppingDebuggerManager().GetSource(.lib, .name) is: false)

		// cannot add break point at pos 0 (not a statement)
		Assert(SteppingDebuggerManager().ToggleBreakPoint(.lib, .name, s, 0) is: false)
		Assert(SteppingDebuggerManager().HasBreakPoints?(.lib, .name) is: false)
		// source should be initialized
		Assert(SteppingDebuggerManager().GetSource(.lib, .name) is: s)

		// add break point on "c = a + b"
		Assert(SteppingDebuggerManager().ToggleBreakPoint(.lib, .name, s, 29))
		Assert(SteppingDebuggerManager().HasBreakPoints?(.lib, .name))
		// add break point on "c *= 2"
		Assert(SteppingDebuggerManager().ToggleBreakPoint(.lib, .name, s, 34))
		Assert(SteppingDebuggerManager().GetBreakPointRanges(.lib, .name)
			is: Object(Object(i: 22, n: 9), Object(i: 34, n: 6)))

		Assert(SteppingDebuggerManager().ToggleBreakPoint(.lib, .name, s, 0) is: false)
		Assert(SteppingDebuggerManager().ToggleBreakPoint(.lib, .name, s, 21) is: false)
		Assert(SteppingDebuggerManager().ToggleBreakPoint(.lib, .name, s, 31) is: false)
		// remove break point on "c = a + b"
		Assert(SteppingDebuggerManager().ToggleBreakPoint(.lib, .name, s, 30))
		Assert(SteppingDebuggerManager().HasBreakPoints?(.lib, .name))
		// remove break point on "c *= 2"
		Assert(SteppingDebuggerManager().ToggleBreakPoint(.lib, .name, s, 34))
		Assert(SteppingDebuggerManager().HasBreakPoints?(.lib, .name) is: false)
		Assert(SteppingDebuggerManager().GetBreakPointRanges(.lib, .name) is: #())

		// add break point on "c = a + b"
		Assert(SteppingDebuggerManager().ToggleBreakPoint(.lib, .name, s, 29))

		fn = Global(.name)
		Assert(fn(1, 2) is: 6)
		Assert(callLogs.Size() is: 3)

		// remove all break points
		SteppingDebuggerManager().RemoveAllBreakPoints(.lib, .name)
		Assert(SteppingDebuggerManager().HasBreakPoints?(.lib, .name) is: false)

		callLogs.Delete(all:)
		fn = Global(.name)
		Assert(fn(1, 2) is: 6)
		Assert(callLogs.Size() is: 0)
		}

	Teardown()
		{
		SteppingDebuggerManager().ClearSource(.lib, .name)
		super.Teardown()
		}
	}
