// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
// like SmtpSendMessage but calls SendMsg instead of SendMail
// returns true or error message
function (server, from, to, message,
	user = false, password = false, helo_arg = "", print = false)
	{
	smtp = SmtpClient(server, user, password, helo_arg: helo_arg, print: print)
	result = smtp.SendMsg(from, to, message)
	smtp.Close()
	return result
	}
