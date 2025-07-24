// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_formatAddressesForSend()
		{
		fn = BookSendEmail.FormatAddressesForSend
		Assert(fn('joeblow@comany.com') is: 'joeblow@comany.com')
		Assert(fn('<joeblow@comany.com>') is: '<joeblow@comany.com>')
		Assert(fn('<<joeblow@comany.com>') is: '<joeblow@comany.com>')
		Assert(fn('Joe Blow <joeblow@comany.com>') is: 'Joe Blow <joeblow@comany.com>')
		Assert(fn('Joe Blow (joe)<joeblow@comany.com>') is:
			 'Joe Blow {joe}<joeblow@comany.com>')
		Assert(fn('Joe Blow (joe)<joeblow@comany.com>;' $
			'Jane Doe (jane)<janedoe@comany.com>') is:
				'Joe Blow {joe}<joeblow@comany.com>,Jane Doe {jane}<janedoe@comany.com>')
		Assert(fn('Joe Blow (joe)<joeblow@comany.com>;' $
			'Jane: Doe (jane)<janedoe@comany.com>') is:
				'Joe Blow {joe}<joeblow@comany.com>,Jane Doe {jane}<janedoe@comany.com>')
		}

	Test_cleanupDisplayName()
		{
		fn = BookSendEmail.CleanupDisplayName
		Assert(fn('test@company.com') is: 'test@company.com')
		// the "&" has special meaning in String.Replace, so need to make sure it's
		// handled correctly
		Assert(fn('T&T <email@address.com>') is: 'T&T <email@address.com>')
		Assert(fn('Joe: Blow <joeblow@comany.com>') is: 'Joe Blow <joeblow@comany.com>')
		Assert(fn('Joe" Blow <joeblow@comany.com>') is: 'Joe Blow <joeblow@comany.com>')
		Assert(fn("Joe' Blow <joeblow@comany.com>") is: 'Joe Blow <joeblow@comany.com>')
		Assert(fn("Joe, Blow <joeblow@comany.com>") is: 'Joe Blow <joeblow@comany.com>')
		Assert(fn('Joe Blow; <joeblow@comany.com>') is: 'Joe Blow <joeblow@comany.com>')
		Assert(fn("Joe,, Blow <joeblow@comany.com>") is: 'Joe Blow <joeblow@comany.com>')
		Assert(fn('Joe< Blow <joeblow@comany.com>') is: 'Joe Blow <joeblow@comany.com>')
		Assert(fn('>Joe Bl>ow< <joeblow@comany.com>') is: 'Joe Blow <joeblow@comany.com>')
		Assert(fn('Joe> Blow <joeblow@comany.com>') is: 'Joe Blow <joeblow@comany.com>')
		Assert(fn('Joe<seph> Blow <joeblow@comany.com>')
			is: 'Joeseph Blow <joeblow@comany.com>')

		}

	Test_emailSES_nonSuccessfulSend()
		{
		mock = Mock(BookSendEmail)
		mock.ErrorMsg = BookSendEmail.ErrorMsg
		mock.When.BookSendEmail_logErr([anyArgs:]).Return(false)
		mock.When.InvalidEmailAddressResult?([anyArgs:]).CallThrough()
		mock.When.BookSendEmail_forwardSendMsg([anyArgs:]).Return(
			// Invalid 400 and 403 response structures
			"400",
			"403",
			"500",
			"500 Unexpected Error: 400 invalid email address format"
			"500 Unexpected Error: 403 blacklisted address"
			// Valid 400 and 403 response structures
			"400 invalid email address format",
			"400 invalid email address format: additional text",
			"403 blacklisted address",
			"403 blacklisted address: additional text",
			"200 request is blocked by web filter"
			)

		f = BookSendEmail.BookSendEmail_emailSES
		mime = FakeObject(ToString: "")
		Assert(mock.Eval(f, '', '', mime, '') is: 'unable to send email'
			msg: "Response: 400")
		Assert(mock.Eval(f, '', '', mime, '') is: 'unable to send email',
			msg: "Response: 403")
		Assert(mock.Eval(f, '', '', mime, '') is: 'unable to send email'
			msg: "Response: 500")
		Assert(mock.Eval(f, '', '', mime, '') is: 'unable to send email'
			msg: "Response: Nested 400")
		Assert(mock.Eval(f, '', '', mime, '') is: 'unable to send email'
			msg: "Response: Nested 403")

		Assert(mock.Eval(f, '', '', mime, '') is: '400 invalid email address format')
		Assert(mock.Eval(f, '', '', mime, '') is: '400 invalid email address format: ' $
			'additional text')

		Assert(mock.Eval(f, '', '', mime, '') is: '403 blacklisted address')
		Assert(mock.Eval(f, '', '', mime, '') is: '403 blacklisted address: ' $
			'additional text')

		Assert(mock.Eval(f, '', '', mime, '') is: '200 request is blocked by web filter')
		}

	Test_fileMissingError()
		{
		cl = BookSendEmail
			{
			BookSendEmail_fileExists?(unused) { return false }
			}
		func = cl.BookSendEmail_fileMissingError
		watch = .WatchTable('suneidolog')

		// fileMissingError gets called from a failed try/catch, should never be ''
		func('')
		calls = .GetWatchTable(watch)
		Assert(calls[0].sulog_message
			is: 'ERROR: Failed to attach folder does not exist: ')

		errorString = `can't open 'C:/thisFolderDoesNotExist/Invoice1.pdf'  in mode 'r'`
		func(errorString)
		calls = .GetWatchTable(watch)
		Assert(calls[1].sulog_message is: 'ERROR: Failed to attach folder does ' $
			`not exist: C:/thisFolderDoesNotExist/`)

		dirSpy = .SpyOn(CheckDirExists)
		dirSpy.Return(true)

		errorString = `can't open '\\sharefolder\eta\Invoice1.pdf' in mode 'r'`
		func(errorString)
		calls = .GetWatchTable(watch)
		Assert(calls[2].sulog_message
			is: 'ERROR: Failed to attach Invoice1.pdf, it does not exist')

		cl = BookSendEmail
			{ BookSendEmail_fileExists?(unused) { return true } }
		func = cl.BookSendEmail_fileMissingError
		errorString = `can't open '\\sharefolder\eta\Invoice1.pdf' in mode 'r'`
		cl.BookSendEmail_fileMissingError(errorString)
		calls = .GetWatchTable(watch)
		Assert(calls[3].sulog_message
			is: 'ERROR: Failed to attach Invoice1.pdf, it exists, ' $
				'but was being held open and cannot attach')

		dirSpy.Close()
		}
	}
