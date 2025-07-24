// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(milliseconds, block, list = false)
		{
		return new this(milliseconds, block, list)
		}

	New(milliseconds, callable, .list)
		{
		.callable = callable
		.timer = SuSetTimer(NULL, 0, milliseconds, this)

		if list isnt false
			list.Add(this)
		}
	Call(hwnd/*unused*/, msg/*unused*/, id/*unused*/, time/*unused*/)
		{
		if .Kill()
			(.callable)()
		return // need this to ensure no return value
		}
	timer: 0
	Kill()
		{
		timer = .timer
		.timer = 0
		if timer is 0
			return false

		SuKillTimer(NULL, timer)
		if .list isnt false
			.list.Delete(.list.Find(this))
		return true
		}
	}
