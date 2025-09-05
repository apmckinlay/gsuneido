// Copyright (C) 2015 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_getLogFile()
		{
		method = Check_LogFiles.Check_LogFiles_getLogFile

		mock = Mock(Check_LogFiles)
		mock.When.getFileWithLimit('testfile').Return(false)
		Assert(mock.Eval(method, 'testfile', 'testFileName') is: '')

		mock.When.getFileWithLimit('testfile').Return(
			Object(str: '', suppressedErrors: 0))
		mock.Check_LogFiles_errorLogByteLimit = 10
		Assert(mock.Eval(method, 'testfile', 'testFileName') is: '')

		mock.When.getFileWithLimit('testfile').Return(
			Object(str: 'Error1\nError2\nError3', suppressedErrors: 0))
		mock.Check_LogFiles_errorLogByteLimit = 50
		Assert(mock.Eval(method, 'testfile', 'testFileName',
			' Log since last Nightly Checks')
			is: '\n\ntestFileName Log since last Nightly Checks:\n' $
				'Error1\nError2\nError3')

		mock.When.getFileWithLimit('testfile').Return(
			Object(str: 'Error1\nErr', suppressedErrors: 0))
		mock.When.truncatedErrors?([anyArgs:]).Return(false)
		mock.Check_LogFiles_errorLogByteLimit = 10
		Assert(mock.Eval(method, 'testfile', 'testFileName',
			' Log since last Nightly Checks')
			is: '\n\ntestFileName Log since last Nightly Checks:\n' $
				'Error1\nErr...\nTOO MANY ERRORS: LOG HAS BEEN TRIMMED')

		mock.When.getFileWithLimit('testfile').Return(
			Object(str: '                  \n\n', suppressedErrors: 0))
		mock.Check_LogFiles_errorLogByteLimit = 10
		Assert(mock.Eval(method, 'testfile', 'testFileName') is: '')
		}

	Test_getLogFile_trimmedHasError()
		{
		method = Check_LogFiles.Check_LogFiles_getLogFile

		mock = Mock(Check_LogFiles)
		mock.When.getFileWithLimit('testfile').Return(
			Object(str:
				'2017-11-02 09:23:04 INFO: closing idle connection user@127.0.0.1\r\n' $
				'2017-11-02 10:17:36 INFO: closing idle connection user@127.0.0.1\r\n' $
				'2017-11-02 12:57:10 PREVIOUS: user@127.0.0.1: fatal error: lost ' $
					'connection: socket Read failed (10053)\r\n' $
				'2017-11-02 16:14:05 INFO: closing idle connection user@127.0.0.1\r\n' $
				'2017-11-02 17:20:33 INFO: closing idle connectio', suppressedErrors: 0))
		mock.When.truncatedErrors?('testfile').Return(true)

		mock.Check_LogFiles_errorLogByteLimit = 350
		Assert(mock.Eval(method, 'testfile', 'testFileName',
			' Log since last Nightly Checks')
			is: '\n\ntestFileName Log since last Nightly Checks:\n' $
				'2017-11-02 09:23:04 INFO: closing idle connection user@127.0.0.1\r\n' $
				'2017-11-02 10:17:36 INFO: closing idle connection user@127.0.0.1\r\n' $
				'2017-11-02 12:57:10 PREVIOUS: user@127.0.0.1: fatal error: lost ' $
					'connection: socket Read failed (10053)\r\n' $
				'2017-11-02 16:14:05 INFO: closing idle connection user@127.0.0.1\r\n' $
				'2017-11-02 17:20:33 INFO: closing idle connectio...\n' $
				'TOO MANY ERRORS: LOG HAS BEEN TRIMMED\n' $
				'ERROR/WARNING found in trimmed section\n')

		mock.When.getFileWithLimit('test2file').Return(
			Object(str: 'simple test\r\n', suppressedErrors: 1))
		result = mock.Eval(method, 'test2file', 'test2FileName', ' Test Log for Testing')
		Assert(result like: 'test2FileName Test Log for Testing:\n' $
			'simple test\r\n' $
			'\nSUPPRESSIONS: 1 Non-Suneido errors were suppressed.\n')
		}

	Test_truncatedErrors?()
		{
		seekError = Check_LogFiles.Check_LogFiles_searchForError

		file1 = MockObject(
			Object(Object('Seek', 15),
			Object(Object('Readline'), result: 'This is a test line'),
			Object(Object('Readline'), result: '2017-11-02 10:18:27 PREVIOUS: ' $
				'april@192.168.0.106: fatal error: IN_PAGE_ERROR'),
			Object(Object('Readline'), result: false)))
		Assert(seekError(file1, 15) is: false)

		file2 = MockObject(
			Object(Object('Seek', 15),
			Object(Object('Readline'), result: 'This is a test line'),
			Object(Object('Readline'), result: '2017-11-02 10:15:26 ERROR: ' $
				'long duration update transaction ut2412029 (31 secs)'),
			Object(Object('Readline'), result: false)))
		Assert(seekError(file2, 15))
		}

	Test_readFileWithSuppressions()
		{
		fn = Check_LogFiles.Check_LogFiles_readFileWithSuppressions

		file = FakeFile('hello\r\n')
		result = fn(file, 15000)
		Assert(result.str is: 'hello\r\n')
		Assert(result.suppressedErrors is: 0)

		file = FakeFile('2021-07-15 07:40:03 192.168.1.30 PREVIOUS: 2021/07/14 ' $
				'07:45:01 192.168.1.30 SetTimer timeout\r\n' $
			'2021-07-15 08:14:45 192.168.1.117 PREVIOUS: (1) Appmon\r\n' $
			'2021-07-15 08:14:45 192.168.1.117 PREVIOUS: 112s359ms405us ' $
				'ERROR: SnmpChannel::Open() _ioChannel.Init() failed\r\n' $
			'2021-07-15 09:22:36 192.168.1.131 PREVIOUS: Exception 0xc0000006' $
				' 0x8 0x103480 0x103480\r\n' $
			'2021-07-15 16:47:25 dbms server: 192.168.1.124: closing idle connection\r\n'$
			'2021-07-15 17:03:59 PREV: 2021-07-15 16:36:25 Amit@wts66 ' $
				'FATAL: lost connection: EOF\r\n' $
			'2021-07-15 16:25:03 10.4.11.155 PREVIOUS: - Error during "DELETE FROM ' $
				'Types WHERE ServerID = 0 AND TypeID NOT IN (SELECT TypeID FROM ' $
				'Temp_Types)"\r\n' $
			'2021-07-16 01:24:57 INFO: SocketServer:3249 executor shutdown\r\n' $
			'2021-07-22 06:42:06 10.60.0.152 PREVIOUS:  - Code: 14, Message: ' $
				'"unable to open database file"\r\n' $
			'2021-07-15 08:19:31 10.4.11.49 PREVIOUS: "\\utilrdk\eta2\log\log.db" - ' $
				'Error during "sqlite3_open" - Code: 14, Message: "unable to ' $
				'open database file"SQLite error\r\n')
		result = fn(file, 15000)
		Assert(result.str is: '2021-07-15 07:40:03 192.168.1.30 PREVIOUS: 2021/07/14 ' $
				'07:45:01 192.168.1.30 SetTimer timeout\r\n' $
			'2021-07-15 09:22:36 192.168.1.131 PREVIOUS: Exception 0xc0000006' $
				' 0x8 0x103480 0x103480\r\n' $
			'2021-07-16 01:24:57 INFO: SocketServer:3249 executor shutdown\r\n')
		Assert(result.suppressedErrors is: 5)

		file = FakeFile('hello\r\nthis is\r\na test')
		result = fn(file, 10)
		Assert(result.str is: 'hello\r\nthis is\r\n')
		Assert(result.suppressedErrors is: 0)

		file.Seek(0)
		result = fn(file, 5)
		Assert(result.str is: 'hello\r\n')
		Assert(result.suppressedErrors is: 0)
		}

	Test_suppressionRegexes()
		{
		fn = Check_LogFiles.SuppressLine?

		testErr = ''
		Assert(fn(testErr) is: false)

		testErr = "2021-07-07 07:54:31 10.4.11.79 PREVIOUS: - Code: 14, Message: " $
			"'unable to open database file'\r\n"
		Assert(fn(testErr))

		testErr = "2021-07-07 07:54:31 10.4.11.79 PREVIOUS: - Error during " $
			"'sqlite3_open''C:\test\PCMILER30\Data\Base\save\avoidfavors.db' - " $
			"Error during 'DROP TABLE IF EXISTS Route' - Code: 8, Message: "
		Assert(fn(testErr))

		testErr = "2021-07-22 06:42:06 10.60.0.152 PREVIOUS:  - Code: 14, " $
			"Message: 'unable to open database file'"
		Assert(fn(testErr))

		testErr = '2021-06-29 09:02:42 10.100.161.15 PREVIOUS: (1) Adobe PDF Port ' $
			'Monitor\r\n'
		Assert(fn(testErr))

		testErr = '2021-06-29 09:02:42 10.100.161.15 PREVIOUS: (123456) C364SeriesFAX ' $
			'Language Monitor\r\n'
		Assert(fn(testErr))

		testErr = '2021-07-02 08:15:36 192.168.1.117 PREVIOUS: 293s572ms754us ERROR: ' $
			' FillPortInfo() USB port in a non-NT5-or-NT6 system\r\n'
		Assert(fn(testErr))

		testErr = '2021-07-02 08:15:36 192.168.1.117 PREVIOUS: 123s456ms789us ERROR:  ' $
			'GetPhysicalPort() failed\r\n'
		Assert(fn(testErr))

		testErr = '2021-07-02 08:15:36 192.168.1.117 PREVIOUS: 123ss456ms789us ERROR:  ' $
			'GetPhysicalPort() failed\r\n'
		Assert(fn(testErr) is: false)

		testErr = '2022-03-15 16:40:03 192.168.100.143 PREVIOUS: 2022-03-15T18:38:36' $
			'.011ZE [10200:NonCelloThread] crash.cc:83:HandleCrashpadLog [13672:10200:' $
			'20220315,133836.010:ERROR crashpad_client_win.cc:522] CreateProcess: ' $
			'The parameter is incorrect. (87)'
		Assert(fn(testErr), msg: 'suppress crash pad error')

		testErr = `2022-08-22 08:12:44 192.168.172.128 PREVIOUS: ` $
			`2022-08-19T14:36:31.766ZE [5968:ShellIpcClient] ` $
			`shell_ipc_client.cc:129:Connect Can't connect to socket at: ` $
			`\\.\Pipe\GoogleDriveFSPipe_lfriesen_shell`
		Assert(fn(testErr), msg: 'suppress GoogleDriveFSPipe error')

		testErr = `2022/12/28 09:17:40 PREVIOUS: 2022-12-28T14:54:15.405ZI ` $
			`[13156:NonCelloThread] ctxmenu.cc:213:GenerateContextMenu Received ` $
			`context menu with 0 menu items.`
		Assert(fn(testErr), msg: 'suppress GenerateContextMenu error')

		testErr = `2023-03-02 06:56:05 10.4.11.141 PREVIOUS: Error during` $
			` 'BEGIN EXCLUSIVE TRANSACTION' - Code: 8, Message: 'attempt to write a ` $
			`readonly database'SQLite error`
		Assert(fn(testErr), msg: 'suppress sqlite error')

		testErr = `2024/03/11 09:46:43 ERROR in rule for custom_000200 ` $
			`"SHOW: There is a problem calculating Hopper Minutes on Site\r\n` $
			`Invalid <Number> value: 1236.183333333333. ` $
			`Maximum digits before decimal is 3"`
		Assert(fn(testErr), msg: 'suppress SHOW error in rule')

		testErr = `2025/01/22 12:14:19 ERROR: in rule for custom_000097 ` $
			`"SHOW: There is a problem calculating Over Dimensional Rate\r\n` $
			`Invalid <Number> value: 130092. Maximum digits before decimal is 5"`
		Assert(fn(testErr), msg: 'suppress SHOW error: in rule')

		testErr = `2024/05/10 07:44:45 PREV: 2024/05/09 12:30:06 ` $
			`ericka@wts19 - 2024-05-09T18:30:06.667ZE [11348:NonCelloThread] ` $
			`registry_win.h:57:GetProtoFromRegistryValue Opening registry key ` $
			`Software\Google\DriveFS\Share failed with 0x2`
		Assert(fn(testErr), msg: 'suppress GetProtoFromRegistryValue')

		// not needed after BuiltDate 2025-08-21
		testErr = `(continuous_tests) 2025/08/21 09:00:06 WARNING: Query1 slow: 120 views
where view_name is "test_query"`
		Assert(fn(testErr), msg: 'suppress views slow query')
		}

	Test_skipRegexes()
		{
		fn = Check_LogFiles.Check_LogFiles_skippedLine?

		str = ''
		Assert(fn(str) is: false)

		str = '2025/01/15 15:34:50 PREV: 2025/01/15 11:55:31 raman@192.168.1.207 ' $
			'SetTimer timeout'
		Assert(fn(str) is: false, msg: 'do not skip')

		str = '2025/01/15 16:45:25 dbms server: 192.168.1.124: closing idle connection'
		Assert(fn(str), msg: 'skip idle connection closed')

		str = '2025/01/15 17:03:59 PREV: 2025/01/15 16:36:25 Amit@wts66 ' $
			'FATAL: lost connection: EOF'
		Assert(fn(str), msg: 'skip prev connection closed')
		}
	}
