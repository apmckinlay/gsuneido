// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		m = MimeText("hello world")
		m.To('apmckinlay@gmail.com')
		m.From('default@suneido.com')
		Assert(m.ToString() is:
			'MIME-Version: 1.0\r\n' $
			"From: default@suneido.com\r\n" $
			"To: apmckinlay@gmail.com\r\n" $
			'Content-Type: text/plain; charset="us-ascii"\r\n' $
			'Content-Transfer-Encoding: 7bit\r\n' $
			"\r\n" $
			"hello world\r\n")
		}

	Test_NewlinesInMIMEParts()
		{
		m = MimeText('hello')
		m.Subject('testing\r\nwith newline')
		Assert(m.ToString() is:
			"MIME-Version: 1.0\r\n" $
			"Subject: testing\r\n\twith newline\r\n" $
			'Content-Type: text/plain; charset="us-ascii"\r\n' $
			'Content-Transfer-Encoding: 7bit\r\n\r\n' $
			'hello\r\n')
		}
	}
