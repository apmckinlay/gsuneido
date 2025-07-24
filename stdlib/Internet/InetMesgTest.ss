// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		mock = Mock()
		mock.When.Readline().Return(
			"From: mr",
			" from",
			"To: miss to",
			"",
			"the",
			"body",
			""
			)

		msg = InetMesg(.msg1)
		Assert(msg.Message is: .msg1)
		Assert(msg.Body is: "the\r\nbody")
		Assert(msg.Header is: "From: mr from\r\nTo: miss to\r\n")
		Assert(msg.ReadHeader(mock) is: #("From: mr from","To: miss to"))
		Assert(msg.Field("Date") is: "")
		Assert(msg.Field("Date", "031102") is: "031102")
		Assert(msg.Field("From") is: "mr from")
		Assert(msg.Field("To") is: "miss to")

		msg = InetMesg(.msg2)
		Assert(msg.Message is: .msg2)
		Assert(msg.Address("From") is: "john@smith.com")
		Assert(msg.DisplayName("From") is: "John Smith")
		Assert(msg.Address("To") is: "sue@you.com")
		Assert(msg.DisplayName("To") is: "")
		}
	msg1:
"From: mr
 from
To: miss to

the
body"
	msg2:
"From: John Smith <john@smith.com>
To: sue@You.com

body"
	Test_100_continue()
		{
		mock = Mock()
		mock.When.Readline().Return(
			'HTTP/1.1 200 ok',
			'Blah: blah'
			'')
		Assert(InetMesg.ReadHeader(mock) is: #('HTTP/1.1 200 ok', 'Blah: blah'))

		mock = Mock()
		mock.When.Readline().Return(
			'HTTP/1.1 100 continue',
			'',
			'HTTP/1.1 100 continue',
			'',
			'HTTP/1.1 200 ok',
			'Blah: blah'
			'')
		Assert(InetMesg.ReadHeader(mock) is: #('HTTP/1.1 200 ok', 'Blah: blah'))
		mock = Mock()
		mock.When.Readline().Return(
			'HTTP/1.1 200 ok',
			'HTTP/1.1 100 continue',
			'Blah: blah'
			'')
		Assert(InetMesg.ReadHeader(mock) is:
			#('HTTP/1.1 200 ok', 'HTTP/1.1 100 continue', 'Blah: blah'))
		}
	}