// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	cl: SvcSocketClient
		{
		CallClass() { }

		ResetAll() { }

		SvcSocketClient_openSocketClient()
			{
			mock = Mock()
			mock.When.Readline().Return(@_readline)
			return mock
			}

		SvcSocketClient_send(@unused)
			{ return _read.PopFirst() }

		SvcSocketClient_subscribe() { }

		SvcSocketClient_logError(unused) { }

		SvcSocketClient_svcSettings()
			{ return _settings }

		SvcSocketClient_setState(state)
			{ .SvcSocketClient_state = state }
		}

	Test_Main()
		{
		_settings = []
		_read = ['key', false]
		_readline = ['invalid server']
		inst = new .cl

		// Invalid Svc Server
		inst.Run({ |unused| /*do nothing*/ })
		.assertInstState(inst, err: inst.InvalidServer)

		// Valid Svc Server, invalid passowrd
		inst.Close()
		_readline = [inst.ValidServer, 'ERR invalid server']
		inst.Run({ |unused| /*do nothing*/ })
		.assertInstState(inst,
			inst.InvalidCredentials $ '\r\n\r\nLogin or Password is not set')

		// Valid Svc Server, unable to get key
		_settings = [svc_server: 'server', svc_userId: 'test', svc_passhash: 'hash']
		_read = [false]
		_readline = [inst.ValidServer]
		inst.SvcSocketClient_updateSettings()
		inst.Run({ |unused| /*do nothing*/ })
		.assertInstState(inst, err: inst.InvalidKey)

		// Valid Svc Server, retrieved key, login failed
		_settings = [svc_server: 'server', svc_userId: 'test', svc_passhash: 'hash']
		_read = ['lib:name', false]
		_readline = [inst.ValidServer]
		inst.SvcSocketClient_updateSettings()
		inst.Run({ |unused| /*do nothing*/ })
		.assertInstState(inst, err: inst.InvalidCredentials)

		// Valid Svc Server, retrieved key, login succeeds
		_settings = [svc_server: 'server', svc_userId: 'test', svc_passhash: 'hash']
		_read = ['lib:name', true]
		_readline = [inst.ValidServer]
		inst.SvcSocketClient_updateSettings()
		inst.Run({ |unused| /*do nothing*/ })
		.assertInstState(inst, free: 1, connected?:)

		// Valid Svc Server, password and table
		inst.Close()
		_read = ['lib:name', true]
		_readline = [inst.ValidServer, 'valid password']
		inst.Run({ |unused| /*do nothing*/ })
		.assertInstState(inst, free: 1, connected?:)

		// Run, will reserve the one available socket,
		// Upon completion, it will be re-added to socket pool
		inst.Run()
			{ | unused |
			.assertInstState(inst, used: 1, connected?:)
			}
		.assertInstState(inst, free: 1, connected?:)
		// Run, will reserve the one available socket, and change the settings
		// Upon completion, as the settings have changed, the socket will be closed and
		// not re-added
		inst.Run()
			{ | unused |
			.assertInstState(inst, used: 1, connected?:)
			inst.SvcSocketClient_updateSettings()
			}
		.assertInstState(inst)

		// As there are no available sockets, a new one is opened
		_read = ['lib:name', true]
		inst.Run({ |unused| /*do nothing*/ })
		.assertInstState(inst, free: 1, connected?:)
		// As another socket opened while none were available, a new one is opened
		inst.Run()
			{ | unused |
			_read = ['lib:name', true]
			.assertInstState(inst, used: 1, connected?:)
			inst.Run()
				{ |unused|
				.assertInstState(inst, used: 2, connected?:)
				}
			// Upon completion, the socket is re-added to free
			.assertInstState(inst, free: 1, used: 1, connected?:)
			// close is called,
			// - all free sockets are closed
			// - all used sockets are closed and no-added to free pool upon completion
			inst.Close()
			.assertInstState(inst, used: 1)
			}
		.assertInstState(inst)

		// -- Code Coverage --
		// Ensuring the below does not throw errors. Socket client could potentially
		// succeed / fail the below depending on the actual settings of the SvcServer
		// falls back to svcSettings
		_read = ['key', true]
		_settings = [svc_server: 'Fake', svc_userId: 'User', svc_passhash: 'Fake']
		inst.SvcSocketClient_userId = inst.SvcSocketClient_passhash = 'Fake'
		Assert(inst.TestConnect(false))

		// set default is used
		_read = ['key', true]
		inst.SvcSocketClient_userId = inst.SvcSocketClient_passhash = ''
		Assert(inst.TestConnect(false) is: false)
		}

	assertInstState(inst, err = '', free = 0, used = 0, connected? = false)
		{
		Assert(inst.Error is: err)
		Assert(inst.Connections.free isSize: free)
		Assert(inst.Connections.used isSize: used)
		Assert(inst.Connected? is: connected?)
		}

	Test_SocketClientFailsToOpen()
		{
		_settings = []
		mock = Mock(SvcSocketClient)
		mock.When.openSocketClient().CallThrough()
		mock.When.logError([anyArgs:]).Do({ })
		mock.SvcSocketClient_connections = Object()
		Assert(mock.SvcSocketClient_openSocketClient() is: false)
		Assert(mock.Error is: mock.SvcSocketClient_invalidConnection $
			':\r\n\r\nServer is not set')

		mock.Close()
		mock.SvcSocketClient_server = 'server is set'
		mock.When.SC().Throw(err = 'ERROR: failed to connect')
		Assert(mock.SvcSocketClient_openSocketClient() is: false)
		Assert(mock.Error is: mock.SvcSocketClient_invalidConnection $ ':\r\n\r\n' $ err)
		}

	Test_fatalError?()
		{
		fn = SvcSocketClient.SvcSocketClient_fatalError?

		Assert(fn('a timeout occured') is: false)
		Assert(fn('somehow lost connection') is: false)
		Assert(fn('socketclient standard error') is: false)
		Assert(fn('a non-standard socketclient error'))
		Assert(fn('a real error occured'))
		Assert(fn(SvcSocketClient.InvalidCredentials) is: false)
		Assert(fn(SvcSocketClient.InvalidServer) is: false)
		Assert(fn(SvcSocketClient.InvalidKey) is: false)
		}
	}