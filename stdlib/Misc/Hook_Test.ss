// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
// TAGS: win32
Test
	{
	Test_NoDoubleHook()
		{
		hook = Hook(WH.DEBUG, .hookProc)
		hook.Hook()
		Assert({ hook.Hook() } throws: "can't Hook() a Hook that is already hooked")
		hook.Unhook()
		}
	Test_Unhook()
		{
		hook = Hook(WH.DEBUG, .hookProc)
		Assert(hook.Unhook() is: false)
		hook.Hook()
		Assert(hook.Hooked?())
		Assert(hook.Open?())
		Assert(hook.Unhook())
		Assert(hook.Open?() is: false)
		Assert(hook.Unhook() is: false)
		}
	Test_Close()
		{
		hook = Hook(WH.DEBUG, .hookProc)
		Assert(hook.Unhook() is: false)
		hook.Hook()
		Assert(hook.Hooked?())
		Assert(hook.Open?())
		hook.Close()
		Assert(hook.Hooked?() is: false)
		Assert(hook.Open?() is: false)
		Assert(hook.Unhook() is: false)
		}
	counter: 0
	hookProc()
		{ }
	}