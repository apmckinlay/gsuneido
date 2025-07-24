// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
IM_MessengerTabBase
	{
	TabName: 'Home'
	New()
		{
		if UserSettings.Get('IM_EnterToSend') is true
			.FindControl('enterToSend').Set(true)
		if UserSettings.Get('IM_AutoLoadHistory') is true
			.FindControl('autoLoadHistory').Set(true)
		}

	Controls()
		{
		return Object('Vert'
			Object('Mshtml', .headerText(), name: 'htmlEditor')
			#(Horz Fill
				(CheckBox '"Enter" key sends message' name: 'enterToSend')
				Fill
				(CheckBox, 'Auto-load History' name:'autoLoadHistory')
				Fill
			))
		}

	headerText()
		{
		return  '<html>' $ .CssStyle() $ '<body>' $
			'<center> <br><br> <h1> Welcome to Instant Messenger! </h1>'  $
			'Please double click on a user to start conversation'  $ '</body></html>'
		}

	AutoLoadHistory?()
		{
		return .FindControl('autoLoadHistory').Get() is true
		}

	EnterToSend?()
		{
		return .FindControl('enterToSend').Get() is true
		}

	Destroy()
		{
		UserSettings.Put('IM_EnterToSend', .FindControl('enterToSend').Get())
		UserSettings.Put('IM_AutoLoadHistory', .FindControl('autoLoadHistory').Get())
		super.Destroy()
		}
	}