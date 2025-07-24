// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// abstract base class
// derived classes supply:
//		Authenticate(username, password)
//		GetMessages(username)
//		Remove(i)
//		Complete()
SocketServer
	{
	Name: "POP3 Server"
	Port: 110
	Run()
		{
		try
			{
			.user = ''
			.state = 'User'
			.Writeline("+OK Suneido POP3 server ready")
			while (.state isnt 'Closed' and false isnt request = .Readline())
				{
				command = request.Extract("^[a-zA-Z]*").Upper()
				argument = request[command.Size() + 1 ..]
				this[.state](command, argument)
				}
			}
		catch (e)
			{
			if not e.Prefix?('socket')
				SuneidoLog("ERROR in Pop3Server " $ e)
			}
		}
	// states
	User(command, argument)
		{
		switch (command)
			{
		case 'USER' :
			.user = argument
			.Writeline("+OK hi " $ .user)
			.state = 'Password'
		case 'QUIT' :
			.Writeline("+OK bye " $ .user)
			.state = 'Closed'
		default :
			.Writeline("-ERR expecting USER or QUIT")
			}
		}
	Password(command, argument)
		{
		switch (command)
			{
		case "PASS" :
			.password = argument
			if (.Authenticate(.user, .password))
				{
				.Writeline("+OK password accepted")
				.messages = .GetMessages(.user)
				.state = 'Transaction'
				}
			else
				{
				.Writeline("-ERR invalid")
				.user = ''
				.state = 'User'
				}
		case "QUIT" :
			.Writeline("+OK bye " $ .user)
			.state = 'Closed'
		default :
			.Writeline("-ERR expecting PASS or QUIT")
			.user = ''
			.state = 'User'
			}
		}
	Transaction(command, argument)
		{
		method = 'Tran_' $ command
		if .isMember(method)
			{
			this[method](:argument)
			return
			}
		.Writeline("-ERR expecting STAT, LIST, RETR, DELE, NOOP, or RSET")
		}
	isMember(method) // for test
		{
		return .Member?(method)
		}
	Tran_STAT()
		{
		size = 0
		for (msg in .messages)
			size += msg.Size()
		.Writeline("+OK " $ .messages.Size() $ " " $ size)
		}
	Tran_LIST(argument)
		{
		if (argument is "")
			{
			.Writeline("+OK " $ .messages.Size() $ " messages")
			for (i in .messages.Members())
				.Writeline((i + 1) $ " " $ .messages[i].Size())
			.Writeline(".")
			}
		else if false isnt i = .findMsg(argument)
			.Writeline("+OK " $ (i + 1) $ " " $ .messages[i].Size())
		else
			.Writeline("-ERR no such message")
		}
	Tran_UIDL(argument)
		{
		if (argument is "")
			{
			.Writeline("+OK " $ .messages.Size() $ " messages")
			for (i in .messages.Members())
				.Writeline((i + 1) $ " " $ .uid(i))
			.Writeline(".")
			}
		else if false isnt i = .findMsg(argument)
			.Writeline("+OK " $ (i + 1) $ " " $ .uid(i))
		else
			.Writeline("-ERR no such message")
		}
	uid(i)
		{
		s = .messages[i]
		n = 0
		for (i = s.Size() - 1; i >= 0; --i)
			n += s[i].Asc()
		return s.Extract("Date: (.*)").Tr(' ') $ n
		}
	Tran_TOP(argument)
		{
		args = argument.Split(' ')
		i = Number(args[0]) - 1
		n = Number(args[1])
		if (.messages.Member?(i))
			{
			.Writeline("+OK")
			// var name hiddenlines will remove messages from
			// SuneidoLogs if exception is caught
			hiddenlines = .messages[i].Split("\r\n")
			for (j = 0; hiddenlines[j] isnt ""; ++j)
				.Writeline(hiddenlines[j])
			.Writeline("")
			for (++j; j < hiddenlines.Size() and n > 0; ++j, --n)
				.Writeline(hiddenlines[j])
			.Writeline(".")
			}
		else
			.Writeline("-ERR no such message")
		}
	Tran_RETR(argument)
		{
		if false isnt i = .findMsg(argument)
			{
			.Writeline("+OK " $ .messages[i].Size() $ " octets")
			// TODO: encode periods
			.Writeline(.messages[i])
			.Writeline(".")
			}
		else
			.Writeline("-ERR no such message")
		}
	Tran_DELE(argument)
		{
		if false isnt i = .findMsg(argument)
			{
			.Remove(i)
			.Writeline("+OK")
			}
		else
			.Writeline("-ERR no such message")
		}
	findMsg(argument)
		{
		i = Number(argument) - 1
		return .messages.Member?(i) ? i : false
		}
	Tran_NOOP()
		{
		.Writeline("+OK")
		}
	Tran_RSET()
		{
		.Writeline("+OK")
		}
	Tran_QUIT()
		{
		.Writeline("+OK bye " $ .user)
		.state = 'Closed'
		.Complete()
		}
	// default methods
	Authenticate(user /*unused*/, password /*unused*/)
		{
		return true
		}
	GetMessages(user /*unused*/)
		{
		return #()
		}
	Remove(i /*unused*/)
		{ }
	Complete()
		{ }
	}
