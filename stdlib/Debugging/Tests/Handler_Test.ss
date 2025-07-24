// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()
		.suneidoUser = Suneido.User
		}

	Test_main()
		{
		logs = Object()
		logs.Add(printLogs = .SpyOn(Print).Return('').CallLogs())
		logs.Add(debuggerLogs = .SpyOn(Debugger.Window).Return('').CallLogs())
		logs.Add(suneidoLogs = .SpyOn(SuneidoLog).Return([suneidolog:]).CallLogs())
		logs.Add(alertLogs = .SpyOn(AlertError).Return('').CallLogs())

		// ------------ Default user errors ------------
		Suneido.User = 'default'

		// Testing: ''
		Handler('', 1, [])
		Assert(debuggerLogs.Extract(0)
			equalsSet: [err: '', hwnd: 1, calls: [], onDestroy: false])
		logs.Each({ Assert(it isSize: 0) })

		// Testing: standard error
		Handler('error is thrown', 1, [])
		Assert(debuggerLogs.Extract(0)
			equalsSet: [err: 'error is thrown', hwnd: 1, calls: [], onDestroy: false])
		logs.Each({ Assert(it isSize: 0) })

		// Testing: interrupt errors
		Handler('interrupt - error is printed', 1, [])
		Assert(printLogs.Extract(0) is: #(args: #("interrupt - error is printed")))
		logs.Each({ Assert(it isSize: 0) })

		// ------------ System user errors ------------
		Suneido.User = 'admin'

		// Testing: ''
		Handler('', 1, [])
		logs.Each({ Assert(it isSize: 0) })

		// Testing: show errors on client side (produces a "Warning" alert)
		Handler('SHOW: this error to user', 1, [arg: 'carried over'])
		Assert(suneidoLogs.Extract(0)
			equalsSet: #(
				message: 'warning: SHOW: this error to user',
				calls: (arg: 'carried over'),
				params: '', switch_prefix_limit: 10, caughtMsg: ''))
		Assert(alertLogs.Extract(0)
			equalsSet: #(msg: 'this error to user', hwnd: 1, log_rec: #()))
		logs.Each({ Assert(it isSize: 0) })

		// Testing: show errors from server side (produces a "Warning" alert)
		Handler('SHOW: this error to user (from server)', 1, [arg: 'carried over'])
		Assert(suneidoLogs.Extract(0)
			equalsSet: #(
				message: 'warning: SHOW: this error to user (from server)',
				calls: (arg: 'carried over'),
				params: '', switch_prefix_limit: 10, caughtMsg: ''))
		Assert(alertLogs.Extract(0)
			equalsSet: #(msg: 'this error to user', hwnd: 1, log_rec: #()))
		logs.Each({ Assert(it isSize: 0) })

		// Testing: normal errors (Alert: An unexpected problem has occurred)
		Handler('unexpected error', 1, [arg: 'carried over'])
		Assert(suneidoLogs.Extract(0)
			equalsSet: #(
				message: 'ERROR: unexpected error',
				calls: (arg: 'carried over'),
				params: '', switch_prefix_limit: 10, caughtMsg: ''))
		Assert(alertLogs.Extract(0)
			equalsSet: #(msg: false, hwnd: 1, log_rec: (suneidolog:)))
		logs.Each({ Assert(it isSize: 0) })

		// Testing: interrupt errors (produces no alerts)
		Handler('interrupt error', 1, [arg: 'carried over'])
		Assert(suneidoLogs.Extract(0)
			equalsSet: #(
				message: 'interrupt error',
				calls: (arg: 'carried over'),
				params: '', switch_prefix_limit: 10, caughtMsg: ''))
		logs.Each({ Assert(it isSize: 0) })
		}

	Teardown()
		{
		Suneido.User = .suneidoUser
		super.Teardown()
		}
	}