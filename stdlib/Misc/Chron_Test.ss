// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// TAGS: win32
Test
	{
	Test_NoDoubleStart()
		{
		chron = Chron(100, .increment)
		chron.Start()
		caught = false
		try
			chron.Start()
		catch (unused, "*already running")
			caught = true
		chron.Stop()
		Assert(caught)
		}
	Test_Stop()
		{
		chron = Chron(1000, .increment, "test stop")
		Assert(chron.Stop() is: false)
		chron.Start()
		Assert(chron.Running?())
		Assert(chron.Open?())
		Assert(chron.Stop())
		Assert(chron.Open?() is: false)
		Assert(chron.Stop() is: false)
		}
	Test_Close()
		{
		chron = Chron(1000, .increment, "test close")
		Assert(chron.Stop() is: false)
		chron.Start()
		Assert(chron.Running?())
		Assert(chron.Open?())
		chron.Close()
		Assert(chron.Running?() is: false)
		Assert(chron.Open?() is: false)
		Assert(chron.Stop() is: false)
		}
	Test_Milliseconds()
		{
		chron = Chron(500, .increment, "test milliseconds")
		Assert(chron.Milliseconds() is: 500)
		chron.SetMilliseconds(250)
		Assert(chron.Milliseconds() is: 250)
		chron.Start()
		chron.SetMilliseconds(125)
		Assert(chron.Milliseconds() is: 125)
		chron.Stop()
		}
	counter: 0
	increment()
		{ ++.counter }
	}
