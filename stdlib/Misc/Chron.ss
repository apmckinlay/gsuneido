// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Closeable
	{
	milliseconds: false
	callable:     false
	id:           false  // any value
	timer:        false  // identifier returned by SetTimer

	New(milliseconds, .callable, .id = false)
		{
		.SetMilliseconds(milliseconds)
		}

	Close()
		{
		if false isnt .timer
			{
			killed? = KillTimer(NULL, .timer)
			if not killed?
				throw "can't KillTimer(" $ .timer $ ")"
			cleared? = ClearCallback(.timerProc)
			if not cleared?
				throw "can't ClearCallback()"
			.timer = false
			super.Close()
			}
		return // Don't return a value
		}

	Start()
		{
		.checkNotRunning()
		timer = SetTimer(NULL, 0, .milliseconds, .timerProc)
		if 0 is timer
			throw "can't SetTimer()"
		.timer = timer
		.Closeable_open()
		return this
		}
	Stop()
		{
		if false isnt .timer
			{
			.Close()
			return true
			}
		return false
		}
	Milliseconds()
		{
		.milliseconds
		}
	Id()
		{
		.id
		}
	Running?()
		{
		false isnt .timer
		}
	MinMilliseconds()
		{
		USER_TIMER_MINIMUM
		}
	MaxMilliseconds()
		{
		USER_TIMER_MAXIMUM
		}
	SetMilliseconds(milliseconds)
		{
		// NOTE: This does not change the interval on running timer. If the
		//       timer is currently running, you will have to .Stop(), .Start()
		//       it to benefit from the change.
		.milliseconds = .checkMilliseconds(milliseconds)
		return
		}

	checkMilliseconds(milliseconds)
		{
		Assert(milliseconds, isInt:)
		if milliseconds < USER_TIMER_MINIMUM
			throw "can't set timer to less than " $ .MinMilliseconds() $ " milliseconds"
		else if USER_TIMER_MAXIMUM < milliseconds
			throw "can't set timer to more than " $ .MaxMilliseconds() $ " milliseconds"
		return milliseconds
		}
	checkNotRunning()
		{
		if false isnt .timer
			throw "can't Start() a Chron that is already running"
		}
	timerProc(hwnd/*unused*/, msg/*unused*/, id/*unused*/, time/*unused*/)
		{
		(.callable)(chron: this)
		}
	}
