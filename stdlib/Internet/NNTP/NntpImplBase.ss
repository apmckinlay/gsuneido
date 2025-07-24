// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
class
	{
	MODE()
		{
		return '201    Posting prohibited'
		}

	ListHeaderLine: '215 list of newsgroups follows'
	XoverHeaderLine: '224 Overview information follows'

	QUIT(args /*unused*/, server)
		{
		server.State = 'Closed'
		return '205    Connection closing'
		}
	}