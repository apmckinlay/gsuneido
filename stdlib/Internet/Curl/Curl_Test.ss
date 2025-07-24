// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	curlClass: Curl
		{
		Curl_runPiped(cmd)
			{
			Suneido.Curl_Test_cmd = cmd
			return ''
			}
		Curl_runPipedWithBlock(cmd, block, _mockPipe = false)
			{
			Suneido.Curl_Test_cmd_block = cmd
			if mockPipe isnt false
				block(mockPipe)
			return ''
			}
		Curl_app()
			{
			return 'curl'
			}
		}
	pre: 'curl -Y 1 -y 480 --connect-timeout 60 '
	suffix: ' -s -S'

	// HTTP --------------------------------------------------------------------

	Test_http()
		{
		_mockPipe = MockObject(#())
		Assert({ .http('XXX', 'http://abc.com/stuff') }
			throws: 'unhandled switch value')
		Assert({ .http('XXX', 'http://abc.com/stuff', content: '...', fromFile: 'dest') }
			throws: 'should not have content AND fromFile')
		}
	Test_http_get()
		{
		Assert(.http('GET', 'http://abc.com/stuff')
			is: '"http://abc.com/stuff" -D -')
		Assert(.http('GET', 'http://abc.com/stuff', toFile: 'dest')
			is: '"http://abc.com/stuff" -o "dest" -D -')
		Assert({ .http('GET', 'http://abc.com/stuff', content: '...') }
			throws: 'GET should not have body')
		Assert({ .http('GET', 'http://abc.com/stuff', fromFile: '...') }
			throws: 'GET should not have body')
		Assert(.http('GET', 'http://abc.com/stuff', header: #('x-hdr': foo))
			is: '"http://abc.com/stuff" -D - -H "x-hdr: foo"')
		}

	Test_http_check()
		{
		.httpCheck(#('GET', 'http://abc.com/stuff'),
			MockObject(#(
				(CloseWrite),
				((Readline), result: 'HTTP/1.1 200'),
				((Readline), result: ''),
				((Read), result: 'Response'),
				((Read), result: false))),
			#(header: 'HTTP/1.1 200', content: 'Response'),
			'"http://abc.com/stuff" -D -')
		.httpCheck(#('GET', 'http://abc.com/stuff', toFile: 'dest'),
			MockObject(#(
				(CloseWrite)
				((Readline), result: 'HTTP/1.1 100'),
				((Readline), result: 'content-type: text/plain;'),
				((Readline), result: ''),
				((Readline), result: 'HTTP/1.1 200'),
				((Readline), result: 'content-type: text/plain;'),
				((Readline), result: ''))),
			#(header: 'HTTP/1.1 200\ncontent-type: text/plain;', content: ''),
			'"http://abc.com/stuff" -o "dest" -D -')
		.httpCheck(#('PUT', 'http://abc.com/stuff', fromFile: 'dest'),
			MockObject(#(
				(CloseWrite),
				((Readline), result: 'HTTP/1.1 200'),
				((Readline), result: ''),
				((Read), result: 'success'),
				((Read), result: false))),
			#(header: 'HTTP/1.1 200', content: 'success'),
			'"http://abc.com/stuff" -D - -g -T "dest"')
		.httpCheck(#('POST', 'http://abc.com/stuff', content: '...'),
			MockObject(#(
				(Write, '...'),
				(CloseWrite),
				((Readline), result: 'HTTP/1.1 200'),
				((Readline), result: ''),
				((Read), result: 'success'),
				((Read), result: false))),
			#(header: 'HTTP/1.1 200', content: 'success'),
			'"http://abc.com/stuff" -D - --data-binary "@-" -H Content-Type:')
		}

	Test_http_put()
		{
		Assert(.http('PUT', 'http://abc.com/stuff', fromFile: 'dest')
			is: '"http://abc.com/stuff" -D - -g -T "dest"')
		Assert(.http('EMPTYPUT', 'http://abc.com/stuff')
			is: '"http://abc.com/stuff" -D - -X PUT')
		Assert(.http('PUT', 'http://abc.com/stuff', content: '...')
			is: '"http://abc.com/stuff" -D - ' $
				'-g -T - -H "Transfer-Encoding: " -H "Content-Length: 3"')
		}
	Test_http_post()
		{
		Assert(.http('POST', 'http://abc.com/stuff', content: '...')
			is: '"http://abc.com/stuff" -D - ' $
				'--data-binary "@-" -H Content-Type:')
		Assert(.http('POST', 'http://abc.com/stuff', fromFile: 'dest')
			is: '"http://abc.com/stuff" -D - ' $
				'--data-binary "@dest" -H Content-Type:')
		}
	Test_http_patch()
		{
		Assert(.http('PATCH', 'http://abc.com/stuff', content: '...')
			is: '"http://abc.com/stuff" -D - ' $
				'--data-binary "@-" -H Content-Type: -X PATCH')
		Assert(.http('PATCH', 'http://abc.com/stuff', fromFile: 'dest')
			is: '"http://abc.com/stuff" -D - ' $
				'--data-binary "@dest" -H Content-Type: -X PATCH')
		}
	http(@args)
		{
		(.curlClass).Http(@args)
		return Suneido.Curl_Test_cmd_block.RemovePrefix(.pre).RemoveSuffix(.suffix)
		}

	httpCheck(args, mockPipe, expectedResult, expectedCmd)
		{
		_mockPipe = mockPipe
		result = (.curlClass).Http(@args)
		Assert(result is: expectedResult)
		Assert(Suneido.Curl_Test_cmd_block.RemovePrefix(.pre).RemoveSuffix(.suffix)
			is: expectedCmd)
		}

	Test_postfiles()
		{
		cl = new .curlClass('https', 'openinvoice.test.com', 'user', 'password',
			options: #(cacert: 'cacert', key: 'key.pem', cert: 'client.pem:pass',
				files: #('receipt=@file1.txt; type=application/json',
					'attach.pdf=@attach.pdf')))
		cl.Http( 'POSTFILES', 'https://openinvoice.test.com',
			header: #(Content_Type: 'Application/json; charset=utf-8'))
		str = Suneido.Curl_Test_cmd_block.RemovePrefix(.pre).RemoveSuffix(.suffix)

		ob = #('-u "user:password"',
			'--cert client.pem:pass', '--key key.pem', '--cacert cacert',
			'-F "receipt=@file1.txt; type=application/json"',
			'-F "attach.pdf=@attach.pdf"',
			'"https://openinvoice.test.com"', '-D -', '-k',
			'-H "Content-Type: Application/json; charset=utf-8"')
		for segment in ob
			Assert(str has: segment)
		}

	// FTP ---------------------------------------------------------------------

	DirText:
