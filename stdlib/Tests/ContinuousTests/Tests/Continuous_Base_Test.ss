// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup ()
		{
		super.Setup()
		.exeFolderExist? = DirExists?('exeFolder')
		}
	Test_htmlFormatResults()
		{
		cl = ContinuousTest_Base
			{
			ContinuousTest_Base_minTestResultSize: 3
			GetChangeFileText()
				{ return '' }
			}
		formatResult = cl.ContinuousTest_Base_htmlFormatResults

		noErrorMsg = 'System Tests as of 2016-11-17 09:24
TestThatIsFine_Test_fine
		Pass Message goes here:
		Duration of Test: .5 sec'
		formattedNoErrorMsg = '<h3>Continuous Tests Results</h3>' $
		'System Tests as of 2016-11-17 09:24
TestThatIsFine_Test_fine
		Pass Message goes here:
		Duration of Test: .5 sec'

		text = formatResult(noErrorMsg).text
		Assert(text like: formattedNoErrorMsg)


		legitErrorMsg = 'System Tests as of 2016-11-17 09:24
ThisIsATest.Test_willhaveissue
		ERROR: this is an error message
		Duration of Test: .5 sec'

		formattedLegitErrorMsg = '<h3>Continuous Tests Results</h3>' $
		'System Tests as of 2016-11-17 09:24
ThisIsATest.Test_willhaveissue
<span style="color: red">        ERROR: this is an error message</span>
		Duration of Test: .5 sec'

		text = formatResult(legitErrorMsg).text
		Assert(text like: formattedLegitErrorMsg)

		noErrorWithFailWord = 'System Tests as of 2016-11-17 09:24
ThisShouldNotFail_Test.Test_Error
		Pass Message goes here:
		Duration of Test: .5 sec'

		formattedNoErrorWithFailWord = '<h3>Continuous Tests Results</h3>' $
		'System Tests as of 2016-11-17 09:24
ThisShouldNotFail_Test.Test_Error
		Pass Message goes here:
		Duration of Test: .5 sec'

		text = formatResult(noErrorWithFailWord).text
		Assert(text like: formattedNoErrorWithFailWord)


		legitFailureMsg = 'ETA Tests - With Test Data
Build date: Dec 19 2016 12:34:23 (vs2015 Release)
Win10 64bit 32gb
FAILURES: query: nonexistent table: foobar
Check ETA Book - OKAY'
		formattedFailureMsg = '<h3>Continuous Tests Results</h3>' $
'ETA Tests - With Test Data
Build date: Dec 19 2016 12:34:23 (vs2015 Release)
Win10 64bit 32gb
<span style="color: red">FAILURES: query: nonexistent table: foobar</span>
Check ETA Book - OKAY'

		text = formatResult(legitFailureMsg).text
		Assert(text like: formattedFailureMsg)

		warningMsg = 'SuneidoLog [#20210517.091810751]: INFO: AmazonIAMviaService.' $
			'iamCredentials - axon - "Http.POST failed: HTTP/1.1 401 Unauthorized\r\n' $
			'SuneidoLog [#20210517.091810752]: WARNING: CircuitBreaker for ' $
				'NetworkService is open, service calls will be suspended\r\n' $
			'SuneidoLog [#20210517.091957205]: Check Invalid Connections SUCCEEDED\r\n'
		formattedWarningMsg = '<h3>Continuous Tests Results</h3>' $
			'SuneidoLog [#20210517.091810751]: INFO: ' $
				'AmazonIAMviaService.iamCredentials - axon - "Http.POST failed: ' $
				'HTTP/1.1 401 Unauthorized\r\n' $
			'<span style="color: orange">' $
				'SuneidoLog [#20210517.091810752]: WARNING: CircuitBreaker for ' $
				'NetworkService is open, service calls will be suspended</span>\r\n' $
			'SuneidoLog [#20210517.091957205]: Check Invalid Connections SUCCEEDED\r\n'
		results = formatResult(warningMsg)
		Assert(results.text like: formattedWarningMsg)
		Assert(results.status is: 'WARNING')

		warningErrorMsg = 'SuneidoLog [#20210517.091810751]: INFO: AmazonIAMviaService.' $
			'iamCredentials - axon - "Http.POST failed: HTTP/1.1 401 Unauthorized\r\n' $
			'ERROR: this is an error message\r\n' $
			'SuneidoLog [#20210517.091810752]: WARNING: CircuitBreaker for ' $
				'NetworkService is open, service calls will be suspended\r\n' $
			'SuneidoLog [#20210517.091957205]: Check Invalid Connections SUCCEEDED\r\n'
		results = formatResult(warningErrorMsg)
		Assert(results.status is: 'FAILED')
		}
	cl: ContinuousTest_Base
		{
		New(.newExe?, .masterChanges, .log) { .putFileLogs = Object() }
		DisableSuneidoVariables(@unused) { }
		EnsureImportHistoryTable()
			{
			}
		ContinuousTest_Base_updateImportHistory(unused)
			{
			}
		ContinuousTest_Base_checkForLatestExe(@unused)
			{
			return .newExe?
			}
		ContinuousTest_Base_getMasterChanges()
			{
			if .masterChanges.Has?('FAILURES')
				throw .masterChanges
			return .masterChanges
			}
		ContinuousTest_Base_checkForceRunFile(forcedRun?)
			{
			return forcedRun?
			}
		ContinuousTest_Base_continuousTestType()
			{
			return 'test'
			}
		SendResults(@args)
			{
			.log.Add(args)
			}
		ContinuousTest_Base_putFile(file, text)
			{
			.putFileLogs.Add(Object(file, text))
			}
		ClearPutFileLogs()
			{
			.putFileLogs = Object()
			}
		GetPutFileLogs()
			{
			return .putFileLogs
			}
		ContinuousTest_Base_deleteFile(unused)
			{
			return ''
			}
		}
	Test_GetVersionControlChanges()
		{
		// newExe?: false; has code changes
		c1 = (.cl)(false, 'changes', logs = Object())
		c1.GetVersionControlChanges()
		Assert(logs isSize: 0)
		Assert(c1.GetPutFileLogs() is: #(("getChangesSucceeded.txt",
			"Version Control Changes:\r\nchanges")))
		c1.ClearPutFileLogs()

		c1.GetVersionControlChanges(forcedRun?:)
		Assert(logs isSize: 0)
		Assert(c1.GetPutFileLogs() is: #(("getChangesSucceeded.txt",
			"Version Control Changes:\r\nchanges")))
		c1.ClearPutFileLogs()

		// newExe?: true; has code changes
		c1 = (.cl)('New Exe', 'changes', logs = Object())
		c1.GetVersionControlChanges()
		Assert(logs isSize: 0)
		Assert(c1.GetPutFileLogs() is: #(("getChangesSucceeded.txt",
			"New Exe\r\n\r\nVersion Control Changes:\r\nchanges")))
		c1.ClearPutFileLogs()

		c1.GetVersionControlChanges(forcedRun?:)
		Assert(logs isSize: 0)
		Assert(c1.GetPutFileLogs() is: #(("getChangesSucceeded.txt",
			"New Exe\r\n\r\nVersion Control Changes:\r\nchanges")))
		c1.ClearPutFileLogs()

		// newExe? is true; no code changes
		c1 = (.cl)('New Exe', '', logs = Object())
		c1.GetVersionControlChanges()
		Assert(logs isSize: 0)
		Assert(c1.GetPutFileLogs() is: #(("getChangesSucceeded.txt", "New Exe")))
		c1.ClearPutFileLogs()

		c1.GetVersionControlChanges(forcedRun?:)
		Assert(logs isSize: 0)
		Assert(c1.GetPutFileLogs() is: #(('getChangesSucceeded.txt',
			'New Exe\r\n\r\nForcing test run - No Changes found')))
		c1.ClearPutFileLogs()

		c1 = (.cl)(false, '', logs = Object())
		c1.GetVersionControlChanges(forcedRun?:)
		Assert(logs isSize: 0)
		Assert(c1.GetPutFileLogs() is: #(('getChangesSucceeded.txt',
			'Forcing test run - No Changes found')))
		c1.ClearPutFileLogs()

		// get master change error
		c1 = (.cl)(false, 'FAILURES: Version Control Conflicts in stdlib',
			logs = Object())
		c1.GetVersionControlChanges(forcedRun?:)
		Assert(logs
			is: #(('test - Import Version Control Changes FAILED',
				'Instance: test\r\nFAILURES: Version Control Conflicts in stdlib')))
		Assert(c1.GetPutFileLogs() isSize: 0)
		}

	Test_checkForLatestExe()
		{
		mock = Mock(ContinuousTest_Base)
		mock.When.checkForLatestExe().CallThrough()
		mock.When.SendResults([anyArgs:]).Do({ })
		mock.When.copyFile([anyArgs:]).Do({ })

		// Test empty mapping / errors
		// Ensure blank strings do NOT result in SendResults
		errors = Object(Object(' ', ''), Object('   '), Object('\r\n', '\t'))
		mock.When.compareExes().Return([copyMap: #(), :errors])
		Assert(mock.checkForLatestExe() is: false)
		mock.Verify.Never().copyFile([anyArgs:])
		mock.Verify.Never().SendResults([anyArgs:])

		// Test bulk
		compareResults =  [
			copyMap: [
				version1: [
					serverExeStd:  'server.exe',
					serverExeTest: 'server_version1.exe',
					clientExeStd:  'client.exe'
					clientExeTest: 'client_version1.exe',
					stdExePath:    `C:/fake/version1`,
					message:       'Successfully copied version1 files'
					]
				version2: [
					serverExeStd:  'server.exe',
					serverExeTest: 'server_version2.exe',
					clientExeStd:  'client.exe'
					clientExeTest: 'client_version2.exe',
					stdExePath:    `C:/fake/version2`,
					message:       'Successfully copied version2 files'
					]
				]
			errors: [
				version1: Object()
				version2: Object()
				version3: Object('Failed to get file date for comparison, ' $
					'(client_version3.exe)')
				]
			]
		mock.When.compareExes().Return(compareResults)
		mock.When.copyFile(`C:/fake/version2/client.exe`, './client_version2.exe').
			Throw('system command failed with result: 1')

		Assert(mock.checkForLatestExe() is: 'Successfully copied version1 files')
		mock.Verify.copyFile(`C:/fake/version1/client.exe`, `./client_version1.exe`)
		mock.Verify.copyFile(`C:/fake/version1/server.exe`, `./server_version1.exe`)
		mock.Verify.copyFile(`C:/fake/version2/client.exe`, `./client_version2.exe`)
		mock.Verify.copyFile(`C:/fake/version2/server.exe`, `./server_version2.exe`)
		mock.Verify.SendResults([has: 'Check Latest Exes FAILED'],
			[startsWith: 'Errors Encountered:'])
		mock.Verify.SendResults([has: 'Check Latest Exes FAILED'],
			[has: 'Failed to get file date for comparison, (client_version3.exe)'])
		mock.Verify.SendResults([has: 'Check Latest Exes FAILED'],
			[has: 'Failed to copy exe: C:/fake/version2/client.exe, ' $
				'(system command failed with result: 1)'])
		Assert(compareResults.copyMap hasntMember: 'version2')
		}

	Test_compareChecksums()
		{
		cl = ContinuousTest_Base
			{
			ContinuousTest_Base_libsToCheck() { return #(stdlib, Help) }
			SetChkSumResult(result) { .ckSumResult = result }
			ContinuousTest_Base_checksumLibraries(@unused)
				{
				return .ckSumResult
				}
			}
		testCl = new cl
		fn = testCl.ContinuousTest_Base_compareChecksums
		testCl.SetChkSumResult(
			#((lib: 'stdlib', cksum: '0x00000001'), (lib: 'Help', cksum: '0x00000003')))
		checksumOb = #((lib: stdlib, cksum: '0x00000001'),
			(lib: Help, cksum: '0x00000002'))
		Assert(fn(checksumOb) is: 'Checksum mismatches:\r\n' $
			'\tHelp - svc 0x00000002; local 0x00000003')
		testCl.SetChkSumResult(#((lib: 'stdlib', cksum: '0x00000001')))
		Assert(fn(checksumOb) is: 'Checksum mismatches:\r\n' $
			'\tHelp - svc 0x00000002; local missing')

		testCl.SetChkSumResult(
			#((lib: 'stdlib', cksum: '0x00000001'), (lib: 'Help', cksum: '0x00000002')))
		checksumOb = #((lib: Help, cksum: '0x00000003'))
		Assert(fn(checksumOb) is: 'Checksum mismatches:\r\n' $
			'\tstdlib - svc    missing; local 0x00000001\r\n' $
			'\tHelp   - svc 0x00000003; local 0x00000002')

		testCl.SetChkSumResult(
			#((lib: 'stdlib', cksum: '0x00000001'), (lib: 'Help', cksum: '0x00000002')))
		checksumOb = #((lib: stdlib, cksum: '0x00000001'),
			(lib: Help, cksum: '0x00000002'))
		Assert(fn(checksumOb) is: '')
		}

	Test_collectChecksums()
		{
		cl = ContinuousTest_Base
			{
			SetGetFileResult(result) { .getFileResult = result }
			ContinuousTest_Base_getFile(unused)
				{
				return .getFileResult
				}
			ContinuousTest_Base_deleteFile(unused) { return '' }
			}

		testCl = new cl
		fn = testCl.ContinuousTest_Base_collectChecksums

		testCl.SetGetFileResult(false)
		checksumOb = fn()
		Assert(checksumOb is: #())

		testCl.SetGetFileResult('{"stdlib":"0x00000001"}')
		checksumOb = fn()
		Assert(checksumOb is: #(stdlib: '0x00000001'))

		testCl.SetGetFileResult('[{"lib":"stdlib","cksum":"0x00000001"},' $
			'{"lib":"otherlib","cksum":"0x00000003"}]')
		checksumOb = fn()
		Assert(checksumOb is:
			#((lib: stdlib, cksum: '0x00000001'), (lib: otherlib, cksum: '0x00000003')))
		}

	Teardown()
		{
		if not .exeFolderExist?
			DeleteDir('exeFolder')
		super.Teardown()
		}
	}
