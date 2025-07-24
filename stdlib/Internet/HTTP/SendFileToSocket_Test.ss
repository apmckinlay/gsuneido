// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.testThrow(Object(), '', 'does not exist')

		sockfile = .MakeFile()
		tempfile = .MakeFile()
		fileString = ''
		PutFile(tempfile, fileString)
		.test(tempfile, Object(),
			'Date: <date>\r\nContent-Length: ' $ fileString.Size() $ '\r\n\r\n',
			sockfile)

		fileString = 'File Content: buffered send\r\nis based on file lines'
		PutFile(tempfile, fileString)
		.test(tempfile, Object('Content-Type': 'text/plain')
			'Date: <date>\r\n' $
			'Content-Type: text/plain\r\n' $
			'Content-Length: ' $ fileString.Size() $ '\r\n\r\n' $
			fileString,
			sockfile)
		}

	test(tempFile, header, result, sockfile)
		{
		File(sockfile, "w")
			{|f|
			SendFileToSocket(f, tempFile, header, 'HTTP/1.0', delete?:)
			}
		log = GetFile(sockfile)
		log = log.Replace('^Date: .*', 'Date: <date>')
		Assert(log is: 'HTTP/1.0\r\n' $ result)
		}

	testThrow(header, file, expectedThrow)
		{
		Assert({SendFileToSocket(false, file, header, 'HTTP/1.0', delete?:)}
			throws: expectedThrow)
		}
	}
