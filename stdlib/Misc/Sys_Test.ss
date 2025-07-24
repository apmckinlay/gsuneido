// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_kill()
		{
		cl = Sys
			{
			Sys_killDbSession(unused){}
			Sys_toKill()
				{
				return _toKill
				}
			Sys_threadList()
				{
				return #(
					'user1@remote1<dXNlcjErMA==>(jsS)',
					'user1@remote1<dXNlcjErMQ==>(jsS)',
					'user2@remote1<dXNlcjIrMA==>(jsS)')
				}
			}

		// kill an exe client session
		_toKill = Object()
		cl.Sys_kill('user3@remote1')
		Assert(_toKill isSize: 0)
		cl.Sys_kill('user4@remote1')
		Assert(_toKill isSize: 0)

		// kill a sujs client session
		cl.Sys_kill('user1@remote1(jsS)')
		Assert(_toKill isSize: 2)
		_toKill = Object()
		cl.Sys_kill('user2@remote1(jsS)')
		Assert(_toKill isSize: 1)
		_toKill = Object()
		cl.Sys_kill('user1@remote1<dXNlcjErMA==>(jsS)')
		Assert(_toKill isSize: 2)
		}
	}