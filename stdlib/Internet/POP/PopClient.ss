// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(pop_server, user, pass, timeout = 60)
		{
		.sc = SocketClient(pop_server, 110, timeout) /*= port */
		.sc.Readline()
		// authenticate
		.sc.Writeline("USER " $ user)
		if (.sc.Readline().Prefix?("-ERR"))
			{
			.Close(); throw "invalid user: " $ user
			}
		.sc.Writeline("PASS " $ pass)
		if (.sc.Readline().Prefix?("-ERR"))
			{
			.Close(); throw "invalid  password for: " $ user
			}
		}
	List(i = false)
		{
		if (.sc is false)
			return false
		if (i is false)
			{
			.sc.Writeline("LIST")
			return .getbody()
			}
		.sc.Writeline("LIST " $ i)
		line = .sc.Readline()
		if (line.Prefix?("-ERR"))
			return false
		return line
		}
	GetMessageSize(i)
		// pre: i is an integer
		// post: the size of the message is returned as a
		// number (octets), or false is returned if there is no such
		// message.
		{
		list = .List(i)
		if (list is false)
			return false

		// strip off the prefix of the server response
		list = list[String(i).Size() + 5 ..] /*= remaining prefix size */
		// do the following in case the server appends
		// additional information.
		list = list[.. list.Find(" ")]
		// this could fail if the server didn't repond to the
		// LIST command in the standard manner.
		try
			size = Number(list)
		catch
			size = false
		return size
		}
	DeleteMessage(i)
		{
		if (.sc is false)
			return false
		.sc.Writeline("DELE " $ i)
		return not .sc.Readline().Prefix?("-ERR")
		}
	Top(i, nlines = 0)
		{
		if (.sc is false)
			return false
		.sc.Writeline("TOP " $ i $ " " $ nlines)
		return .getbody()
		}
	GetMessage(i)
		{
		if (.sc is false)
			return false
		// retrieve message
		.sc.Writeline("RETR " $ i)
		return .getbody()
		}
	getbody()
		{
		if ((line = .sc.Readline()).Prefix?("-ERR"))
			return false
		s = ""
		while ((line = .sc.Readline()) isnt '.')
			s $= line $ '\r\n'
		return s
		}
	Close()
		{
		if (.sc is false)
			return
		.sc.Writeline("QUIT")
		.sc.Close()
		.sc = false
		return
		}
	}
