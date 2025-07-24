// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(msg = '', title = 'Warning', hwnd = 0, flags = false, uniqueId = false,
		block = false)
		{
		// Alert has been overridden, skip the delay e.g. DoWithAlertRedirected
		if false isnt Suneido.GetDefault(#Alert, false)
			{
			.callAlert(msg, title, hwnd, flags, block)
			return
			}

		.uniqueSetup(uniqueId)

		if flags is false
			flags = MB.ICONWARNING
		timer = Defer({ .alert(msg, :title, :hwnd, :flags, :uniqueId, :block) })
		.uniqueSet(timer, uniqueId)
		}

	uniqueSetup(uniqueId)
		{
		if uniqueId is false
			return
		timers = Suneido.GetInit(#alertDelayedTimers, { Object() })
		if timers.Member?(uniqueId)
			timers[uniqueId].Kill()
		}

	uniqueSet(timer, uniqueId)
		{
		if uniqueId is false
			return
		Suneido.alertDelayedTimers[uniqueId] = timer
		}

	alert(msg, title, hwnd, flags, uniqueId = false, block = false)
		{
		.uniqueClear(uniqueId)
		if not IsWindow(hwnd)
			hwnd = 0

		.callAlert(msg, title, hwnd, flags, block)
		}

	callAlert(msg, title, hwnd, flags, block)
		{
		if block is false
			Alert(:msg, :title, :hwnd, :flags)
		else
			block(Alert(:msg, :title, :hwnd, :flags))
		}

	uniqueClear(uniqueId)
		{
		if uniqueId is false
			return
		Suneido.alertDelayedTimers.Delete(uniqueId)
		}

	}
