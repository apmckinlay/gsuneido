// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.delay, .action)
		{
		}
	Reset()
		{
		.Kill()
		.timer = Delay(.delay, .idleAfterChange)
		return this
		}
	idleAfterChange()
		{
		PrintExceptions(.action)
		}
	timer: false
	Kill()
		{
		if .timer is false
			return
		.timer.Kill()
		.timer = false
		}
	}