// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// abstract base class
// derived classes supply:
//		ValidateRecipient(rcpt)
//		Send(from, to, msg)
SocketServer
	{
	Name: "SMTP Server"
	Port: 25
	Run()
		{
		.state = 'Helo'
		.Writeline("220 Suneido SMTP server ready")
		while (.state isnt 'Closed' and false isnt request = .Readline())
			{
			command = request.Extract("^[a-zA-Z]*").Upper()
			this[.state](command, request)
			}
		}
	// states
	Helo(command, request /*unused*/)
		{
		switch (command)
			{
		case 'HELO' :
			.Writeline("250 Suneido")
			.state = 'Mail'
		case 'NOOP' :
			.Writeline("250 OK")
		case 'QUIT' :
			.Writeline("221 Suneido goodbye")
			.state = 'Closed'
		default :
			.Writeline("500 expecting HELO")
			}
		}
	Mail(command, request)
		{
		switch (command)
			{
		case 'MAIL' :
			.from = request.AfterFirst('<').BeforeLast('>') // strip "MAIL FROM: <...>"
			.to = Object()
			.msg = Object()
			.Writeline("250 OK")
			.state = 'Rcpt'
		case 'RSET' :
			.from = ""
			.msg = Object()
			.to = Object()
			.Writeline("250 OK")
		case 'NOOP' :
			.Writeline("250 OK")
		case 'QUIT' :
			.Writeline("221 Suneido goodbye")
			.state = 'Closed'
		default :
			.Writeline("500 expecting MAIL")
			}
		}
	Rcpt(command, request)
		{
		switch (command)
			{
		case "RCPT" :
			.Handle_RCPT(request)
		case "DATA" :
			.Writeline("354 Start mail input, end with .")
			.state = 'Data'
		case 'RSET' :
			.from = ""
			.msg = Object()
			.to = Object()
			.Writeline("250 OK")
			.state = 'Mail'
		case 'NOOP' :
			.Writeline("250 OK")
		case 'QUIT' :
			.Writeline("221 Suneido goodbye")
			.state = 'Closed'
		default :
			.Writeline("500 expecting RCPT")
			}
		}
	Handle_RCPT(request)
		{
		rcpt = request.AfterFirst('<').BeforeLast('>') // strip "RCPT TO: <...>"
		if '' is msg = .ValidateRecipient(rcpt)
			{
			.AddToRcpt(rcpt)
			.Writeline("250 OK")
			}
		else
			{
			.Writeline("550 " $ msg)
			.state = 'Closed'
			}
		}

	AddToRcpt(rcpt)
		{
		.to.Add(rcpt)
		}

	Data(command /*unused*/, request)
		{
		if (request isnt ".")
			{
			request = request.RemovePrefix('.')
			if .msg.Size() < 1000 /*= max lines */
				.msg.Add(request)
			}
		else
			{
			.Writeline(result = .Send(.from, .to, .msg))
			.state = Numberable?(result[0]) and Number(result[0]) < 4/*= non-error codes*/
				? 'Mail'
				: 'Closed'
			}
		}
	// default methods
	ValidateRecipient(rcpt /*unused*/)
		{
		return ''
		}
	Send(from /*unused*/, to /*unused*/, msg /*unused*/)
		{
		return "250 OK"
		}
	}
