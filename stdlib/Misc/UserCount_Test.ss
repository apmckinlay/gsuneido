// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		cl = UserCount
			{
			UserCount_userConnections()
				{
				return #()
				}

			UserCount_skipSessionsToKill(conns)
				{
				return _skip.Members().Each({ conns.Remove(it) })
				}
			}

		_skip = #()
		Assert(0 is: cl(#()))
		Assert(1 is: cl(#('test@wts5')))
		Assert(1 is: cl(#('test@wts5', 'test@wts5')))
		Assert(2 is: cl(#('test@wts5', 'test2@wts5')))
		Assert(1 is: cl(#('PRELOGIN_test@wts5')))
		Assert(2 is: cl(#('test@wts5', 'PRELOGIN_test@wts5', 'test2@wts5')))
		Assert(2 is: cl(#('test@wts5', 'test@wts6')))
		Assert(2 is: cl(#('test@wts5', 'test2@wts6')))

		// one jsS session stuck; test exe user count
		_skip = #('user1@remote1<dXNlcjErMA==>(jsS)':)
		Assert(0 is: cl(#()))
		Assert(1 is: cl(#('test@wts5')))
		Assert(1 is: cl(#('test@wts5', 'test@wts5')))
		Assert(2 is: cl(#('test@wts5', 'test2@wts5')))
		Assert(1 is: cl(#('PRELOGIN_test@wts5')))
		Assert(2 is: cl(#('test@wts5', 'PRELOGIN_test@wts5', 'test2@wts5')))
		Assert(2 is: cl(#('test@wts5', 'test@wts6')))
		Assert(2 is: cl(#('test@wts5', 'test2@wts6')))

		// one jsS session stuck;
		_skip = #('user1@remote1<dXNlcjErMA==>(jsS)':)
		Assert(0 is: cl(Object('user1@remote1<dXNlcjErMA==>(jsS)')))
		Assert(1 is: cl(Object('user1@remote1<dXNlcjErMA==>(jsS)',
			'user1@remote1<dXNlcjErMQ==>(jsS)')))
		Assert(2 is: cl(Object('user1@remote1<dXNlcjErMA==>(jsS)',
			'user1@remote1<dXNlcjErMQ==>(jsS)',
			'user1@remote2<dXNlcjErMA==>(jsS)')))
		Assert(2 is: cl(Object('user1@remote1<dXNlcjErMA==>(jsS)',
			'user2@remote1<dXNlcjIrMA==>(jsS)',
			'user2@remote2<dXNlcjIrMA==>(jsS)')))
		Assert(2 is: cl(Object('PRELOGIN_user1@remote1<dXNlcjErMQ==>(jsS)',
			'user2@remote1<dXNlcjIrMA==>(jsS)',
			'user2@remote2<dXNlcjIrMA==>(jsS)')))
		Assert(2 is: cl(Object('PRELOGIN_user2@remote1<dXNlcjErMQ==>(jsS)',
			'user2@remote1<dXNlcjIrMA==>(jsS)',
			'user2@remote2<dXNlcjIrMA==>(jsS)')))

		// the same session returns from being stuck
		_skip = #()
		Assert(1 is: cl(Object('user1@remote1<dXNlcjErMA==>(jsS)')))
		Assert(1 is: cl(Object('user1@remote1<dXNlcjErMA==>(jsS)',
			'user1@remote1<dXNlcjErMQ==>(jsS)')))
		Assert(2 is: cl(Object('user1@remote1<dXNlcjErMA==>(jsS)',
			'user1@remote1<dXNlcjErMQ==>(jsS)',
			'user1@remote2<dXNlcjErMA==>(jsS)')))
		Assert(2 is: cl(Object('user1@remote1<dXNlcjErMA==>(jsS)',
			'user2@remote1<dXNlcjIrMA==>(jsS)',
			'user2@remote2<dXNlcjIrMA==>(jsS)')))
		Assert(3 is: cl(Object('user1@remote1<dXNlcjErMA==>(jsS)',
			'user2@remote1<dXNlcjIrMA==>(jsS)',
			'user2@remote2<dXNlcjIrMQ==>(jsS)',
			'user2@remote3<dXNlcjIrMA==>(jsS)')))
		Assert(3 is: cl(Object('user1@remote1<dXNlcjErMA==>(jsS)',
			'user2@remote1<dXNlcjIrMA==>(jsS)',
			'user2@remote2<dXNlcjIrMQ==>(jsS)',
			'user2@remote3<dXNlcjIrMA==>(jsS)',
			'user1@remote1',
			'user1@remote2',
			'user1@remote3')))
		Assert(1 is: cl(Object('PRELOGIN_user1@remote1<dXNlcjErMA==>(jsS)',
			'user1@remote1<dXNlcjErMA==>(jsS)')))
		Assert(2 is: cl(Object('PRELOGIN_user1@remote1<dXNlcjErMA==>(jsS)',
			'user1@remote1<dXNlcjErMA==>(jsS)',
			'user2@remote1<dXNlcjIrMA==>(jsS)')))
		}
	}
