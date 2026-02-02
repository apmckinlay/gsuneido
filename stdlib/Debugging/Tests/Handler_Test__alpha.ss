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
		_alertLogs = Object()
		_printLogs = Object()
		_debuggerLogs = Object()
		_suneidoLog = Object()
		cl = Handler
			{
			Handler_print(err)
				{
				_printLogs.Add(err)
				}
			Handler_debuggerWindow(hwnd, err, calls)
				{
				_debuggerLogs.Add([:hwnd, :err, :calls, onDestroy: false])
				}
			Handler_alert(err, hwnd, log_rec, logOnly? = false)
				{
				if not logOnly?
					_alertLogs.Add([msg: err, :hwnd, :log_rec])
				_suneidoLog.Add([sulog_message: log_rec.sulog_message,
					calls: log_rec.calls])
				}
			}

		// ------------ Default user errors ------------
		Suneido.User = 'default'

		// Testing: ''
		cl('', 1, [])
		Assert(_debuggerLogs.Extract(0)
			equalsSet: [err: '', hwnd: 1, calls: [], onDestroy: false])

		// Testing: standard error
		cl('error is thrown', 1, [])
		Assert(_debuggerLogs.Extract(0)
			equalsSet: [err: 'error is thrown', hwnd: 1, calls: [], onDestroy: false])

		// Testing: interrupt errors
		cl('interrupt - error is printed', 1, [])
		Assert(_printLogs.Extract(0) is: "interrupt - error is printed")

		// ------------ System user errors ------------
		Suneido.User = 'admin'

		// Testing: ''
		cl('', 1, [])
		Assert(_suneidoLog isSize: 0)

		// Testing: show errors on client side (produces a "Warning" alert)
		cl('SHOW: this error to user', 1, [arg: 'carried over'])
		Assert(_suneidoLog isSize: 1)
		Assert(_suneidoLog[0] equalsSet:
			#(sulog_message: 'warning: SHOW: this error to user',
				calls: #(arg: 'carried over')))
		Assert(_alertLogs isSize: 1)
		Assert(_alertLogs[0] equalsSet: Object(msg: 'this error to user', hwnd: 1,
			log_rec: #(sulog_message: 'warning: SHOW: this error to user',
			calls: #(arg: 'carried over'))))

		// Testing: show errors from server side (produces a "Warning" alert)
		cl('SHOW: this error to user (from server)', 1, [arg: 'carried over'])
		Assert(_suneidoLog isSize: 2)
		Assert(_suneidoLog[1] equalsSet:
			#(sulog_message: 'warning: SHOW: this error to user (from server)',
				calls: #(arg: 'carried over')))
		Assert(_alertLogs isSize: 2)
		Assert(_alertLogs[1] equalsSet: Object(msg: 'this error to user', hwnd: 1,
			log_rec: #(sulog_message: 'warning: SHOW: this error to user (from server)',
			calls: #(arg: 'carried over'))))

		// Testing: normal errors (Alert: An unexpected problem has occurred)
		cl('unexpected error', 1, [arg: 'carried over'])
		Assert(_suneidoLog isSize: 3)
		Assert(_suneidoLog[2] equalsSet:
			#(sulog_message: 'ERROR: unexpected error',
				calls: #(arg: 'carried over')))
		Assert(_alertLogs isSize: 3)
		Assert(_alertLogs[2] equalsSet: Object(msg: false, hwnd: 1,
			log_rec: #(sulog_message: 'ERROR: unexpected error',
			calls: #(arg: 'carried over'))))

		// Testing: interrupt errors (produces no alerts)
		cl('interrupt error', 1, [arg: 'carried over'])
		Assert(_suneidoLog isSize: 4)
		Assert(_suneidoLog[3] equalsSet:
			#(sulog_message: 'interrupt error',
				calls: #(arg: 'carried over')))
		Assert(_alertLogs isSize: 3)
		Assert(_alertLogs isSize: 3)
		}

	Teardown()
		{
		Suneido.User = .suneidoUser
		super.Teardown()
		}
	}