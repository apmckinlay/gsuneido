// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		cl = Mock(SvcSettingsIconsControl)
		cl.SvcSettingsIconsControl_serverIcon = serverIcon = Mock()
		cl.SvcSettingsIconsControl_passwordIcon = passwordIcon = Mock()
		cl.When.set([anyArgs:]).CallThrough()
		cl.When.svcSocketClient([anyArgs:]).Do({ })
		cl.When.connected?().Return(false)

		// No settings
		cl.SvcSettingsIconsControl_settings = []
		cl.set()
		insufficient = cl.SvcSettingsIconsControl_insufficientSettings
		serverIcon.Verify.ToolTip('Server: ' $ insufficient)
		passwordIcon.Verify.ToolTip('Password: ' $ insufficient)
		passwordIcon.Verify.SetImage('locked')
		cl.Verify.Never().svcSocketClient([anyArgs:])

		// Standalone
		cl.SvcSettingsIconsControl_settings = [svc_local?:]
		cl.set()
		serverIcon.Verify.ToolTip('Server: Standalone')
		passwordIcon.Verify.ToolTip('Password: Not required')
		passwordIcon.Verify.SetImage('unlocked')
		// svcSocketClient(close?:) is called
		cl.Verify.svcSocketClient([anyArgs:])

		// Server, authentication failed
		cl.SvcSettingsIconsControl_settings = [svc_server: 'test']
		cl.When.connectionError().Return(connectionErr = 'Failed to Connect')
		cl.set()
		authFailure = cl.SvcSettingsIconsControl_authFailure
		serverIcon.Verify.ToolTip('Server: ' $ authFailure $ connectionErr)
		passwordIcon.Verify.ToolTip('Password: ' $ authFailure $ connectionErr)
		passwordIcon.Verify.SetImage('invalid_lock')
		// Ensure svcSocketClient(...) is not called since the last time
		cl.Verify.svcSocketClient([anyArgs:])

		// Server, authentication failed, incorrect password
		cl.When.connectionError().Return(passwordErr = SvcSocketClient.InvalidCredentials)
		cl.set()
		err = authFailure $ passwordErr
		serverIcon.Verify.ToolTip('Server: ' $ err)
		passwordIcon.Verify.ToolTip('Password: ' $ err)
		passwordIcon.Verify.Times(2).SetImage('invalid_lock')
		// Ensure svcSocketClient(...) is still not called
		cl.Verify.svcSocketClient([anyArgs:])

		// Server, authentication passed, with password
		cl.When.connectionError().Return('')
		cl.set()
		serverIcon.Verify.ToolTip('Server: test')
		passwordIcon.Verify.ToolTip('Password: Verified')
		passwordIcon.Verify.SetImage('valid_lock')
		// svcSocketClient(...) is called because it is not connected and there is
		// no current connection errors
		cl.Verify.Times(2).svcSocketClient([anyArgs:])
		}
	}
