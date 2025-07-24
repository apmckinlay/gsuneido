// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()

		// Do not need to test encoding of file content as it has its own test
		.MakeLibraryRecord([
			name: 'Base64BufferEncodeFile'
			text: 'function(fileRead /*=unused*/, fileWrite)
				{ fileWrite.Writeline("<encoded content>") }'])
		}

	mime(testName)
		{
		cl = new MimeMultiBuffered { ValidateFile(unused) { } }
		// Subject is a parent class method, ensure it works with this child class
		return cl.Subject('This is test: ' $ testName)
		}

	Test_simple()
		{
		cl = .mime(testName = 'Simple')

		file = FakeFile('')
		boundary = cl.MimeMultiBuffered_buildBoundary()
		cl.MimeMultiBuffered_outputParts(file, boundary)

		.assertHeader(file, testName, boundary)
		.assertFileEnd(file, boundary)
		}

	assertHeader(file, testName, boundary)
		{
		Assert(file.Readline() is: 'MIME-Version: 1.0')
		Assert(file.Readline() is: 'Subject: This is test: ' $ testName)
		Assert(file.Readline() is: 'Content-Type: multipart/mixed; ')
		Assert(file.Readline() is: '\tboundary=' $ Display(boundary[2..]))
		Assert(file.Readline() is: '')
		}

	assertFileEnd(file, boundary)
		{
		Assert(file.Readline() is: boundary $ '--')
		Assert(file.Readline() is: false)
		}

	Test_Attach()
		{
		cl = .mime(testName = 'Attach')
		cl.Attach(MimeText(content = 'Hello\r\nThis is a test\r\nNot too complex'))

		file = FakeFile('')
		boundary = cl.MimeMultiBuffered_buildBoundary()
		cl.MimeMultiBuffered_outputParts(file, boundary)

		.assertHeader(file, testName, boundary)
		Assert(file.Readline() is: boundary)

		.assertAttach(file, content)
		.assertFileEnd(file, boundary)
		}

	assertAttach(file, content)
		{
		Assert(file.Readline() is: 'Content-Type: text/plain; charset="us-ascii"')
		Assert(file.Readline() is: 'Content-Transfer-Encoding: 7bit')
		Assert(file.Readline() is: '')
		for line in content.Lines()
			Assert(file.Readline() is: line)
		Assert(file.Readline() is: '')
		}

	Test_AttachFile()
		{
		cl = .mime(testName = 'AttachFile')
		cl.AttachFile('test.pdf', attachFileName: 'renameTest.pdf')

		file = FakeFile('')
		boundary = cl.MimeMultiBuffered_buildBoundary()
		cl.MimeMultiBuffered_outputParts(file, boundary)

		.assertHeader(file, testName, boundary)
		Assert(file.Readline() is: boundary)

		.assertAttachFile(file, 'application/pdf', 'renameTest.pdf')
		.assertFileEnd(file, boundary)

		cl = .mime('testAttachFile')

		getTempName = EmailAttachment_Mime.EmailAttachment_Mime_getTempName
		cl.AttachFile(getTempName('C:\cleanUpOriginal.pdf'),
			attachFileName: 'renameCleanUpOriginal.pdf')

		att = cl.MimeMultiBuffered_attachFiles
		Assert(att isSize: 1)
		Assert(att['renameCleanUpOriginal.pdf']['cleanupOriginal?'])
		}

	assertAttachFile(file, type, name)
		{
		Assert(file.Readline() is: 'Content-Type: ' $ type $ '; ')
		Assert(file.Readline() is: '\tname=' $ Display(name))
		Assert(file.Readline() is: 'Content-Transfer-Encoding: base64')
		Assert(file.Readline() is: 'Content-Disposition: attachment; ')
		Assert(file.Readline() is: '\tfilename=' $ Display(name))
		Assert(file.Readline() is: '')
		Assert(file.Readline() is: '<encoded content>')
		Assert(file.Readline() is: '')
		}

	Test_Complex()
		{
		cl = .mime(testName = 'Complex')
		cl.Attach(MimeText(content1 = 'Hello\r\nThis is a test'))
		cl.Attach(MimeText(content2 = 'This test\r\nIs decently complex'))
		cl.AttachFile('test.pdf', attachFileName: 'renameTest.pdf')
		cl.AttachFile('test.jpeg')

		file = FakeFile('')
		boundary = cl.MimeMultiBuffered_buildBoundary()
		cl.MimeMultiBuffered_outputParts(file, boundary)
		.assertHeader(file, testName, boundary)
		Assert(file.Readline() is: boundary)

		.assertAttach(file, content1)
		Assert(file.Readline() is: boundary)

		.assertAttach(file, content2)
		Assert(file.Readline() is: boundary)

		.assertAttachFile(file, 'application/pdf', 'renameTest.pdf')
		Assert(file.Readline() is: boundary)

		.assertAttachFile(file, 'image/jpeg', 'test.jpeg')
		.assertFileEnd(file, boundary)
		}

	Test_GetMimeTextMessageContent()
		{
		lines = MimeMultiBuffered().
			Attach(MimeText('This will be line 1')).
			Attach(MimeText('This will be line 2')).
			Attach(MimeText('This will be line 3')).
			Attach(MimeText('This will be line 4')).
			Attach(MimeText('This will be line 5')).
			GetMimeTextMessageContent().Lines()

		Assert(lines isSize: 9)
		Assert(lines[0] is: 'This will be line 1')
		Assert(lines[1] is: '')
		Assert(lines[2] is: 'This will be line 2')
		Assert(lines[3] is: '')
		Assert(lines[4] is: 'This will be line 3')
		Assert(lines[5] is: '')
		Assert(lines[6] is: 'This will be line 4')
		Assert(lines[7] is: '')
		Assert(lines[8] is: 'This will be line 5')
		}
	}