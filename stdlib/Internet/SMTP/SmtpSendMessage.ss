// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
function (server, from, to, subject, message, header_extra = "",
	user = false, password = false, helo_arg = "")
	{
	smtp = SmtpClient(server, user, password, helo_arg: helo_arg)
	// SmtpClient requires a "header_from" argument. For simplicity, this
	// function sends its from argument for both of SendMail's arguments.
	result = smtp.SendMail(from, from, to, subject, message, header_extra)
	smtp.Close()
	return result
	}
