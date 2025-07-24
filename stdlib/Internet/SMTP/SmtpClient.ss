// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(smtp_server, user = false, password = false,
		timeout = 60, helo_arg = "", print = false)
		{
		if "" is port = smtp_server.AfterFirst(':')
			port = 25
		smtp_server = smtp_server.BeforeFirst(':')

		.print = print is true
		.sc = SocketClient(smtp_server, port, timeout)
		if (.sc is false)
			return
		.expect('220', "connect failed")
		.writeline(((user is false ? "HELO " : "EHLO ") $ helo_arg).Trim())
		.expect('250', "server does not support authentication")
		if user isnt false
			{
			.writeline("AUTH LOGIN")
			.expect('334', "server does not support AUTH LOGIN")
			.writeline(Base64.Encode(user))
			.expect('334', "error during login")
			.writeline(Base64.Encode(password))
			.expect("235", "login failed")
			}
		}
	expect(prefix, error)
		{
		line = s = .readline()
		while s[3] is '-'
			line $= s = .readline()
		if not line.Prefix?(prefix)
			{
			.sc.Close()
			throw error $ "\nexpected " $ prefix $ " got: " $ line
			}
		}
	SendMail(from, header_from, to, subject, message,
		header_extra = "", header_to = false) // returns true or error message
		{
		msg =
			"Date: " $ Date().Format("d MMM yyyy HH:mm:ss") $ '\r\n' $
			"From: " $ header_from $ '\r\n' $
			"To: " $ (header_to is false ? to : header_to) $ '\r\n' $
			"Subject: " $ subject $ '\r\n' $
			header_extra.Trim() $ '\r\n\r\n' $
			message.ChangeEol('\r\n')
		return .SendMsg(from, to, msg)
		}
	SendMsg(from, to, msg) // returns true or error message
		{
		if .sc is false
			return "no connection"
		// reset transaction
		.writeline("RSET")
		.readline()
		.writeline("MAIL FROM:<" $ .addr(from) $ ">")
		if (false is (s = .readline()).Upper().Has?("OK") and not s.Prefix?("250"))
			return 'MAIL FROM failed: ' $ s
		to_list = to.Split(',')
		for rcpt in to_list
			{
			.writeline("RCPT TO:<" $ .addr(rcpt) $ ">")
			if (false is (s = .readline()).Upper().Has?("OK") and not s.Prefix?("250"))
				return 'RCPT TO failed: ' $ s
			}
		.writeline("DATA")
		.readline()
		.writeline(msg)
		.writeline(".")
		if not (s = .readline()).Prefix?("250")
			return s
		return true
		}
	addr(addr)
		{ // "andrew <a@b.c>" => "a@b.c"
		if false isnt s = addr.Extract('^[^<>]*<(.*)>$')
			return s.Trim()
		return addr.Trim()
		}
	writeline(line)
		{
		if .print
			Print('>', line)
		.sc.Writeline(line)
		}
	readline()
		{
		line = .sc.Readline()
		if .print
			Print('<', line)
		return line
		}
	Close()
		{
		if .sc is false
			return
		.writeline("QUIT")
		.sc.Close()
		.sc = false
		return
		}
	}
