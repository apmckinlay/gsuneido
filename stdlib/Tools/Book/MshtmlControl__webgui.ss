// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Controller
	{
	Name:		'Mshtml'
	Xstretch:	1
	Ystretch:	1

	New(.text = '', .allowReadOnly = false, style /*unused*/ = false)
		{
		.ctrl = .FindControl('WebBrowser')
		.ctrl.Load(.convertValue(text))
		.Send('Data')
		}

	convertValue(text)
		{
		if text.Prefix?('http://') or text.Prefix?('https://')
			return text
		else
			return "MSHTML:" $ text
		}

	Controls: #('WebBrowser', false)
	Set(.text)
		{
		.ctrl.Load(.convertValue(text))
		}
	Get()
		{
		return .text
		}
	SetReadOnly(readonly)
		{
		if .allowReadOnly is true
			super.SetReadOnly(readonly)
		}
	Dirty?(dirty/*unused*/ = "")
		{
		return false
		}
	Valid?(@unused)
		{
		return true
		}
	// need to RedirAccels for Ctrl + F
	On_Find()
		{
		SuneidoLog('INFO: Detect On_Find in use', calls:)
		}

	Default(@args)
		{
		return (.ctrl[args[0]])(@+1args)
		}

	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}