"drwxr-xr-x    2 99       99           4096 01-26-05  11:53AM vti_cnf
-rwxr--r--     1 99       99          65866 01-26-05  11:53AM ajith030602.zip
-rwxr--r--     1 99       99       12649430 02-08-05  05:27PM 20050201.zip
-rwxr--r--     1 99       99         123456 02-08-05  05:27PM 20050125.zip
-rwxr--r--     1 99       99        9398287 12-21-04  02:34PM 20050115.zip"

	Test_dirList()
		{
		.testOneMatch()
		.testManyMatch()
		}

	testOneMatch()
		{
		fmask = FtpClient.BuildFMask('vti_cnf', caseSense: true)
		files = Curl.Curl_dirList(.DirText, fmask, '', details:)
		Assert(files isSize: 1)
		Assert(files has: #(name: 'vti_cnf', size: 4096))

		fmask = FtpClient.BuildFMask('Vti_Cnf', caseSense: false)
		files = Curl.Curl_dirList(.DirText, fmask, '', details:)
		Assert(files isSize: 1)
		Assert(files has: #(name: 'vti_cnf', size: 4096))

		fmask = FtpClient.BuildFMask('Vti_Cnf', caseSense: true)
		files = Curl.Curl_dirList(.DirText, fmask, '', details:)
		Assert(files isSize: 0)
		}

	testManyMatch()
		{
		fmask = FtpClient.BuildFMask('*.*', caseSense: false)
		// list all files
		files = Curl.Curl_dirList(.DirText, fmask, '', details:)
		Assert(files isSize: 5)
		Assert(files has: #(name: '20050201.zip', size: 12649430))
		Assert(files has: #(name: '20050125.zip' size: 123456))
		Assert(files has: #(name: '20050115.zip', size: 9398287))
		Assert(files has: #(name: 'ajith030602.zip', size: 65866))
		Assert(files has: #(name: 'vti_cnf', size: 4096))

		fmask = FtpClient.BuildFMask('*.zip', caseSense: false)
		// list all zip files
		files = Curl.Curl_dirList(.DirText, fmask, '', details:)
		Assert(files isSize: 4)
		Assert(files has: #(name: '20050201.zip', size: 12649430))
		Assert(files has: #(name: '20050125.zip' size: 123456))
		Assert(files has: #(name: '20050115.zip', size: 9398287))
		Assert(files has: #(name: 'ajith030602.zip', size: 65866))
		Assert(files hasnt: #(name: 'vti_cnf', size: 4096))

		// list zip files that have yyyyMMdd.zip as the file names
		files = Curl.Curl_dirList(.DirText, fmask, '^20\d\d\d\d\d\d.zip$', details:)
		Assert(files isSize: 3)
		Assert(files has: #(name: '20050201.zip', size: 12649430))
		Assert(files has: #(name: '20050125.zip' size: 123456))
		Assert(files has: #(name: '20050115.zip', size: 9398287))
		Assert(files hasnt: #(name: 'ajith030602.zip', size: 65866))
		Assert(files hasnt: #(name: 'vti_cnf', size: 4096))
		}

	Test_Dir()
		{
		curlClass = Curl
			{
			Curl_protocol: 'ftp'
			Curl_server: '127.0.0.1/'
			Curl_timeout: 60
			Curl_runCommand(cmd)
				{
				if cmd.Has?('127.0.0.1/ignoredir')
					return 'curl: (19) RETR response: 226'
				if cmd.Has?('127.0.0.1/baddir')
					return 'curl: (19) RETR response: 530'
				if cmd.Has?('127.0.0.1/emptydir')
					return ''
				if cmd.Has?('127.0.0.1/fulldir')
					return 'file1.txt\r\nfile2.txt\r\nfile3.txt'
				return 'unexpected call'
				}
			}

		Assert(curlClass.Dir('ignoredir/') is: #())
		Assert(curlClass.Dir('baddir/') is: false)
		Assert(curlClass.Dir('emptydir/') is: #())
		Assert(curlClass.Dir('fulldir/') is: #('file1.txt', 'file2.txt', 'file3.txt'))
		}

	Test_multiple_files()
		{
		file_list = #('file3.I07', 'test_folder/file (4).I08')
		curl = Curl('ftp', 'ftpserver')

		get = curl.Curl_buildGetScripts(file_list, '\\testpath\\', 'usersfolder')
		Assert(get
			like: '-o "/testpath/file3.I07"
				url = "ftp://ftpserver/usersfolder/file3.I07"
				-o "/testpath/file (4).I08"
				url = "ftp://ftpserver/usersfolder/test_folder/file%20%284%29.I08"')

		delete = curl.Curl_buildDeleteScripts(file_list, 'usersfolder', false)
		Assert(delete
			like: '-Q "CWD /usersfolder"
				-Q "DELE file3.I07"
				-Q "DELE test_folder/file (4).I08"')

		delete = curl.Curl_buildDeleteScripts(file_list, 'usersfolder', true)
		Assert(delete
			like: '-Q "CWD usersfolder"
				-Q "DELE file3.I07"
				-Q "DELE test_folder/file (4).I08"')

		curl = Curl('sftp', 'ftpserver')
		delete = curl.Curl_buildDeleteScripts(file_list, 'usersfolder', true)
		Assert(delete
			like: '-Q "rm /usersfolder/file3.I07"
				-Q "rm /usersfolder/test_folder/file (4).I08"')

		}

	Test_FtpGet()
		{
		// regular ftp
		ftp = (.curlClass)('ftp', '999.999.999.999', 'uSer', 'pAssw0rd')
		ftp.Get('abc_test.txt', 'abc_test_renamed.txt')
		Assert(Suneido.Curl_Test_cmd
			is: .pre $ '-u "uSer:pAssw0rd" -o "abc_test_renamed.txt" ' $
				'"ftp://999.999.999.999/abc_test.txt"' $ .suffix)
		}

	Test_FtpPut()
		{
		// regular ftp put
		ftp = (.curlClass)('ftp', '999.999.999.999', 'uSer', 'pAssw0rd')
		ftp.Put('abc_test.txt', 'abc_test_renamed.txt')
		Assert(Suneido.Curl_Test_cmd
			is: .pre $ '-u "uSer:pAssw0rd" -T "abc_test.txt" ' $
				'"ftp://999.999.999.999/abc_test_renamed.txt"' $ .suffix)
		}

	Test_Del()
		{
		ftp = (.curlClass)('ftp', '999.999.999.999', 'uSer', 'pAssw0rd')
		ftp.Del('abc_test.txt')
		Assert(Suneido.Curl_Test_cmd
			is: .pre $ '-u "uSer:pAssw0rd"  ' $
				'ftp://999.999.999.999/ -Q "DELE abc_test.txt"' $ .suffix)

		ftp = (.curlClass)('ftp', '999.999.999.999', 'uSer', 'pAssw0rd')
		ftp.Del('abc_test.txt', 'subdir')
		Assert(Suneido.Curl_Test_cmd
			is: .pre $ '-u "uSer:pAssw0rd"  ftp://999.999.999.999/ ' $
				'-Q "CWD /subdir" -Q "DELE abc_test.txt"' $ .suffix)

		ftp = (.curlClass)('ftp', '999.999.999.999', 'uSer', 'pAssw0rd')
		ftp.Del('abc_test.txt', 'subdir', true)
		Assert(Suneido.Curl_Test_cmd
			is: .pre $ '-u "uSer:pAssw0rd"  ftp://999.999.999.999/ ' $
				'-Q "CWD subdir" -Q "DELE abc_test.txt"' $ .suffix)

		ftp = (.curlClass)('sftp', '999.999.999.999', 'uSer', 'pAssw0rd')
		ftp.Del('abc_test.txt')
		Assert(Suneido.Curl_Test_cmd
			is: .pre $ '-u "uSer:pAssw0rd"  ' $
				'sftp://999.999.999.999/ -Q "rm /abc_test.txt"' $ .suffix)

		ftp = (.curlClass)('sftp', '999.999.999.999', 'uSer', 'pAssw0rd')
		ftp.Del('abc_test.txt', 'subdir')
		Assert(Suneido.Curl_Test_cmd
			is: .pre $ '-u "uSer:pAssw0rd"  ' $
				'sftp://999.999.999.999/ -Q "rm /subdir/abc_test.txt"' $ .suffix)
		}

	Test_buildDirectoryScripts()
		{
		curl = Curl('ftp', 'ftpserver')
		directories = Object('dirA', 'dirB', 'dirC')
		result = curl.Curl_buildDirectoryScripts(directories, 'testingCurlScripts')
		Assert(result
			is: '-l\n' $
				'url = "ftp://ftpserver/dirA"\noutput = "testingCurlScripts0"\n' $
				'url = "ftp://ftpserver/dirB"\noutput = "testingCurlScripts1"\n' $
				'url = "ftp://ftpserver/dirC"\noutput = "testingCurlScripts2"\n')
		}
	Test_cookies()
		{
		Assert(.http('POST', 'http://abc.com/stuff', content: '...',
			cookies: 'cookies.txt')
			is: '"http://abc.com/stuff" -D - ' $
				'-b cookies.txt -c cookies.txt ' $
				'--data-binary "@-" -H Content-Type:')
		}

	Test_nocurlexe()
		{
		.SpyOn(ExternalApp).Return(false)
		.SpyOn(Sys.Linux?).Return(false)
		result = Curl.Curl_runCommand('')
		Assert(result is: 'missing curl.exe')
		}

	Test_runCommand_retry()
		{
		mock = .runCommandMock()
		mock.When.runPiped([anyArgs:]).Return('') // Successfully responses are ''
		mock.runCommand('')
		mock.Verify.runPiped([anyArgs:])
		mock.Verify.extractDebugging('')

		expectedError = `curl: (35) OpenSSL SSL_connect: SSL_ERROR_SYSCALL in ` $
			`connection to s3.us-east-1.amazonaws.com:443`
		.testSslErrorRetries(mock, expectedError)
		expectedError = `curl: (35) OpenSSL SSL_connect: ` $
			`Connection was reset in connection to s3.us-east-1.amazonaws.com:443 `
		.testSslErrorRetries(mock, expectedError)
		}

	testSslErrorRetries(mock, expectedError)
		{
		mock = .runCommandMock()
		mock.When.runPiped([anyArgs:]).Return(expectedError, expectedError, '')
		mock.runCommand('')
		mock.Verify.Times(3).runPiped([anyArgs:])
		mock.Verify.extractDebugging('')

		mock = .runCommandMock()
		mock.When.runPiped([anyArgs:]).Return(expectedError)
		mock.runCommand('')
		mock.Verify.Times(4).runPiped([anyArgs:])
		mock.Verify.extractDebugging(expectedError)

		// should always retry 4 times on the expected error
		mock = .runCommandMock()
		mock.When.runPiped([anyArgs:]).Return(expectedError)
		mock.runCommand('', retries: 0)
		mock.Verify.Times(4).runPiped([anyArgs:])

		mock = .runCommandMock()
		mock.When.runPiped([anyArgs:]).Return(expectedError)
		mock.runCommand('', retries: 4)
		mock.Verify.Times(4).runPiped([anyArgs:])

		mock = .runCommandMock()
		mock.When.runPiped([anyArgs:]).Return(expectedError)
		mock.runCommand('', retries: 5)
		mock.Verify.Times(5).runPiped([anyArgs:])

		mock = .runCommandMock()
		mock.When.runPiped([anyArgs:]).Return('curl: other error')
		mock.runCommand('', retries: 5)
		mock.Verify.Times(5).runPiped([anyArgs:])

		mock = .runCommandMock()
		mock.When.runPiped([anyArgs:]).Return(expectedError, expectedError,
			'curl: other error', 'curl: other error', 'curl: other error',
			'curl: other error')
		mock.runCommand('', retries: 6)
		mock.Verify.Times(6).runPiped([anyArgs:])

		mock = .runCommandMock()
		mock.When.runPiped([anyArgs:]).Return(expectedError, 'curl: other error', '')
		mock.runCommand('', retries: 6)
		mock.Verify.Times(3).runPiped([anyArgs:])
		}

	Test_runCommand_retries_Connection_reset()
		{
		expectedError = `curl: (56) OpenSSL SSL_read: Connection was reset`
		mock = .runCommandMock()

		mock.When.runPiped([anyArgs:]).Return(expectedError, expectedError, '')
		mock.runCommand('')
		mock.Verify.Times(3).runPiped([anyArgs:])
		mock.Verify.extractDebugging('')

		mock = .runCommandMock()
		mock.When.runPiped([anyArgs:]).Return(expectedError)
		mock.runCommand('')
		mock.Verify.Times(4).runPiped([anyArgs:])
		mock.Verify.extractDebugging(expectedError)

		// should always retry 4 times on the expected error
		mock = .runCommandMock()
		mock.When.runPiped([anyArgs:]).Return(expectedError)
		mock.runCommand('', retries: 0)
		mock.Verify.Times(4).runPiped([anyArgs:])

		mock = .runCommandMock()
		mock.When.runPiped([anyArgs:]).Return(expectedError)
		mock.runCommand('', retries: 4)
		mock.Verify.Times(4).runPiped([anyArgs:])

		mock = .runCommandMock()
		mock.When.runPiped([anyArgs:]).Return(expectedError)
		mock.runCommand('', retries: 5)
		mock.Verify.Times(5).runPiped([anyArgs:])
		}

	runCommandMock()
		{
		mock = Mock(Curl)
		mock.When.runCommand([anyArgs:]).CallThrough()
		mock.When.app().Return('curl')
		return mock
		}

	Test_initUserPass()
		{
		initUsrPass = Curl.Curl_initUserPass
		// normal user/pass
		up = initUsrPass('user','pass')
		Assert(up is: ' -u "user:pass"')

		// no user/pass
		up = initUsrPass('','')
		Assert(up is: '')

		// user but no pass
		up = initUsrPass('user','')
		Assert(up is: ' -u "user"')

		// no user but pass
		up = initUsrPass('','pass')
		Assert(up is: '')

		// special char in pwd
		up = initUsrPass('user','Fred^2(8@sZq')
		Assert(up is: ' -u "user:Fred^2(8@sZq"')
		// pwd with QUOTE in them
		up = initUsrPass('user','Fred"2(8@sZq')
		Assert(up is: ' -u "user:Fred\\"2(8@sZq"')

		up = initUsrPass("user","Fred'2(8@sZq")
		Assert(up is: ' -u "user:Fred\'2(8@sZq"')
		}

	Test_version_checking()
		{
		// valid version
		cl = Curl
			{
			Curl_runPiped(@unused) { return 'curl 8.8.0 followed by the extra stuff' }
			Curl_minVersion() { return #(8, 8, 0) }
			}
		Assert(cl.Curl_version('unused') is Object(8, 8, 0))
		Suneido.Curl_VersionChecked = false
		cl.Curl_checkVersion('unused')

		// less than min version
		cl = Curl
			{
			Curl_runPiped(@unused) { return 'curl 7.68.0 followed by the extra stuff' }
			Curl_minVersion() { return #(8, 8, 0) }
			}
		Assert(cl.Curl_version('unused') is Object(7, 68, 0))
		Suneido.Curl_VersionChecked = false
		ServerSuneido.Set('TestRunningExpectedErrors', Object('ERROR: (CAUGHT) ' $
			'The curl (7.68.0) is lower than the minimum requirement (8.8.0)'))
		cl.Curl_checkVersion('unused')
		Assert(ServerSuneido.Get('TestRunningExpectedErrors') is: #())

		// invalid DEV version
		cl = Curl
			{
			Curl_runPiped(@unused) { return 'curl 8.8.0-DEV followed by the extra stuff' }
			Curl_minVersion() { return #(8, 8, 0) }
			}
		Assert(cl.Curl_version('unused') is '8.8.0-DEV')
		Suneido.Curl_VersionChecked = false
		ServerSuneido.Set('TestRunningExpectedErrors',
			Object('ERROR: (CAUGHT) Unexpected curl version: 8.8.0-DEV'))
		cl.Curl_checkVersion('unused')
		Assert(ServerSuneido.Get('TestRunningExpectedErrors') is: #())
		}

	Teardown()
		{
		super.Teardown()
		Suneido.Delete(#Curl_Test_cmd)
		}
	}
