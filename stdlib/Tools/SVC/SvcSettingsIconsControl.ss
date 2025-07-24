// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 			'SVC Status'
	Name: 			SvcIcons
	serverIcon: 	false
	passwordIcon: 	false
	New(.skipEtchedLine = false)
		{
		.serverIcon = .FindControl('svc_server_icon')
		.passwordIcon = .FindControl('svc_password_icon')
		.settingsChanged()
		.sub = PubSub.Subscribe('SvcSocketClient_StateChanged',
			{ .Defer(.settingsChanged, uniqueID: 'svcsettingsicons') })
		}

	Controls()
		{
		controls = Object('Vert')
		if not .skipEtchedLine
			controls.Add(#(EtchedLine, before: 0, after: 0))
		controls.Add(#(Horz svc_server_icon, Skip, svc_password_icon, Skip))
		return controls
		}

	settingsChanged()
		{
		.settings = SvcSettings.Get()
		if .serverIcon isnt false and .passwordIcon isnt false
			.set()
		}

	setting?: false // Prevent multiple calls as a result of inital SvcSocketClient open
	set()
		{
		if .setting?
			return
		.setting? = true
		if .insufficientSettings?()
			.invalid()
		else if .settings.svc_local? is true
			.standalone()
		else
			.serverBased()
		.setting? = false
		}

	insufficientSettings?()
		{
		if .settings is false
			return true
		if not Object?(.settings)
			return true
		return .settings.svc_local? is '' and .settings.svc_server is ''
		}

	insufficientSettings: 'Insufficient SVC Settings.'
	invalid()
		{
		.setServerIcon(.insufficientSettings, CLR.GRAY)
		.setPasswordIcon(.insufficientSettings, CLR.GRAY, 'locked')
		}

	serverStatus: ''
	setServerIcon(msg, clr)
		{ .setIcon(.serverIcon, .serverStatus = 'Server: ' $ msg, clr) }

	passwordStatus: ''
	setPasswordIcon(msg, clr, iconImage = false)
		{ .setIcon(.passwordIcon, .passwordStatus = 'Password: ' $ msg, clr, iconImage) }

	setIcon(icon, tip, color, iconImage = false)
		{
		icon.ToolTip(tip)
		icon.SetImageColor(color, color)
		if iconImage isnt false
			icon.SetImage(iconImage)
		}

	standalone()
		{
		.svcSocketClient(close?:)
		.setServerIcon('Standalone', CLR.EnhancedButtonFace)
		.setPasswordIcon('Not required', CLR.GRAY, 'unlocked')
		}

	svcSocketClient(close? = false)
		{
		if close?
			SvcSocketClient().Close()
		else if .connectionError() is ''
			SvcSocketClient().TestConnect(.settings.svc_server)
		}

	connectionError()
		{ return SvcSocketClient().Error }

	serverBased()
		{
		if not .connected?() and '' is err = .connectionError()
			.svcSocketClient()
		if '' isnt err = .connectionError()
			.invalidCredentials(err)
		else
			{
			.setServerIcon(.settings.svc_server, CLR.EnhancedButtonFace)
			.setPasswordIcon('Verified', CLR.EnhancedButtonFace, 'valid_lock')
			}
		}

	connected?()
		{ return SvcSocketClient().Connected? }

	authFailure: 'Authentication failed. '
	invalidCredentials(error)
		{
		servIconClr = passIconClr = CLR.DARKRED
		if error.Prefix?(SvcSocketClient.InvalidCredentials)
			servIconClr = CLR.EnhancedButtonFace
		else if error isnt SvcSocketClient.InvalidServer
			servIconClr = passIconClr = CLR.GRAY
		.setServerIcon(.authFailure $ error, servIconClr)
		.setPasswordIcon(.authFailure $ error, passIconClr, 'invalid_lock')
		}

	retry()
		{
		.svcSocketClient(close?:)
		.set()
		}

	On_Status() // Called Via On_IconMenu
		{
		if '' isnt err = .connectionError()
			{
			if YesNo('SvcSocketClient is not connected.\r\nLast error: ' $ err $
				'\r\n\r\nDo you want to retry?', .Title)
				.retry()
			}
		else
			.AlertInfo(.Title, .serverStatus.Has?(.insufficientSettings)
				? .noSettingsMsg()
				: .serverStatus $ '\r\n' $ .passwordStatus)
		}

	noSettingsMsg()
		{
		return .insufficientSettings $ '\r\n\r\nPlease verify Version Control Settings.'
		}

	On_Settings() // Called Via On_IconMenu
		{ SvcSettings(openDialog:) }

	Destroy()
		{
		.sub.Unsubscribe()
		super.Destroy()
		}
	}
