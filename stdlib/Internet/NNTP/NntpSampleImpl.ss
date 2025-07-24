// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
NntpImplBase
	{
	LIST()
		{
		return Object(.ListHeaderLine,
			'm.test 9 7 n' // group, high, low, status
			'.') // end
		}

	GROUP()
		{
		return Object('211 3 7 9 m.test') // 211, count, low, high, group
		}

	XOVER()
		{
		// article number, subject, user, date, id, byte count, lines
		return Object(.XoverHeaderLine
			'7\thello 7\tuser\t3 Jun 2019 04:38:40\t7\t\t\t',
			'8\thello 8\tuser\t3 Jun 2019 04:38:40\t8\t\t\t',
			'9\thello 9\tuser\t3 Jun 2019 04:38:40\t9\t\t\t',
			'.')
		}

	HEAD(args)
		{
		number = args
		return Object('221 0 ' $ number
			'From: "user'
			'Subject: hello' $ number
			'Newsgroups: m.test'
			'Date: 3 Jun 2019 04:38:40'
			'Message-ID: ' $ number,
			'.'
			)
		}

	ARTICLE(args)
		{
		number = args
		return Object('220 ' $ number $ ' ' $ number,
			'From: user',
			'Newsgroups: m.test'
			'Subject: hello ' $ number
			'Date: 3 Jun 2019 04:38:40'
			'Message-ID: ' $ number,
			'',
			'This is just a test article for id ' $ number
			'This is testing new line'
			'.'
			)
		}
	}