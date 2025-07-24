// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_WithoutSpecial()
		{
		Assert(UserConnections.WithoutSpecial(#()) is: #())
		Assert(UserConnections.WithoutSpecial(#('(scheduler)', '(smtp)')) is: #())
		Assert(UserConnections.WithoutSpecial(#('ab@wts5', 'cd@wts5'))
			is: #('ab@wts5', 'cd@wts5'))
		Assert(UserConnections.WithoutSpecial(#('ab@wts5', '127.0.0.1', 'cd@wts5'))
			is: #('ab@wts5', 'cd@wts5'))
		Assert(#('ab@wts5', 'cd@wts5')
			is: UserConnections.WithoutSpecial(#('ab@wts5', '(xy)', 'cd@wts5')))
		Assert(UserConnections.WithoutSpecial(#('joe@127.0.0.1')) is: #('joe@127.0.0.1'))
		Assert(UserConnections.WithoutSpecial(#('ab@255.255.255.1', '127.0.0.1:main',
			'127.0.0.1:pool-2-thread-1', 'cd@wts5', 'ab@255.255.255.1:Thread-99'))
			is: #('ab@255.255.255.1', 'cd@wts5'))
		Assert(UserConnections.WithoutSpecial(#('ab@wts5', '192.168.1.129', 'cd@wts5'))
			is: #('ab@wts5', 'cd@wts5'))
		Assert(UserConnections.WithoutSpecial(
			#('ab@wts5', 'joe@192.168.1.129', '192.168.1.129', 'cd@wts5',
				'ab@wts5:Thread-1', 'excluded@wts5:Thread-27'))
			is: #('ab@wts5', 'joe@192.168.1.129', 'cd@wts5'))
		}
	}
