// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Pop3Server
	{
	GetMessages(user /*unused*/)
		{ return #("header\r\n\r\none\r\n") }
	}