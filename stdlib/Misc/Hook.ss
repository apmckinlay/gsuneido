// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Closeable
	{
	callable:  false		// any callable
	idHook:    false        // valid member of WH enumeration
	idThread:  false        // an integer in the DWORD range
	hhook:     false		// returned by SetWindowsHookEx()

	New(idHook, .callable, idThread = false)
		{
		.idHook = .checkIdHook(idHook)
		.idThread = .checkIdThread(idThread)
		}

	Close()
		{
		if false isnt .hhook
			{
			unhooked? = UnhookWindowsHookEx(.hhook)
			if not unhooked?
				throw "can't UnhookWindowsHookEx(" $ .hhook $ ")"
			cleared? = ClearCallback(.hookProc)
			if not cleared?
				throw "can't ClearCallback()"
			.hhook = false
			super.Close()
			}
		}

	Hook()
		{
		.checkNotHooked()
		hhook = SetWindowsHookEx(.idHook, .hookProc, NULL, .getThreadId())
		if NULL is hhook
			throw "can't SetWindowsHookEx()"
		.hhook = hhook
		.Closeable_open()
		return this
		}

	Unhook()
		{
		if false isnt .hhook
			{
			.Close()
			return true
			}
		return false
		}

	Hooked?()
		{ false isnt .hhook }

	checkIdHook(idHook)
		{
		Assert(WH, has: idHook)
		return idHook
		}

	checkIdThread(idThread)
		{
		rangeAbs = 2147483648
		if false isnt idThread
			Assert(idThread, isIntInRange: Object(-rangeAbs, rangeAbs))
		return idThread
		}

	checkNotHooked()
		{
		if false isnt .hhook
			throw "can't Hook() a Hook that is already hooked"
		}

	getThreadId()
		{ false is .idThread ? GetCurrentThreadId() : .idThread }

	hookProc(nCode, wParam, lParam)
		{
		// The MSDN documentation for the various hookProc callback flavours
		// invariably states: "If nCode is less than zero, the hook procedure
		// must pass the message to the CallNextHookEx function without further
		// processing and should return the value returned by CallNextHookEx."
		if 0 <= nCode
			(.callable)(hook: this, :nCode, :wParam, :lParam)
		return CallNextHookEx(.hhook, nCode, wParam, lParam)
		}
	}