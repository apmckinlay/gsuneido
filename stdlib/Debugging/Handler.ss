// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(err, hwnd, calls)
		{
		if hwnd is 0
			hwnd = GetActiveWindow()
		if Suneido.User is 'default'
			.programmerError(err, hwnd, calls)
		else
			.userError(err, hwnd, calls)
		}

	programmerError(err, hwnd, calls)
		{
		if .interrupt?(err)
			Print(err)
		else
			Debugger.Window(hwnd, err, calls)
		}

	interrupt?(err)
		{
		return err.Prefix?('interrupt')
		}

	showPrefix: 'SHOW: '
	serverSuffix: ' (from server)'
	userError(err, hwnd, calls)
		{
		if err is ''
			return
		logPrefix = (show? = err.Prefix?(.showPrefix))
			? 'warning: '
			: .interrupt?(err)
				? ''
				: 'ERROR: '
		log_rec = SuneidoLog(logPrefix $ err, calls)
		err = err.RemovePrefix(.showPrefix).RemoveSuffix(.serverSuffix)
		if not .interrupt?(err)
			.alertError(show?, err, log_rec, hwnd)
		}

	alertError(show?, err, log_rec, hwnd)
		{
		if show?
			log_rec = #()
		else
			err = false
		AlertError(err, :hwnd, :log_rec)
		}
	}
